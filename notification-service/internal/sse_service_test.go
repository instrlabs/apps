package internal

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// SSE Service Tests - Focus on core business logic
// ====================

func TestNewSSEService(t *testing.T) {
	cfg := &Config{
		Environment:                 "test",
		Port:                        ":3001",
		Origins:                     "http://localhost:8000",
		JWTSecret:                   "test-secret",
		AuthService:                 "http://auth-service:3000",
		NatsURI:                     "nats://localhost:4222",
		NatsSubjectNotificationsSSE: "notifications.sse",
	}

	service := NewSSEService(cfg)

	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.cfg)
	assert.NotNil(t, service.clients)
	assert.Len(t, service.clients, 0)
}

func TestSSEService_NotificationUser_Success(t *testing.T) {
	cfg := &Config{JWTSecret: "test-secret"}
	service := NewSSEService(cfg)

	// Setup mock client
	messageChan := make(chan []byte, 10)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      "user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients["user-123"] = client

	// Create notification
	notification := InstructionNotification{
		UserID:              "user-123",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}

	messageBytes, err := json.Marshal(notification)
	require.NoError(t, err)

	// Test
	service.NotificationUser(messageBytes)

	// Assert - Message delivered to client
	select {
	case receivedMessage := <-messageChan:
		var receivedNotification InstructionNotification
		err := json.Unmarshal(receivedMessage, &receivedNotification)
		assert.NoError(t, err)
		assert.Equal(t, "user-123", receivedNotification.UserID)
		assert.Equal(t, "instr-456", receivedNotification.InstructionID)
		assert.Equal(t, "detail-789", receivedNotification.InstructionDetailID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Message not received within timeout")
	}
}

func TestSSEService_NotificationUser_ClientNotFound(t *testing.T) {
	cfg := &Config{JWTSecret: "test-secret"}
	service := NewSSEService(cfg)

	// No clients registered

	notification := InstructionNotification{
		UserID:              "nonexistent-user",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}

	messageBytes, err := json.Marshal(notification)
	require.NoError(t, err)

	// Test - Should not panic or error
	service.NotificationUser(messageBytes)

	// Assert - No clients affected
	assert.Len(t, service.clients, 0)
}

func TestSSEService_NotificationUser_InvalidJSON(t *testing.T) {
	cfg := &Config{JWTSecret: "test-secret"}
	service := NewSSEService(cfg)

	// Setup mock client
	messageChan := make(chan []byte, 10)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      "user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients["user-123"] = client

	// Invalid JSON
	invalidJSON := []byte(`{invalid json`)

	// Test - Should not panic
	service.NotificationUser(invalidJSON)

	// Assert - No message sent to client
	select {
	case <-messageChan:
		t.Fatal("Should not receive any message for invalid JSON")
	case <-time.After(50 * time.Millisecond):
		// Expected - no message sent
	}
}

func TestSSEService_NotificationUser_EmptyMessage(t *testing.T) {
	cfg := &Config{JWTSecret: "test-secret"}
	service := NewSSEService(cfg)

	// Setup mock client
	messageChan := make(chan []byte, 10)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      "user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients["user-123"] = client

	// Empty message
	emptyMessage := []byte(`{}`)

	// Test
	service.NotificationUser(emptyMessage)

	// Assert - No message sent (no user_id in payload)
	select {
	case <-messageChan:
		t.Fatal("Should not receive message for empty user_id")
	case <-time.After(50 * time.Millisecond):
		// Expected - no message sent
	}
}

func TestSSEService_NotificationUser_ConcurrentNotifications(t *testing.T) {
	cfg := &Config{JWTSecret: "test-secret"}
	service := NewSSEService(cfg)

	// Setup mock client with buffered channel
	messageChan := make(chan []byte, 100)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      "user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients["user-123"] = client

	// Send multiple notifications concurrently
	notificationCount := 10
	done := make(chan bool)

	go func() {
		for i := 0; i < notificationCount; i++ {
			notification := InstructionNotification{
				UserID:              "user-123",
				InstructionID:       "instr-" + string(rune('0'+i)),
				InstructionDetailID: "detail-" + string(rune('0'+i)),
			}
			messageBytes, _ := json.Marshal(notification)
			service.NotificationUser(messageBytes)
		}
		done <- true
	}()

	<-done

	// Assert - All messages received
	receivedCount := 0
	timeout := time.After(500 * time.Millisecond)

loop:
	for {
		select {
		case <-messageChan:
			receivedCount++
			if receivedCount == notificationCount {
				break loop
			}
		case <-timeout:
			break loop
		}
	}

	assert.Equal(t, notificationCount, receivedCount, "Should receive all notifications")
}

func TestSSEService_NotificationUser_ChannelFull(t *testing.T) {
	cfg := &Config{JWTSecret: "test-secret"}
	service := NewSSEService(cfg)

	// Setup client with small buffer
	messageChan := make(chan []byte, 1)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      "user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients["user-123"] = client

	// Fill the channel
	notification1 := InstructionNotification{
		UserID:              "user-123",
		InstructionID:       "instr-1",
		InstructionDetailID: "detail-1",
	}
	msg1, _ := json.Marshal(notification1)
	service.NotificationUser(msg1)

	// Try to send another (channel full)
	notification2 := InstructionNotification{
		UserID:              "user-123",
		InstructionID:       "instr-2",
		InstructionDetailID: "detail-2",
	}
	msg2, _ := json.Marshal(notification2)

	// Should not block or panic
	service.NotificationUser(msg2)

	// Assert - At least first message received
	select {
	case receivedMessage := <-messageChan:
		var receivedNotification InstructionNotification
		err := json.Unmarshal(receivedMessage, &receivedNotification)
		assert.NoError(t, err)
		assert.Equal(t, "instr-1", receivedNotification.InstructionID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Should receive at least the first message")
	}
}

func TestSSEService_MultipleClients(t *testing.T) {
	cfg := &Config{JWTSecret: "test-secret"}
	service := NewSSEService(cfg)

	// Setup multiple clients
	messageChan1 := make(chan []byte, 10)
	messageChan2 := make(chan []byte, 10)

	client1 := &SSEClient{
		userId:      "user-1",
		connection:  messageChan1,
		done:        make(chan bool),
		connectedAt: time.Now().UTC(),
	}

	client2 := &SSEClient{
		userId:      "user-2",
		connection:  messageChan2,
		done:        make(chan bool),
		connectedAt: time.Now().UTC(),
	}

	service.clients["user-1"] = client1
	service.clients["user-2"] = client2

	// Send notification to user-1
	notification := InstructionNotification{
		UserID:              "user-1",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}
	messageBytes, _ := json.Marshal(notification)
	service.NotificationUser(messageBytes)

	// Assert - Only user-1 receives message
	select {
	case msg := <-messageChan1:
		var n InstructionNotification
		json.Unmarshal(msg, &n)
		assert.Equal(t, "user-1", n.UserID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("User-1 should receive message")
	}

	// User-2 should not receive anything
	select {
	case <-messageChan2:
		t.Fatal("User-2 should not receive message intended for user-1")
	case <-time.After(50 * time.Millisecond):
		// Expected - no message for user-2
	}
}
