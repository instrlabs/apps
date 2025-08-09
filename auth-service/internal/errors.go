package internal

// Error messages
const (
	// General errors
	ErrInternalServer     = "Internal server error. Please contact support."
	ErrInvalidRequestBody = "Invalid request body"

	// Authentication errors
	ErrInvalidCredentials = "Invalid email or password"
	ErrUnauthorized       = "Unauthorized access"
	ErrTokenExpired       = "Token has expired"
	ErrInvalidToken       = "Invalid token"

	// Validation errors
	ErrEmailRequired        = "Email is required"
	ErrInvalidEmail         = "Invalid email format"
	ErrPasswordRequired     = "Password is required"
	ErrRefreshTokenRequired = "Refresh token is required"

	// User errors
	ErrUserNotFound       = "User not found"
	ErrEmailAlreadyExists = "Email already exists"
)
