package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SSEClient struct {
	userId      string
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
	var envelope struct {
		UserID string `json:"user_id"`
	}
	if err := json.Unmarshal(msg, &envelope); err == nil {
		s.mutex.Lock()
		client := s.clients[envelope.UserID]
		s.mutex.Unlock()
		if client != nil {
			select {
			case client.connection <- msg:
			default:
			}
		}
		return
	}
}

func NewSSEService(cfg *Config) *SSEService {
	return &SSEService{
		cfg:     cfg,
		clients: make(map[string]*SSEClient),
	}
}

func (s *SSEService) HandleSSE(c *fiber.Ctx) error {
	userId := c.Locals("UserID").(string)
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	messageChan := make(chan []byte, 16)
	doneChan := make(chan bool)

	client := &SSEClient{
		userId:      userId,
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	s.mutex.Lock()
	if existingClient, exists := s.clients[userId]; exists {
		log.Infof("Closing existing connection for user %s", userId)
		existingClient.done <- true
	}
	s.clients[userId] = client
	s.mutex.Unlock()

	log.Infof("New SSE client connected for user %s. Total clients: %d", userId, len(s.clients))

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
				if s.clients[userId] == client {
					delete(s.clients, userId)
				}
				s.mutex.Unlock()
				return
			case <-ctx.Done():
				s.mutex.Lock()
				if s.clients[userId] == client {
					delete(s.clients, userId)
				}
				s.mutex.Unlock()
				client.done <- true
				log.Infof("SSE client disconnected for user %s. Total clients: %d", userId, len(s.clients))
				return
			case data := <-client.connection:
				var msg InstructionNotification
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
