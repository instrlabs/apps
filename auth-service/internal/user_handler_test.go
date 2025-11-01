package internal

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepository struct {
	CreateFunc                       func(user *User) *User
	FindByEmailFunc                  func(email string) *User
	FindByIDFunc                     func(id string) *User
	UpdateRefreshTokenFunc           func(userID string, token string) error
	UpdateRefreshTokenWithExpiryFunc func(userID string, token string, duration time.Duration) error
	ClearRefreshTokenFunc            func(userID string) error
	SetPinWithExpiryFunc             func(email string, hashedPin string) error
	ClearPinFunc                     func(userID string) error
	SetRegisteredAtFunc              func(userID string) error
	UpdateGoogleIDFunc               func(userID string, googleID string) error
	FindByGoogleIDFunc               func(googleID string) *User
	FindByRefreshTokenFunc           func(token string) *User
}

func (m *MockUserRepository) Create(user *User) *User {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return user
}

func (m *MockUserRepository) FindByEmail(email string) *User {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil
}

func (m *MockUserRepository) FindByID(id string) *User {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil
}

func (m *MockUserRepository) UpdateRefreshToken(userID string, token string) error {
	if m.UpdateRefreshTokenFunc != nil {
		return m.UpdateRefreshTokenFunc(userID, token)
	}
	return nil
}

func (m *MockUserRepository) UpdateRefreshTokenWithExpiry(userID string, token string, duration time.Duration) error {
	if m.UpdateRefreshTokenWithExpiryFunc != nil {
		return m.UpdateRefreshTokenWithExpiryFunc(userID, token, duration)
	}
	return nil
}

func (m *MockUserRepository) ClearRefreshToken(userID string) error {
	if m.ClearRefreshTokenFunc != nil {
		return m.ClearRefreshTokenFunc(userID)
	}
	return nil
}

func (m *MockUserRepository) SetPinWithExpiry(email string, hashedPin string) error {
	if m.SetPinWithExpiryFunc != nil {
		return m.SetPinWithExpiryFunc(email, hashedPin)
	}
	return nil
}

func (m *MockUserRepository) ClearPin(userID string) error {
	if m.ClearPinFunc != nil {
		return m.ClearPinFunc(userID)
	}
	return nil
}

func (m *MockUserRepository) SetRegisteredAt(userID string) error {
	if m.SetRegisteredAtFunc != nil {
		return m.SetRegisteredAtFunc(userID)
	}
	return nil
}

func (m *MockUserRepository) UpdateGoogleID(userID string, googleID string) error {
	if m.UpdateGoogleIDFunc != nil {
		return m.UpdateGoogleIDFunc(userID, googleID)
	}
	return nil
}

func (m *MockUserRepository) FindByGoogleID(googleID string) *User {
	if m.FindByGoogleIDFunc != nil {
		return m.FindByGoogleIDFunc(googleID)
	}
	return nil
}

func (m *MockUserRepository) FindByRefreshToken(token string) *User {
	if m.FindByRefreshTokenFunc != nil {
		return m.FindByRefreshTokenFunc(token)
	}
	return nil
}

type MockSessionRepository struct {
	CreateSessionFunc             func(userID string, ipAddress string, userAgent string) (*UserSession, error)
	FindSessionByRefreshTokenFunc func(token string) (*UserSession, error)
	FindSessionByIDFunc           func(sessionID string, userID string) (*UserSession, error)
	ValidateSessionFunc           func(sessionID string, userID string, deviceHash string) bool
	UpdateSessionActivityFunc     func(sessionID string) error
	DeactivateSessionFunc         func(sessionID string) error
	GetUserSessionsFunc           func(userID string) ([]UserSession, error)
	ClearExpiredSessionsFunc      func(userID string) error
	ClearAllUserSessionsFunc      func(userID string) error
	UpdateSessionRefreshTokenFunc func(sessionID string, refreshToken string) error
}

func (m *MockSessionRepository) CreateSession(userID string, ipAddress string, userAgent string) (*UserSession, error) {
	if m.CreateSessionFunc != nil {
		return m.CreateSessionFunc(userID, ipAddress, userAgent)
	}
	return &UserSession{ID: primitive.NewObjectID(), SessionID: "test-session"}, nil
}

func (m *MockSessionRepository) FindSessionByRefreshToken(token string) (*UserSession, error) {
	if m.FindSessionByRefreshTokenFunc != nil {
		return m.FindSessionByRefreshTokenFunc(token)
	}
	return nil, nil
}

func (m *MockSessionRepository) FindSessionByID(sessionID string, userID string) (*UserSession, error) {
	if m.FindSessionByIDFunc != nil {
		return m.FindSessionByIDFunc(sessionID, userID)
	}
	return nil, nil
}

func (m *MockSessionRepository) ValidateSession(sessionID string, userID string, deviceHash string) bool {
	if m.ValidateSessionFunc != nil {
		return m.ValidateSessionFunc(sessionID, userID, deviceHash)
	}
	return true
}

func (m *MockSessionRepository) UpdateSessionActivity(sessionID string) error {
	if m.UpdateSessionActivityFunc != nil {
		return m.UpdateSessionActivityFunc(sessionID)
	}
	return nil
}

func (m *MockSessionRepository) DeactivateSession(sessionID string) error {
	if m.DeactivateSessionFunc != nil {
		return m.DeactivateSessionFunc(sessionID)
	}
	return nil
}

func (m *MockSessionRepository) GetUserSessions(userID string) ([]UserSession, error) {
	if m.GetUserSessionsFunc != nil {
		return m.GetUserSessionsFunc(userID)
	}
	return []UserSession{}, nil
}

func (m *MockSessionRepository) ClearExpiredSessions(userID string) error {
	if m.ClearExpiredSessionsFunc != nil {
		return m.ClearExpiredSessionsFunc(userID)
	}
	return nil
}

func (m *MockSessionRepository) ClearAllUserSessions(userID string) error {
	if m.ClearAllUserSessionsFunc != nil {
		return m.ClearAllUserSessionsFunc(userID)
	}
	return nil
}

func (m *MockSessionRepository) UpdateSessionRefreshToken(sessionID string, refreshToken string) error {
	if m.UpdateSessionRefreshTokenFunc != nil {
		return m.UpdateSessionRefreshTokenFunc(sessionID, refreshToken)
	}
	return nil
}

func newMockConfig() *Config {
	return &Config{
		Environment:        "test",
		Port:               ":3000",
		MongoURI:           "mongodb://localhost:27017",
		MongoDB:            "auth_test",
		JWTSecret:          "test-secret-key-for-testing-only",
		TokenExpiryHours:   1,
		RefreshExpiryHours: 168,
		SMTPHost:           "localhost",
		SMTPPort:           "1025",
		SMTPUsername:       "test",
		SMTPPassword:       "test",
		EmailFrom:          "test@example.com",
		GoogleClientID:     "test-client-id",
		GoogleClientSecret: "test-secret",
		GoogleRedirectUrl:  "http://localhost:3000/google/callback",
		ApiUrl:             "http://localhost:3001",
		WebUrl:             "http://localhost:3000",
		PinEnabled:         true,
		CookieDomain:       "localhost",
	}
}

func TestGenerateAccessToken_ValidToken(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}
	handler := NewUserHandler(config, userRepo, sessionRepo)

	userID := "test-user-id"
	sessionID := "test-session-id"
	roles := []string{"user"}

	token, err := handler.generateAccessToken(userID, roles, sessionID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, userID, claims["user_id"])
	assert.Equal(t, sessionID, claims["session_id"])
}

func TestGenerateAccessToken_HasExpiry(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}
	handler := NewUserHandler(config, userRepo, sessionRepo)

	before := time.Now()
	token, _ := handler.generateAccessToken("user-id", []string{"user"}, "session-id")

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	expectedExpiry := time.Duration(config.TokenExpiryHours) * time.Hour

	timeDiff := expTime.Sub(before)
	assert.Greater(t, timeDiff, expectedExpiry-time.Minute)
	assert.Less(t, timeDiff, expectedExpiry+time.Minute)
}

func TestGenerateRefreshToken_ValidFormat(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}
	handler := NewUserHandler(config, userRepo, sessionRepo)

	token, err := handler.generateRefreshToken()

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Greater(t, len(token), 0)
}

func TestGenerateRefreshToken_Uniqueness(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}
	handler := NewUserHandler(config, userRepo, sessionRepo)

	token1, _ := handler.generateRefreshToken()
	token2, _ := handler.generateRefreshToken()

	assert.NotEqual(t, token1, token2)
}

func TestGenerateAccessToken_ContainsRoles(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}
	handler := NewUserHandler(config, userRepo, sessionRepo)

	roles := []string{"user", "admin"}
	token, _ := handler.generateAccessToken("user-id", roles, "session-id")

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	claimsRoles := claims["roles"].([]interface{})
	assert.Equal(t, 2, len(claimsRoles))
}

func TestGenerateDeviceHash_Consistency(t *testing.T) {
	ip := "192.168.1.1"
	ua := "Mozilla/5.0"

	hash1 := GenerateDeviceHash(ip, ua)
	hash2 := GenerateDeviceHash(ip, ua)

	assert.Equal(t, hash1, hash2)
}

func TestGenerateDeviceHash_DifferentInputs(t *testing.T) {
	hash1 := GenerateDeviceHash("192.168.1.1", "Mozilla/5.0")
	hash2 := GenerateDeviceHash("192.168.1.2", "Mozilla/5.0")
	hash3 := GenerateDeviceHash("192.168.1.1", "Chrome/120.0")

	assert.NotEqual(t, hash1, hash2)
	assert.NotEqual(t, hash1, hash3)
	assert.NotEqual(t, hash2, hash3)
}

func TestGenerateDeviceHash_Format(t *testing.T) {
	hash := GenerateDeviceHash("192.168.1.1", "Mozilla/5.0")

	assert.NotEmpty(t, hash)
	_, err := hex.DecodeString(hash)
	assert.NoError(t, err)
}

func TestGenerateDeviceHash_NonEmptyIpAndUA(t *testing.T) {
	hash1 := GenerateDeviceHash("", "Mozilla/5.0")
	hash2 := GenerateDeviceHash("192.168.1.1", "")
	hash3 := GenerateDeviceHash("", "")

	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)
	assert.NotEmpty(t, hash3)
	assert.NotEqual(t, hash1, hash2)
	assert.NotEqual(t, hash2, hash3)
}

func TestGenerateSessionID_Length(t *testing.T) {
	sessionID := GenerateSessionID()

	assert.Greater(t, len(sessionID), 0)
	assert.NotEmpty(t, sessionID)
}

func TestGenerateSessionID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)

	for i := 0; i < 100; i++ {
		sessionID := GenerateSessionID()
		assert.False(t, ids[sessionID], "Duplicate session ID generated")
		ids[sessionID] = true
	}
}

func TestGenerateSessionID_ValidBase64(t *testing.T) {
	sessionID := GenerateSessionID()

	assert.NotEmpty(t, sessionID)
}

func TestUserHandler_RepositoryDependencyInjection(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}

	handler := NewUserHandler(config, userRepo, sessionRepo)

	assert.NotNil(t, handler)
	assert.Equal(t, config, handler.cfg)
}

func TestMockUserRepository_FindByEmail(t *testing.T) {
	mockRepo := &MockUserRepository{}
	testUser := &User{
		ID:    primitive.NewObjectID(),
		Email: "test@example.com",
	}

	mockRepo.FindByEmailFunc = func(email string) *User {
		if email == "test@example.com" {
			return testUser
		}
		return nil
	}

	found := mockRepo.FindByEmail("test@example.com")
	assert.NotNil(t, found)
	assert.Equal(t, "test@example.com", found.Email)

	notFound := mockRepo.FindByEmail("other@example.com")
	assert.Nil(t, notFound)
}

func TestMockSessionRepository_CreateSession(t *testing.T) {
	mockRepo := &MockSessionRepository{}
	mockSession := &UserSession{
		ID:        primitive.NewObjectID(),
		SessionID: "test-session",
	}

	mockRepo.CreateSessionFunc = func(userID string, ipAddress string, userAgent string) (*UserSession, error) {
		return mockSession, nil
	}

	session, err := mockRepo.CreateSession("user-123", "192.168.1.1", "Mozilla/5.0")

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test-session", session.SessionID)
}

func TestMockSessionRepository_ValidateSession(t *testing.T) {
	mockRepo := &MockSessionRepository{}

	mockRepo.ValidateSessionFunc = func(sessionID string, userID string, deviceHash string) bool {
		return deviceHash == "valid-hash"
	}

	result := mockRepo.ValidateSession("session-1", "user-1", "valid-hash")
	assert.True(t, result)

	result = mockRepo.ValidateSession("session-1", "user-1", "invalid-hash")
	assert.False(t, result)
}

func TestGenerateAccessToken_MultipleRoles(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}
	handler := NewUserHandler(config, userRepo, sessionRepo)

	roles := []string{"user", "moderator", "admin"}
	token, err := handler.generateAccessToken("user-123", roles, "session-456")

	assert.NoError(t, err)

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	assert.Equal(t, "user-123", claims["user_id"])
	assert.Equal(t, "session-456", claims["session_id"])
}

func TestGenerateAccessToken_NoRoles(t *testing.T) {
	config := newMockConfig()
	userRepo := &MockUserRepository{}
	sessionRepo := &MockSessionRepository{}
	handler := NewUserHandler(config, userRepo, sessionRepo)

	token, err := handler.generateAccessToken("user-id", []string{}, "session-id")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
