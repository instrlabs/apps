package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/shared/modelx"
)

type SSEClient struct {
	identifier  string
	connection  chan []byte
	done        chan bool
	connectedAt time.Time
}

type SSEService struct {
	cfg     *Config
	clients map[string]*SSEClient
	mutex   sync.Mutex
}

func (s *SSEService) NotificationUser(msg []byte) {
	var notification modelx.InstructionNotification
	if err := json.Unmarshal(msg, &notification); err == nil {
		var identifier string
		if notification.UserID != nil {
			identifier = notification.UserID.Hex()
		} else if notification.GuestID != nil {
			identifier = *notification.GuestID
		}

		if identifier != "" {
			s.mutex.Lock()
			client := s.clients[identifier]
			s.mutex.Unlock()
			if client != nil {
				select {
				case client.connection <- msg:
				default:
				}
			}
		}
	}
}

func NewSSEService(cfg *Config) *SSEService {
	return &SSEService{
		cfg:     cfg,
		clients: make(map[string]*SSEClient),
	}
}

func (s *SSEService) HandleSSE(c *fiber.Ctx) error {
	userId := c.Get("x-user-id")
	guestId := c.Get("x-guest-id")

	var identifier string

	if userId != "" {
		identifier = userId
	} else if guestId != "" {
		identifier = guestId
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing user identification - provide either x-user-id or x-guest-id header",
			"errors":  nil,
			"data":    nil,
		})
	}

	c.Set("content-type", "text/event-stream")
	c.Set("cache-control", "no-cache")
	c.Set("connection", "keep-alive")

	messageChan := make(chan []byte, 16)
	doneChan := make(chan bool)

	client := &SSEClient{
		identifier:  identifier,
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	s.mutex.Lock()
	if existingClient, exists := s.clients[identifier]; exists {
		log.Infof("Closing existing connection for %s", identifier)
		existingClient.done <- true
	}
	s.clients[identifier] = client
	s.mutex.Unlock()

	log.Infof("New SSE client connected for %s. Total clients: %d", identifier, len(s.clients))

	ctx := c.Context()

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		info := map[string]interface{}{"connected": true}
		infoJson, _ := json.Marshal(info)
		_, _ = fmt.Fprintf(w, "event: connected\ndata: %s\n\n", string(infoJson))
		_ = w.Flush()

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-client.done:
				s.mutex.Lock()
				if s.clients[identifier] == client {
					delete(s.clients, identifier)
				}
				s.mutex.Unlock()
				return
			case <-ctx.Done():
				s.mutex.Lock()
				if s.clients[identifier] == client {
					delete(s.clients, identifier)
				}
				s.mutex.Unlock()
				client.done <- true
				log.Infof("SSE client disconnected for %s. Total clients: %d", identifier, len(s.clients))
				return
			case data := <-client.connection:
				var msg modelx.InstructionNotification
				if err := json.Unmarshal(data, &msg); err != nil {
					log.Infof("RunInstructionMessage: unmarshal error: %v", err)
					return
				}

				msgJson, _ := json.Marshal(msg)
				_, _ = fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(msgJson))
				_ = w.Flush()
			case <-ticker.C:
				info := map[string]interface{}{"time": time.Now().UTC()}
				infoJson, _ := json.Marshal(info)
				_, _ = fmt.Fprintf(w, "event: ping\ndata: %s\n\n", string(infoJson))
				_ = w.Flush()
			}
		}
	})

	return nil
}
