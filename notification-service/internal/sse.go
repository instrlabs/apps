package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
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

func (s *SSEService) HandleSSE(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("SSE connection attempt without token")
		http.Error(w, "Authentication token required", http.StatusUnauthorized)
		return
	}

	userId := extractUserIdFromToken(token)
	if userId == "" {
		log.Printf("Invalid token provided: %s", token)
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	ctx := r.Context()

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("ResponseWriter does not implement http.Flusher")
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	fmt.Fprintf(w, "event: connected\ndata: {\"connected\": true}\n\n")
	flusher.Flush()

	go func() {
		select {
		case <-ctx.Done():
			s.mutex.Lock()
			if s.clients[userId] == client {
				delete(s.clients, userId)
			}
			s.mutex.Unlock()
			client.done <- true
			log.Printf("SSE client disconnected for user %s. Total clients: %d", userId, len(s.clients))
		case <-client.done:
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-client.done:
			return
		case <-ctx.Done():
			return
		case msg := <-client.connection:
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": keepalive %v\n\n", time.Now())
			flusher.Flush()
		}
	}
}

func (s *SSEService) BroadcastNotification(ctx context.Context, notification *JobNotificationMessage) error {
	msgBytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if notification.UserID != "" {
		if client, exists := s.clients[notification.UserID]; exists {
			select {
			case client.connection <- msgBytes:
				log.Printf("Sent notification to user %s", notification.UserID)
			default:
				log.Printf("Failed to send notification to user %s: channel full", notification.UserID)
				client.done <- true
				delete(s.clients, notification.UserID)
				return fmt.Errorf("client channel full")
			}
		} else {
			log.Printf("User %s not connected, notification not delivered", notification.UserID)
		}

		return nil
	}

	sentCount := 0
	for userId, client := range s.clients {
		select {
		case client.connection <- msgBytes:
			sentCount++
		default:
			log.Printf("Failed to send message to user %s: channel full", userId)
			client.done <- true
			delete(s.clients, userId)
		}
	}

	log.Printf("Broadcasted notification to %d clients", sentCount)
	return nil
}

func extractUserIdFromToken(token string) string {
	return token
}
