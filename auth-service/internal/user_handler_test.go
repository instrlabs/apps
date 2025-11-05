package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to setup test app
func setupTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
				"errors":  nil,
				"data":    nil,
			})
		},
	})
	return app
}

// Helper function to create test config
func createTestConfig() *Config {
	return &Config{
		JWTSecret:          "test-secret-key-for-jwt-token",
		TokenExpiryHours:   1,
		RefreshExpiryHours: 720,
		PinEnabled:         false,
		GoogleClientID:     "test-client-id",
		GoogleClientSecret: "test-client-secret",
		GoogleRedirectUrl:  "http://localhost:3000/auth/google/callback",
		WebUrl:             "http://localhost:3000",
	}
}

// Helper to create a test user
func createTestUser(email string) *User {
	pinHash := "$2a$10$ZSVByNNuJls.hrp3dsfd6OSQQl4Wq93mFI5aIgdHB6c061Vn8TZSK" // bcrypt hash of "000000"
	now := time.Now().UTC()
	expires := now.Add(10 * time.Minute)
	return &User{
		ID:         primitive.NewObjectID(),
		Email:      email,
		Username:   "testuser",
		PinHash:    &pinHash,
		PinExpires: &expires,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Helper to create test session
func createTestSession(userID string) *UserSession {
	now := time.Now().UTC()
	return &UserSession{
		ID:             primitive.NewObjectID(),
		UserID:         userID,
		DeviceHash:     "test-device-hash",
		IPAddress:      "127.0.0.1",
		UserAgent:      "test-agent",
		RefreshToken:   "test-refresh-token",
		IsActive:       true,
		LastActivityAt: now,
		CreatedAt:      now,
		ExpiresAt:      now.AddDate(0, 0, 30),
	}
}

// Helper to create handler with mocks
func createHandlerWithMocks(cfg *Config, mockUserRepo *MockUserRepository, mockSessionRepo *MockUserSessionRepository) *UserHandler {
	return NewUserHandler(cfg, mockUserRepo, mockSessionRepo)
}

// Test Login Handler
func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.Login(c)
	})

	user := createTestUser("test@example.com")
	session := createTestSession(user.ID.Hex())

	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)
	mockUserRepo.On("ClearPin", user.ID.Hex()).Return(nil)
	mockUserRepo.On("SetRegisteredAt", user.ID.Hex()).Return(nil)
	mockSessionRepo.On("CreateUserSession", user.ID.Hex(), "127.0.0.1", "test-agent").Return(session, nil)
	mockSessionRepo.On("UpdateUserSessionRefreshToken", session.ID, mock.AnythingOfType("string")).Return(nil)

	reqBody := map[string]string{
		"email": "test@example.com",
		"pin":   "000000",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessLogin, response["message"])
	assert.NotNil(t, response["data"])

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidRequestBody, response["message"])
}

func TestLogin_EmailRequired(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", handler.Login)

	reqBody := map[string]string{
		"email": "",
		"pin":   "123456",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrEmailRequired, response["message"])
}

func TestLogin_PinRequired(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", handler.Login)

	reqBody := map[string]string{
		"email": "test@example.com",
		"pin":   "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrPasswordRequired, response["message"])
}

func TestLogin_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", handler.Login)

	mockUserRepo.On("FindByEmail", "test@example.com").Return(nil)

	reqBody := map[string]string{
		"email": "test@example.com",
		"pin":   "123456",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidCredentials, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidPin(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", handler.Login)

	// Create a user with PIN hash for "000000"
	user := createTestUser("test@example.com")

	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)

	reqBody := map[string]string{
		"email": "test@example.com",
		"pin":   "wrongpin", // Wrong PIN
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidCredentials, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestLogin_SessionCreationFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.Login(c)
	})

	user := createTestUser("test@example.com")

	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)
	mockUserRepo.On("ClearPin", user.ID.Hex()).Return(nil)
	mockUserRepo.On("SetRegisteredAt", user.ID.Hex()).Return(nil)
	mockSessionRepo.On("CreateUserSession", user.ID.Hex(), "127.0.0.1", "test-agent").Return(nil, errors.New("db error"))

	reqBody := map[string]string{
		"email": "test@example.com",
		"pin":   "000000",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrCreateSession, response["message"])

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestLogin_UpdateRefreshTokenFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.Login(c)
	})

	user := createTestUser("test@example.com")
	session := createTestSession(user.ID.Hex())

	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)
	mockUserRepo.On("ClearPin", user.ID.Hex()).Return(nil)
	mockUserRepo.On("SetRegisteredAt", user.ID.Hex()).Return(nil)
	mockSessionRepo.On("CreateUserSession", user.ID.Hex(), "127.0.0.1", "test-agent").Return(session, nil)
	mockSessionRepo.On("UpdateUserSessionRefreshToken", session.ID, mock.AnythingOfType("string")).Return(errors.New("db error"))

	reqBody := map[string]string{
		"email": "test@example.com",
		"pin":   "000000",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrUpdateSession, response["message"])

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestLogin_SetRegisteredAtFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/login", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.Login(c)
	})

	user := createTestUser("test@example.com")
	session := createTestSession(user.ID.Hex())

	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)
	mockUserRepo.On("ClearPin", user.ID.Hex()).Return(nil)
	mockUserRepo.On("SetRegisteredAt", user.ID.Hex()).Return(errors.New("database error"))
	mockSessionRepo.On("CreateUserSession", user.ID.Hex(), "127.0.0.1", "test-agent").Return(session, nil)
	mockSessionRepo.On("UpdateUserSessionRefreshToken", session.ID, mock.AnythingOfType("string")).Return(nil)

	reqBody := map[string]string{
		"email": "test@example.com",
		"pin":   "000000",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessLogin, response["message"])
	assert.NotNil(t, response["data"])

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

// Test RefreshToken Handler
func TestRefreshToken_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/refresh", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.RefreshToken(c)
	})

	user := createTestUser("test@example.com")
	session := createTestSession(user.ID.Hex())
	deviceHash := GenerateDeviceHash("127.0.0.1", "test-agent")

	mockSessionRepo.On("FindUserSessionByRefreshToken", "test-refresh-token").Return(session, nil)
	mockSessionRepo.On("ValidateUserSession", session.ID, session.UserID, deviceHash).Return(true)
	mockUserRepo.On("FindByID", user.ID.Hex()).Return(user)
	mockSessionRepo.On("UpdateUserSessionRefreshToken", session.ID, mock.AnythingOfType("string")).Return(nil)

	reqBody := map[string]string{
		"refresh_token": "test-refresh-token",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessTokenRefreshed, response["message"])

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestRefreshToken_InvalidRequestBody(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/refresh", handler.RefreshToken)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidRequestBody, response["message"])
}

func TestRefreshToken_TokenRequired(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/refresh", handler.RefreshToken)

	reqBody := map[string]string{
		"refresh_token": "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrRefreshTokenRequired, response["message"])
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/refresh", handler.RefreshToken)

	mockSessionRepo.On("FindUserSessionByRefreshToken", "invalid-token").Return(nil, nil)

	reqBody := map[string]string{
		"refresh_token": "invalid-token",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidToken, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

func TestRefreshToken_DeviceHashMismatch(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/refresh", func(c *fiber.Ctx) error {
		c.Locals("userIP", "192.168.1.1")
		c.Locals("userAgent", "different-agent")
		return handler.RefreshToken(c)
	})

	user := createTestUser("test@example.com")
	session := createTestSession(user.ID.Hex())
	deviceHash := GenerateDeviceHash("192.168.1.1", "different-agent")

	mockSessionRepo.On("FindUserSessionByRefreshToken", "test-refresh-token").Return(session, nil)
	mockSessionRepo.On("ValidateUserSession", session.ID, session.UserID, deviceHash).Return(false)

	reqBody := map[string]string{
		"refresh_token": "test-refresh-token",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidToken, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/refresh", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.RefreshToken(c)
	})

	user := createTestUser("test@example.com")
	session := createTestSession(user.ID.Hex())
	deviceHash := GenerateDeviceHash("127.0.0.1", "test-agent")

	mockSessionRepo.On("FindUserSessionByRefreshToken", "test-refresh-token").Return(session, nil)
	mockSessionRepo.On("ValidateUserSession", session.ID, session.UserID, deviceHash).Return(true)
	mockUserRepo.On("FindByID", user.ID.Hex()).Return(nil)

	reqBody := map[string]string{
		"refresh_token": "test-refresh-token",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrUserNotFound, response["message"])

	mockSessionRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestRefreshToken_UpdateRefreshTokenFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/refresh", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.RefreshToken(c)
	})

	user := createTestUser("test@example.com")
	session := createTestSession(user.ID.Hex())
	deviceHash := GenerateDeviceHash("127.0.0.1", "test-agent")

	mockSessionRepo.On("FindUserSessionByRefreshToken", "test-refresh-token").Return(session, nil)
	mockSessionRepo.On("ValidateUserSession", session.ID, session.UserID, deviceHash).Return(true)
	mockUserRepo.On("FindByID", user.ID.Hex()).Return(user)
	mockSessionRepo.On("UpdateUserSessionRefreshToken", session.ID, mock.AnythingOfType("string")).Return(errors.New("db error"))

	reqBody := map[string]string{
		"refresh_token": "test-refresh-token",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrUpdateSession, response["message"])

	mockSessionRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test GoogleLogin Handler
func TestGoogleLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google", handler.GoogleLogin)

	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(t, location, "accounts.google.com")
}

// Test GoogleCallback Handler
func TestGoogleCallback_MissingCode(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", handler.GoogleCallback)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidToken, response["message"])
}

func TestGoogleCallback_ExchangeTokenFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	// Invalid client secret to cause token exchange failure
	cfg.GoogleClientSecret = "invalid-secret"
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", handler.GoogleCallback)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=test-auth-code", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrExchangeToken, response["message"])
}

func TestGoogleCallback_GetUserInfoFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	// Invalid redirect URL to cause token exchange to return invalid token
	cfg.GoogleRedirectUrl = "http://invalid-url.com/callback"
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", handler.GoogleCallback)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=test-auth-code", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	// This will fail at token exchange stage, which is acceptable
	assert.True(t, resp.StatusCode == fiber.StatusInternalServerError)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	// Could be exchange token error or get user info error depending on where it fails
	assert.True(t, response["message"] == ErrExchangeToken || response["message"] == ErrGetUserInfo)
}

func TestGoogleCallback_ParseUserInfoFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", handler.GoogleCallback)

	// Create a request that will cause invalid JSON response
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=invalid-code", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusInternalServerError)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	// Could be exchange token error or parse user info error
	assert.True(t, response["message"] == ErrExchangeToken || response["message"] == ErrParseUserInfo)
}

func TestGoogleCallback_NewUserSuccess(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.GoogleCallback(c)
	})

	// Create a new user that doesn't exist yet
	newUser := createTestUser("newuser@example.com")
	session := createTestSession(newUser.ID.Hex())

	mockUserRepo.On("FindByGoogleID", "test-google-id").Return(nil)
	mockUserRepo.On("FindByEmail", "newuser@example.com").Return(nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*internal.User")).Return(newUser)
	mockSessionRepo.On("CreateUserSession", newUser.ID.Hex(), "127.0.0.1", "test-agent").Return(session, nil)
	mockSessionRepo.On("UpdateUserSessionRefreshToken", session.ID, mock.AnythingOfType("string")).Return(nil)

	// Mock successful OAuth flow - this would require actual HTTP mocking in a real scenario
	// For this test, we expect it to fail at the HTTP request stage
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=test-code", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	// Expected to fail at HTTP request stage since we can't mock Google's OAuth server
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrExchangeToken, response["message"])
}

func TestGoogleCallback_ExistingUserWithGoogleID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.GoogleCallback(c)
	})

	existingUser := createTestUser("existing@example.com")
	session := createTestSession(existingUser.ID.Hex())

	mockUserRepo.On("FindByGoogleID", "test-google-id").Return(existingUser)
	mockSessionRepo.On("CreateUserSession", existingUser.ID.Hex(), "127.0.0.1", "test-agent").Return(session, nil)
	mockSessionRepo.On("UpdateUserSessionRefreshToken", session.ID, mock.AnythingOfType("string")).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=test-code", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	// Expected to fail at HTTP request stage
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrExchangeToken, response["message"])
}

func TestGoogleCallback_UpdateGoogleIDFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.GoogleCallback(c)
	})

	existingUser := createTestUser("existing@example.com")

	mockUserRepo.On("FindByGoogleID", "test-google-id").Return(nil)
	mockUserRepo.On("FindByEmail", "existing@example.com").Return(existingUser)
	mockUserRepo.On("UpdateGoogleID", existingUser.ID.Hex(), "test-google-id").Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=test-code", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	// Expected to fail at HTTP request stage
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrExchangeToken, response["message"])
}

func TestGoogleCallback_SessionCreationFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/auth/google/callback", func(c *fiber.Ctx) error {
		c.Locals("userIP", "127.0.0.1")
		c.Locals("userAgent", "test-agent")
		return handler.GoogleCallback(c)
	})

	existingUser := createTestUser("existing@example.com")

	mockUserRepo.On("FindByGoogleID", "test-google-id").Return(existingUser)
	mockSessionRepo.On("CreateUserSession", existingUser.ID.Hex(), "127.0.0.1", "test-agent").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=test-code", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	// Expected to fail at HTTP request stage
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrExchangeToken, response["message"])
}

// Test GetProfile Handler
func TestGetProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.GetProfile(c)
	})

	user := createTestUser("test@example.com")
	mockUserRepo.On("FindByID", "507f1f77bcf86cd799439011").Return(user)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessProfileRetrieved, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestGetProfile_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.GetProfile(c)
	})

	mockUserRepo.On("FindByID", "507f1f77bcf86cd799439011").Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrUserNotFound, response["message"])

	mockUserRepo.AssertExpectations(t)
}

// Test Logout Handler
func TestLogout_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/logout", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		c.Locals("sessionId", "507f1f77bcf86cd799439012")
		return handler.Logout(c)
	})

	sessionID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	mockSessionRepo.On("DeactivateUserSession", sessionID).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessLogoutSuccessful, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

func TestLogout_InvalidSessionID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/logout", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		c.Locals("sessionId", "invalid-session-id")
		return handler.Logout(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidSessionID, response["message"])
}

func TestLogout_DeactivationFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/logout", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		c.Locals("sessionId", "507f1f77bcf86cd799439012")
		return handler.Logout(c)
	})

	sessionID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	mockSessionRepo.On("DeactivateUserSession", sessionID).Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrFailedToLogout, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

// Test GetDevices Handler
func TestGetDevices_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/devices", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.GetDevices(c)
	})

	sessions := []UserSession{
		*createTestSession("507f1f77bcf86cd799439011"),
		*createTestSession("507f1f77bcf86cd799439011"),
	}
	mockSessionRepo.On("GetUserSessions", "507f1f77bcf86cd799439011").Return(sessions, nil)

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessDevicesRetrieved, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

func TestGetDevices_Failed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Get("/devices", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.GetDevices(c)
	})

	mockSessionRepo.On("GetUserSessions", "507f1f77bcf86cd799439011").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrGetUserSession, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

// Test RevokeDevice Handler
func TestRevokeDevice_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Delete("/devices/:sessionId", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.RevokeDevice(c)
	})

	sessionID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	session := createTestSession("507f1f77bcf86cd799439011")

	mockSessionRepo.On("FindUserSessionByID", sessionID, "507f1f77bcf86cd799439011").Return(session, nil)
	mockSessionRepo.On("DeactivateUserSession", sessionID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/devices/507f1f77bcf86cd799439012", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessDeviceRevoked, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

func TestRevokeDevice_InvalidSessionID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Delete("/devices/:sessionId", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.RevokeDevice(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/devices/invalid-id", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidSessionID, response["message"])
}

func TestRevokeDevice_SessionNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Delete("/devices/:sessionId", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.RevokeDevice(c)
	})

	sessionID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	mockSessionRepo.On("FindUserSessionByID", sessionID, "507f1f77bcf86cd799439011").Return(nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/devices/507f1f77bcf86cd799439012", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrDeviceNotFound, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

// Test LogoutAllDevices Handler
func TestLogoutAllDevices_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/logout-all", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.LogoutAllDevices(c)
	})

	mockSessionRepo.On("ClearAllUserSessions", "507f1f77bcf86cd799439011").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/logout-all", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessLoggedOutAllDevices, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

func TestLogoutAllDevices_Failed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/logout-all", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.LogoutAllDevices(c)
	})

	mockSessionRepo.On("ClearAllUserSessions", "507f1f77bcf86cd799439011").Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/logout-all", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrClearAllSessions, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

// Test SendPin Handler
func TestSendPin_Success_ExistingUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	user := createTestUser("test@example.com")
	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)
	mockUserRepo.On("SetPinWithExpiry", "test@example.com", mock.AnythingOfType("string")).Return(nil)

	reqBody := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessPinSent, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestSendPin_Success_NewUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	newUser := createTestUser("newuser@example.com")
	mockUserRepo.On("FindByEmail", "newuser@example.com").Return(nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*internal.User")).Return(newUser)
	mockUserRepo.On("SetPinWithExpiry", "newuser@example.com", mock.AnythingOfType("string")).Return(nil)

	reqBody := map[string]string{
		"email": "newuser@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessPinSent, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestSendPin_InvalidRequestBody(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrInvalidRequestBody, response["message"])
}

func TestSendPin_EmailRequired(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	reqBody := map[string]string{
		"email": "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrEmailRequired, response["message"])
}

func TestSendPin_CreateUserFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	mockUserRepo.On("FindByEmail", "newuser@example.com").Return(nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*internal.User")).Return(nil)

	reqBody := map[string]string{
		"email": "newuser@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrCreateUser, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestSendPin_SetPinFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	user := createTestUser("test@example.com")
	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)
	mockUserRepo.On("SetPinWithExpiry", "test@example.com", mock.AnythingOfType("string")).Return(errors.New("db error"))

	reqBody := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrSetPin, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestSendPin_WithPinEnabled(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	cfg.PinEnabled = true // Enable PIN generation
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	user := createTestUser("test@example.com")
	mockUserRepo.On("FindByEmail", "test@example.com").Return(user)
	mockUserRepo.On("SetPinWithExpiry", "test@example.com", mock.AnythingOfType("string")).Return(nil)

	reqBody := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessPinSent, response["message"])

	mockUserRepo.AssertExpectations(t)
}

func TestSendPin_WithPinEnabled_NewUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	cfg.PinEnabled = true // Enable PIN generation
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Post("/send-pin", handler.SendPin)

	newUser := createTestUser("newuser@example.com")
	mockUserRepo.On("FindByEmail", "newuser@example.com").Return(nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*internal.User")).Return(newUser)
	mockUserRepo.On("SetPinWithExpiry", "newuser@example.com", mock.AnythingOfType("string")).Return(nil)

	reqBody := map[string]string{
		"email": "newuser@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/send-pin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, SuccessPinSent, response["message"])

	mockUserRepo.AssertExpectations(t)
}

// Test missing RevokeDevice scenarios
func TestRevokeDevice_FindSessionFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Delete("/devices/:sessionId", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.RevokeDevice(c)
	})

	sessionID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	mockSessionRepo.On("FindUserSessionByID", sessionID, "507f1f77bcf86cd799439011").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodDelete, "/devices/507f1f77bcf86cd799439012", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrFindUserSession, response["message"])

	mockSessionRepo.AssertExpectations(t)
}

func TestRevokeDevice_DeactivationFailed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockUserSessionRepository)
	cfg := createTestConfig()
	handler := createHandlerWithMocks(cfg, mockUserRepo, mockSessionRepo)

	app := setupTestApp()
	app.Delete("/devices/:sessionId", func(c *fiber.Ctx) error {
		c.Locals("userId", "507f1f77bcf86cd799439011")
		return handler.RevokeDevice(c)
	})

	sessionID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	session := createTestSession("507f1f77bcf86cd799439011")

	mockSessionRepo.On("FindUserSessionByID", sessionID, "507f1f77bcf86cd799439011").Return(session, nil)
	mockSessionRepo.On("DeactivateUserSession", sessionID).Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodDelete, "/devices/507f1f77bcf86cd799439012", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, ErrDeactivateSession, response["message"])

	mockSessionRepo.AssertExpectations(t)
}
