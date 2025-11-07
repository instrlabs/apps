package internal

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
	AddRefreshToken(userID, token string) error
	RemoveRefreshToken(userID, token string) error
	ValidateRefreshToken(userID, token string) bool
	ClearAllRefreshTokens(userID string) error
}
