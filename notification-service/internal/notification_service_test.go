package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockConfig() *Config {
	return &Config{
		Environment:                 "test",
		Port:                        ":3001",
		Origins:                     "http://localhost:8000",
		JWTSecret:                   "test-secret-key-for-testing-only",
		AuthService:                 "http://auth-service:3000",
		NatsURI:                     "nats://localhost:4222",
		NatsSubjectNotificationsSSE: "notifications.sse",
	}
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Environment)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.Origins)
	assert.NotEmpty(t, cfg.JWTSecret)
	assert.NotEmpty(t, cfg.AuthService)
	assert.NotEmpty(t, cfg.NatsURI)
	assert.NotEmpty(t, cfg.NatsSubjectNotificationsSSE)
}

func TestNewSSEService(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.cfg)
	assert.NotNil(t, service.clients)
	assert.Len(t, service.clients, 0)
}

func TestSSEService_NotificationUser_ValidMessage(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	// Setup a mock client
	messageChan := make(chan []byte, 16)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      "test-user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients["test-user-123"] = client

	// Create test notification
	notification := InstructionNotification{
		UserID:              "test-user-123",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}

	messageBytes, _ := json.Marshal(notification)

	// Test
	service.NotificationUser(messageBytes)

	// Assert - Message should be sent to client
	select {
	case receivedMessage := <-messageChan:
		var receivedNotification InstructionNotification
		err := json.Unmarshal(receivedMessage, &receivedNotification)
		assert.NoError(t, err)
		assert.Equal(t, notification.UserID, receivedNotification.UserID)
		assert.Equal(t, notification.InstructionID, receivedNotification.InstructionID)
		assert.Equal(t, notification.InstructionDetailID, receivedNotification.InstructionDetailID)
	case <-time.After(100 * time.Millisecond):
		t.Error("Message was not received within timeout")
	}
}

func TestSSEService_NotificationUser_InvalidMessage(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	// Setup a mock client
	messageChan := make(chan []byte, 16)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      "test-user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients["test-user-123"] = client

	// Test with invalid JSON
	invalidMessage := []byte("{ invalid json }")

	// Should not panic
	service.NotificationUser(invalidMessage)

	// Assert - No message should be sent due to invalid JSON
	select {
	case <-messageChan:
		t.Error("Message should not be sent for invalid JSON")
	case <-time.After(50 * time.Millisecond):
		// Expected behavior - no message sent
	}
}

func TestSSEService_NotificationUser_ClientNotFound(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	notification := InstructionNotification{
		UserID:              "non-existent-user",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}

	messageBytes, _ := json.Marshal(notification)

	// Should not panic when client doesn't exist
	service.NotificationUser(messageBytes)

	// Expected behavior - no panic, message ignored
	assert.True(t, true)
}

func TestExtractTokenInfo_ValidToken(t *testing.T) {
	secret := "test-secret-key-for-testing-only"
	userID := "test-user-123"
	roles := []string{"user", "admin"}

	// Create valid JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Test
	tokenInfo, err := ExtractTokenInfo(secret, tokenString)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tokenInfo)
	assert.Equal(t, userID, tokenInfo.UserID)
	assert.Equal(t, roles, tokenInfo.Roles)
}

func TestExtractTokenInfo_ExpiredToken(t *testing.T) {
	secret := "test-secret-key-for-testing-only"
	userID := "test-user-123"

	// Create expired JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	})

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Test
	tokenInfo, err := ExtractTokenInfo(secret, tokenString)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired") // JWT library returns wrapped error
	assert.Nil(t, tokenInfo)
}

func TestExtractTokenInfo_InvalidToken(t *testing.T) {
	secret := "test-secret-key-for-testing-only"
	wrongSecret := "wrong-secret-key"

	// Create valid token with one secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "test-user-123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Try to extract with wrong secret
	tokenInfo, err := ExtractTokenInfo(wrongSecret, tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tokenInfo)
}

func TestExtractTokenInfo_EmptyToken(t *testing.T) {
	secret := "test-secret-key-for-testing-only"

	// Test with empty token
	tokenInfo, err := ExtractTokenInfo(secret, "")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrTokenEmpty, err)
	assert.Nil(t, tokenInfo)
}

func TestExtractTokenInfo_EmptySecret(t *testing.T) {
	tokenString := "some-token-string"

	// Test with empty secret - should fail to parse
	tokenInfo, err := ExtractTokenInfo("", tokenString)

	// Assert - will fail to parse token with empty secret
	assert.Error(t, err)
	assert.Nil(t, tokenInfo)
}

func TestExtractTokenInfo_WhitespaceToken(t *testing.T) {
	secret := "test-secret-key-for-testing-only"

	// Test with whitespace-only token
	tokenInfo, err := ExtractTokenInfo(secret, "   ")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrTokenEmpty, err)
	assert.Nil(t, tokenInfo)
}

func TestExtractRoles_StringSlice(t *testing.T) {
	roles := []string{"user", "admin", "moderator"}
	result := extractRoles(roles)

	assert.Equal(t, roles, result)
}

func TestExtractRoles_InterfaceSlice(t *testing.T) {
	roles := []interface{}{"user", "admin", "moderator", "", "  ", "editor"}
	result := extractRoles(roles)

	expected := []string{"user", "admin", "moderator", "editor"}
	assert.Equal(t, expected, result)
}

func TestExtractRoles_StringWithCommas(t *testing.T) {
	rolesString := "user,admin, moderator, editor ,  "
	result := extractRoles(rolesString)

	expected := []string{"user", "admin", "moderator", "editor"}
	assert.Equal(t, expected, result)
}

func TestExtractRoles_NilInput(t *testing.T) {
	result := extractRoles(nil)
	assert.Nil(t, result)
}

func TestExtractRoles_EmptyString(t *testing.T) {
	result := extractRoles("")
	assert.Equal(t, []string{}, result) // The function returns empty slice, not nil
}

func TestExtractRoles_WhitespaceString(t *testing.T) {
	result := extractRoles("   ")
	assert.Equal(t, []string{}, result) // The function returns empty slice, not nil
}

func TestExtractRoles_EmptySlice(t *testing.T) {
	var emptySlice []interface{}
	result := extractRoles(emptySlice)
	assert.Equal(t, []string{}, result) // The function returns empty slice, not nil
}

func TestExtractRoles_UnexpectedType(t *testing.T) {
	result := extractRoles(123)
	assert.Nil(t, result)
}

func TestToString_StringInput(t *testing.T) {
	input := "test string"
	result := toString(input)

	assert.Equal(t, input, result)
}

func TestToString_Float64Input(t *testing.T) {
	input := float64(123.456)
	result := toString(input)

	assert.Equal(t, "123.456", result)
}

func TestToString_Float64IntegerInput(t *testing.T) {
	input := float64(123)
	result := toString(input)

	assert.Equal(t, "123", result)
}

func TestToString_Int64Input(t *testing.T) {
	input := int64(456)
	result := toString(input)

	assert.Equal(t, "456", result)
}

func TestToString_IntInput(t *testing.T) {
	input := int(789)
	result := toString(input)

	assert.Equal(t, "789", result)
}

func TestToString_NilInput(t *testing.T) {
	result := toString(nil)
	assert.Equal(t, "", result)
}

func TestToString_UnexpectedType(t *testing.T) {
	result := toString([]string{"test"})
	assert.Equal(t, "", result)
}

func TestFmtFloat(t *testing.T) {
	assert.Equal(t, "123.456", fmtFloat(123.456))
	assert.Equal(t, "123", fmtFloat(123.0))
	assert.Equal(t, "0", fmtFloat(0.0))
	assert.Equal(t, "-123.456", fmtFloat(-123.456))
}

func TestFmtInt(t *testing.T) {
	assert.Equal(t, "123", fmtInt(123))
	assert.Equal(t, "0", fmtInt(0))
	assert.Equal(t, "-456", fmtInt(-456))
	assert.Equal(t, "9223372036854775807", fmtInt(9223372036854775807)) // Max int64
}

func TestInstructionNotification_Validation(t *testing.T) {
	notification := InstructionNotification{
		UserID:              "user-123",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}

	assert.NotEmpty(t, notification.UserID)
	assert.NotEmpty(t, notification.InstructionID)
	assert.NotEmpty(t, notification.InstructionDetailID)
}

func TestInstructionNotification_JSONSerialization(t *testing.T) {
	notification := InstructionNotification{
		UserID:              "user-123",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}

	// Test JSON marshaling
	jsonBytes, err := json.Marshal(notification)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Test JSON unmarshaling
	var unmarshaled InstructionNotification
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, notification.UserID, unmarshaled.UserID)
	assert.Equal(t, notification.InstructionID, unmarshaled.InstructionID)
	assert.Equal(t, notification.InstructionDetailID, unmarshaled.InstructionDetailID)
}

func TestSSEClient_Creation(t *testing.T) {
	messageChan := make(chan []byte, 16)
	doneChan := make(chan bool)
	now := time.Now().UTC()

	client := &SSEClient{
		userId:      "test-user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: now,
	}

	assert.Equal(t, "test-user-123", client.userId)
	assert.Equal(t, messageChan, client.connection)
	assert.Equal(t, doneChan, client.done)
	assert.Equal(t, now, client.connectedAt)
}

func TestSSEService_ConcurrentAccess(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	// Test concurrent client additions
	done := make(chan bool, 2)

	go func() {
		messageChan := make(chan []byte, 16)
		doneChan := make(chan bool)
		client := &SSEClient{
			userId:      "user-1",
			connection:  messageChan,
			done:        doneChan,
			connectedAt: time.Now().UTC(),
		}
		service.mutex.Lock()
		service.clients["user-1"] = client
		service.mutex.Unlock()
		done <- true
	}()

	go func() {
		messageChan := make(chan []byte, 16)
		doneChan := make(chan bool)
		client := &SSEClient{
			userId:      "user-2",
			connection:  messageChan,
			done:        doneChan,
			connectedAt: time.Now().UTC(),
		}
		service.mutex.Lock()
		service.clients["user-2"] = client
		service.mutex.Unlock()
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	service.mutex.Lock()
	assert.Len(t, service.clients, 2)
	service.mutex.Unlock()
}

func TestExtractTokenInfo_MalformedToken(t *testing.T) {
	secret := "test-secret-key-for-testing-only"

	// Test with completely invalid token
	malformedTokens := []string{
		"not.a.jwt.token",
		"invalid",
		"Bearer malformed",
		"eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.",
	}

	for _, token := range malformedTokens {
		tokenInfo, err := ExtractTokenInfo(secret, token)
		assert.Error(t, err, "Token %q should be invalid", token)
		assert.Nil(t, tokenInfo, "Token %q should return nil token info", token)
	}
}

func TestExtractTokenInfo_MissingUserID(t *testing.T) {
	secret := "test-secret-key-for-testing-only"

	// Create token without user_id claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roles": []string{"user"},
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Test
	tokenInfo, err := ExtractTokenInfo(secret, tokenString)

	// Assert - Should succeed but with empty userID
	assert.NoError(t, err)
	assert.NotNil(t, tokenInfo)
	assert.Equal(t, "", tokenInfo.UserID)
}

func TestExtractTokenInfo_TokenWithoutExpiration(t *testing.T) {
	secret := "test-secret-key-for-testing-only"
	userID := "test-user-123"

	// Create token without exp claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"roles":   []string{"user"},
	})

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Test
	tokenInfo, err := ExtractTokenInfo(secret, tokenString)

	// Assert - Should succeed without expiration check
	assert.NoError(t, err)
	assert.NotNil(t, tokenInfo)
	assert.Equal(t, userID, tokenInfo.UserID)
}

func TestSSEService_ClientReplacement(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	// Create initial client
	messageChan1 := make(chan []byte, 16)
	doneChan1 := make(chan bool, 1) // Buffered to prevent blocking
	client1 := &SSEClient{
		userId:      "test-user-123",
		connection:  messageChan1,
		done:        doneChan1,
		connectedAt: time.Now().UTC(),
	}

	service.clients["test-user-123"] = client1
	assert.Len(t, service.clients, 1)

	// Create new client with same userId
	messageChan2 := make(chan []byte, 16)
	doneChan2 := make(chan bool, 1) // Buffered to prevent blocking
	client2 := &SSEClient{
		userId:      "test-user-123",
		connection:  messageChan2,
		done:        doneChan2,
		connectedAt: time.Now().UTC(),
	}

	// Simulate replacing existing client
	service.mutex.Lock()
	if existingClient, exists := service.clients["test-user-123"]; exists {
		select {
		case existingClient.done <- true: // Signal old client to close
		case <-time.After(10 * time.Millisecond):
			// Skip if channel is blocked
		}
	}
	service.clients["test-user-123"] = client2
	service.mutex.Unlock()

	assert.Len(t, service.clients, 1)
	assert.Equal(t, client2, service.clients["test-user-123"])
}

// Middleware tests for token extraction
func TestSetupMiddleware_TokenExtraction(t *testing.T) {
	cfg := newMockConfig()
	app := fiber.New()
	SetupMiddleware(app, cfg)

	// Add a test endpoint that uses x-user-id header
	app.Get("/test", func(c *fiber.Ctx) error {
		userID := c.Get("x-user-id")
		return c.JSON(fiber.Map{
			"user_id": userID,
		})
	})

	// Test with valid JWT token
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "test-user-456",
		"roles":   []string{"user"},
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := validToken.SignedString([]byte(cfg.JWTSecret))
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "test-user-456", result["user_id"])
}

func TestSetupMiddleware_MissingToken(t *testing.T) {
	cfg := newMockConfig()
	app := fiber.New()
	SetupMiddleware(app, cfg)

	app.Get("/test", func(c *fiber.Ctx) error {
		userID := c.Get("x-user-id")
		return c.JSON(fiber.Map{
			"user_id": userID,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	// User ID should be empty when no token is provided
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Should return empty user_id since no token (middleware sets it to empty string)
	assert.Equal(t, "", result["user_id"])
}

func TestSetupMiddleware_InvalidTokenFormat(t *testing.T) {
	cfg := newMockConfig()
	app := fiber.New()
	SetupMiddleware(app, cfg)

	app.Get("/test", func(c *fiber.Ctx) error {
		userID := c.Get("x-user-id")
		return c.JSON(fiber.Map{
			"user_id": userID,
		})
	})

	// Test with invalid Bearer format
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "BearerInvalidFormat")

	resp, err := app.Test(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Should have empty user_id since token is malformed
	assert.Equal(t, "", result["user_id"])
}

func TestSetupMiddleware_ExpiredToken(t *testing.T) {
	cfg := newMockConfig()
	app := fiber.New()
	SetupMiddleware(app, cfg)

	app.Get("/test", func(c *fiber.Ctx) error {
		userID := c.Get("x-user-id")
		return c.JSON(fiber.Map{
			"user_id": userID,
		})
	})

	// Create expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "test-user",
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	})

	tokenString, err := expiredToken.SignedString([]byte(cfg.JWTSecret))
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := app.Test(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Should have empty user_id since token is expired
	assert.Equal(t, "", result["user_id"])
}

// Additional SSE Handler tests
func TestSSEService_HandleSSE_Success(t *testing.T) {
	// This test verifies SSE connection setup
	// Note: Testing streaming endpoints is complex in unit tests since they block
	// This test verifies the service initialization only
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	// Verify service is properly initialized
	assert.NotNil(t, service)
	assert.NotNil(t, service.clients)
	assert.Equal(t, 0, len(service.clients))

	// Verify SSE client creation works
	messageChan := make(chan []byte, 16)
	doneChan := make(chan bool, 1)
	client := &SSEClient{
		userId:      "test-user-123",
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.mutex.Lock()
	service.clients["test-user-123"] = client
	service.mutex.Unlock()

	assert.Len(t, service.clients, 1)
	assert.Equal(t, client, service.clients["test-user-123"])
}

// Multiple notification test
func TestSSEService_NotificationUser_MultipleClients(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	// Create multiple clients
	users := []string{"user-1", "user-2", "user-3"}
	clientChans := make(map[string]chan []byte)

	for _, userID := range users {
		messageChan := make(chan []byte, 16)
		doneChan := make(chan bool)
		client := &SSEClient{
			userId:      userID,
			connection:  messageChan,
			done:        doneChan,
			connectedAt: time.Now().UTC(),
		}
		service.clients[userID] = client
		clientChans[userID] = messageChan
	}

	// Send notification to specific user
	notification := InstructionNotification{
		UserID:              "user-2",
		InstructionID:       "instr-456",
		InstructionDetailID: "detail-789",
	}

	messageBytes, _ := json.Marshal(notification)
	service.NotificationUser(messageBytes)

	// Verify only user-2 receives message
	select {
	case msg := <-clientChans["user-2"]:
		var received InstructionNotification
		json.Unmarshal(msg, &received)
		assert.Equal(t, "user-2", received.UserID)
	case <-time.After(100 * time.Millisecond):
		t.Error("User-2 should receive message")
	}

	// Verify user-1 and user-3 don't receive message
	select {
	case <-clientChans["user-1"]:
		t.Error("User-1 should not receive message")
	case <-time.After(50 * time.Millisecond):
		// Expected - user-1 should not receive
	}

	select {
	case <-clientChans["user-3"]:
		t.Error("User-3 should not receive message")
	case <-time.After(50 * time.Millisecond):
		// Expected - user-3 should not receive
	}
}

// Test InstructionNotification with different data
func TestInstructionNotification_VariousPayloads(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		instrID  string
		detailID string
	}{
		{"simple ids", "user-1", "instr-1", "detail-1"},
		{"uuid-like ids", "550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001", "550e8400-e29b-41d4-a716-446655440002"},
		{"numeric ids", "123456", "789012", "345678"},
		{"long ids", "user-very-long-id-that-is-still-valid", "instruction-very-long-id", "detail-very-long-id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notification := InstructionNotification{
				UserID:              tt.userID,
				InstructionID:       tt.instrID,
				InstructionDetailID: tt.detailID,
			}

			// Test JSON marshaling
			jsonBytes, err := json.Marshal(notification)
			assert.NoError(t, err)

			// Test unmarshaling
			var decoded InstructionNotification
			err = json.Unmarshal(jsonBytes, &decoded)
			assert.NoError(t, err)

			assert.Equal(t, notification.UserID, decoded.UserID)
			assert.Equal(t, notification.InstructionID, decoded.InstructionID)
			assert.Equal(t, notification.InstructionDetailID, decoded.InstructionDetailID)
		})
	}
}

// Test token extraction edge cases
func TestExtractTokenInfo_NumericUserID(t *testing.T) {
	secret := "test-secret-key-for-testing-only"
	userID := float64(12345) // JWT treats numbers as float64

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Test
	tokenInfo, err := ExtractTokenInfo(secret, tokenString)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tokenInfo)
	assert.Equal(t, "12345", tokenInfo.UserID) // Should convert to string
}

// Config loading tests
func TestNewConfig_DefaultValues(t *testing.T) {
	cfg := NewConfig()

	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.Origins)
	assert.NotEmpty(t, cfg.NatsURI)
	assert.NotEmpty(t, cfg.NatsSubjectNotificationsSSE)
}

// Concurrent notification delivery test
func TestSSEService_ConcurrentNotifications(t *testing.T) {
	cfg := newMockConfig()
	service := NewSSEService(cfg)

	userID := "concurrent-user"
	messageChan := make(chan []byte, 32)
	doneChan := make(chan bool)
	client := &SSEClient{
		userId:      userID,
		connection:  messageChan,
		done:        doneChan,
		connectedAt: time.Now().UTC(),
	}

	service.clients[userID] = client

	// Send multiple notifications concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			notification := InstructionNotification{
				UserID:              userID,
				InstructionID:       "instr-" + fmt.Sprintf("%d", index),
				InstructionDetailID: "detail-" + fmt.Sprintf("%d", index),
			}
			messageBytes, _ := json.Marshal(notification)
			service.NotificationUser(messageBytes)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify messages received
	receivedCount := 0
	timeout := time.NewTimer(500 * time.Millisecond)
	defer timeout.Stop()

	for receivedCount < 10 {
		select {
		case <-messageChan:
			receivedCount++
		case <-timeout.C:
			break
		}
	}

	assert.Equal(t, 10, receivedCount, "Should receive all 10 notifications")
}

// Test handling empty or malformed notifications
func TestSSEService_NotificationUser_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		message     []byte
		shouldPanic bool
	}{
		{"empty bytes", []byte{}, false},
		{"null bytes", nil, false},
		{"random data", []byte("random data that is not json"), false},
		{"incomplete json", []byte("{\"user_id\": \"user\""), false},
		{"json array", []byte("[1,2,3]"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newMockConfig()
			service := NewSSEService(cfg)

			// Add a client
			messageChan := make(chan []byte, 16)
			doneChan := make(chan bool)
			client := &SSEClient{
				userId:      "test-user",
				connection:  messageChan,
				done:        doneChan,
				connectedAt: time.Now().UTC(),
			}
			service.clients["test-user"] = client

			// Should not panic
			if tt.message == nil {
				// Skip nil test
				return
			}

			service.NotificationUser(tt.message)
			assert.True(t, true, "Should not panic")
		})
	}
}
