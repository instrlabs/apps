package internal

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(user *User) error
	FindByEmail(email string, user *User) error
	FindByID(id string, user *User) error
	FindByGoogleID(googleID string, user *User) error
	UpdateGoogleID(userID string, googleID string) error
	SetPinWithExpiry(email, hashedPin string) error
	ClearPin(userID string) error
	SetRegisteredAt(userID string) error
	AddRefreshToken(userID, token string) error
	RemoveRefreshToken(userID, token string) error
	ValidateRefreshToken(userID, token string) error
	ClearAllRefreshTokens(userID string) error
}
