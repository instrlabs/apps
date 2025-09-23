package internal

// Error messages
const (
	// General errors
	ErrInternalServer     = "Internal server error. Please contact support."
	ErrInvalidRequestBody = "Invalid request body"

	// Authentication errors
	ErrInvalidCredentials = "Invalid email or pin"
	ErrUnauthorized       = "Unauthorized access"
	ErrTokenExpired       = "Token has expired"
	ErrInvalidToken       = "Invalid token"

	// Validation errors
	ErrEmailRequired        = "Email is required"
	ErrEmailAlreadyInUse    = "Email already in use"
	ErrEmailNotFound        = "Email not found"
	ErrInvalidEmail         = "Invalid email format"
	ErrPasswordRequired     = "Pin is required"
	ErrRefreshTokenRequired = "Refresh token is required"

	// User errors
	ErrUserNotFound       = "User not found"
	ErrEmailAlreadyExists = "Email already exists"
)
