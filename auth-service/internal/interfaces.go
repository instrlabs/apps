package internal

import "time"

type IUserRepository interface {
	Create(user *User) *User
	FindByEmail(email string) *User
	FindByID(id string) *User
	UpdateRefreshToken(userID string, token string) error
	UpdateRefreshTokenWithExpiry(userID string, token string, duration time.Duration) error
	ClearRefreshToken(userID string) error
	SetPinWithExpiry(email string, hashedPin string) error
	ClearPin(userID string) error
	SetRegisteredAt(userID string) error
	UpdateGoogleID(userID string, googleID string) error
	FindByGoogleID(googleID string) *User
	FindByRefreshToken(token string) *User
}

type ISessionRepository interface {
	CreateSession(userID string, ipAddress string, userAgent string) (*UserSession, error)
	FindSessionByRefreshToken(token string) (*UserSession, error)
	FindSessionByID(sessionID string, userID string) (*UserSession, error)
	ValidateSession(sessionID string, userID string, deviceHash string) bool
	UpdateSessionActivity(sessionID string) error
	DeactivateSession(sessionID string) error
	GetUserSessions(userID string) ([]UserSession, error)
	ClearExpiredSessions(userID string) error
	ClearAllUserSessions(userID string) error
	UpdateSessionRefreshToken(sessionID string, refreshToken string) error
}
