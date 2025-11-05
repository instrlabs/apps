package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(user *User) *User
	FindByEmail(email string) *User
	FindByID(id string) *User
	FindByGoogleID(googleID string) *User
	UpdateGoogleID(userID string, googleID string) error
	SetPinWithExpiry(email, hashedPin string) error
	ClearPin(userID string) error
	SetRegisteredAt(userID string) error
}

// UserSessionRepositoryInterface defines the interface for session repository operations
type UserSessionRepositoryInterface interface {
	CreateUserSession(userID, ipAddress, userAgent string) (*UserSession, error)
	FindUserSessionByRefreshToken(refreshToken string) (*UserSession, error)
	FindUserSessionByID(id primitive.ObjectID, userID string) (*UserSession, error)
	ValidateUserSession(id primitive.ObjectID, userID, deviceHash string) bool
	UpdateUserSessionActivity(id primitive.ObjectID) error
	DeactivateUserSession(id primitive.ObjectID) error
	GetUserSessions(userID string) ([]UserSession, error)
	ClearExpiredUserSessions(userID string) error
	ClearAllUserSessions(userID string) error
	UpdateUserSessionRefreshToken(id primitive.ObjectID, newRefreshToken string) error
}
