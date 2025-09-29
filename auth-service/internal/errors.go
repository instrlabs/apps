package internal

// Error messages
const (
	// General errors
	ErrInternalServer     = "Internal server error. Please contact support."
	ErrInvalidRequestBody = "Invalid request body"

	// Authentication errors
	ErrInvalidCredentials = "Invalid email or pin"
	ErrInvalidToken       = "Invalid token"

	// Validation errors
	ErrEmailRequired        = "Email is required"
	ErrPasswordRequired     = "Pin is required"
	ErrRefreshTokenRequired = "Refresh token is required"

	// User errors
	ErrUserNotFound = "User not found"
)
