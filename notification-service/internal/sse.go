package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (s *SSEService) Broadcast(msg []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, client := range s.clients {
		select {
		case client.connection <- msg:
		default:
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
	userId := c.Locals("UserID").(string)
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte, 16)
	doneChan := make(chan bool)

	client := &SSEClient{
		userId:      userId,
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now(),
	}

	s.mutex.Lock()
	if existingClient, exists := s.clients[userId]; exists {
		log.Printf("Closing existing connection for user %s", userId)
		existingClient.done <- true
	}
	s.clients[userId] = client
	s.mutex.Unlock()

	log.Printf("New SSE client connected for user %s. Total clients: %d", userId, len(s.clients))

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
				log.Printf("SSE client disconnected for user %s. Total clients: %d", userId, len(s.clients))
				return
			case msg := <-client.connection:
				msgJson, _ := json.Marshal(msg)
				_, _ = fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(msgJson))
				_ = w.Flush()
			case <-ticker.C:
				info := map[string]interface{}{"time": time.Now()}
				infoJson, _ := json.Marshal(info)
				_, _ = fmt.Fprintf(w, "event: ping\ndata: %s\n\n", string(infoJson))
				_ = w.Flush()
			}
		}
	})

	return nil
}
