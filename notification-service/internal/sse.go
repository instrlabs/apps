package internal

import (
	"bufio"
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

	messageChan := make(chan []byte)
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
		_, _ = fmt.Fprintf(w, "event: connected\ndata: %s\n\n", `{"connected": true}`)
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
				_, _ = fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
				_ = w.Flush()
			case <-ticker.C:
				_, _ = fmt.Fprintf(w, "event: ping\ndata:  %v\n\n", time.Now())
				_ = w.Flush()
			}
		}
	})

	return nil
}
