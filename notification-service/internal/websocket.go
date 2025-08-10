package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	cfg      *Config
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mutex    sync.Mutex
}

func NewWebSocketService(cfg *Config) *WebSocketService {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins in development
			// In production, this should be more restrictive
			return true
		},
	}

	return &WebSocketService{
		cfg:      cfg,
		upgrader: upgrader,
		clients:  make(map[*websocket.Conn]bool),
	}
}

// HandleWebSocket handles WebSocket connections
func (ws *WebSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to WebSocket: %v", err)
		return
	}

	// Add client to the clients map
	ws.mutex.Lock()
	ws.clients[conn] = true
	ws.mutex.Unlock()

	log.Printf("New WebSocket client connected. Total clients: %d", len(ws.clients))

	// Handle disconnection
	go func() {
		defer func() {
			conn.Close()
			ws.mutex.Lock()
			delete(ws.clients, conn)
			ws.mutex.Unlock()
			log.Printf("WebSocket client disconnected. Total clients: %d", len(ws.clients))
		}()

		// Simple ping/pong to keep connection alive
		for {
			// Read message (we don't actually use the content)
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Client disconnected or error occurred
				break
			}
		}
	}()
}

func (ws *WebSocketService) BroadcastNotification(ctx context.Context, notification *JobNotificationMessage) error {
	msgBytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			client.Close()
			delete(ws.clients, client)
		}
	}

	log.Printf("Broadcasted notification to %d clients", len(ws.clients))
	return nil
}
