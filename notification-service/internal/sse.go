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
		fmt.Fprintf(w, "event: connected\ndata: %s\n\n", `{"connected": true}`)
		w.Flush()

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
				fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
				w.Flush()
			case <-ticker.C:
				fmt.Fprintf(w, ": keepalive %v\n\n", time.Now())
				w.Flush()
			}
		}
	})

	return nil
}

// SendTestNotification handles posting a test notification to a specific user or broadcasting to all.
func (s *SSEService) SendTestNotification(c *fiber.Ctx) error {
	type reqBody struct {
		UserID  string `json:"userId"`
		Message string `json:"message"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if body.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is required"})
	}

	payload := map[string]any{
		"message":   body.Message,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	}
	bytes, _ := json.Marshal(payload)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if body.UserID != "" {
		client, ok := s.clients[body.UserID]
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not connected"})
		}
		select {
		case client.connection <- bytes:
			return c.JSON(fiber.Map{"status": "sent", "userId": body.UserID})
		default:
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "client channel blocked"})
		}
	}

	// Broadcast to all clients
	count := 0
	for _, client := range s.clients {
		select {
		case client.connection <- bytes:
			count++
		default:
			// skip blocked clients
		}
	}
	return c.JSON(fiber.Map{"status": "broadcast", "delivered": count})
}
