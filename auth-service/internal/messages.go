package internal

// Error messages
const (
	// General errors
	ErrInvalidRequestBody = "Invalid request body"
	ErrInternalServer     = "Internal server error"

	// Authentication errors
	ErrInvalidCredentials = "Invalid email or pin"
	ErrInvalidToken       = "Invalid token"

	// Validation errors
	ErrEmailRequired        = "Email is required"
	ErrPasswordRequired     = "Pin is required"
	ErrRefreshTokenRequired = "Refresh token is required"

	// User errors
	ErrUserNotFound = "User not found"

	// Session errors
	ErrInvalidSessionID  = "Invalid session ID"
	ErrFailedToLogout    = "Failed to logout user"
	ErrDeviceNotFound    = "Device not found"
	ErrGetUserSession    = "Failed to get user session"
	ErrFindUserSession   = "Failed to find user session"
	ErrDeactivateSession = "Failed to deactivate session"
	ErrClearAllSessions  = "Failed to clear all user sessions"

	// Token and session management errors
	ErrGenerateAccessToken  = "Error generate access token"
	ErrGenerateRefreshToken = "Error generate refresh token"
	ErrCreateSession        = "Error create user session"
	ErrUpdateSession        = "Error update user session"

	// OAuth errors
	ErrExchangeToken    = "Error exchange authorization code"
	ErrGetUserInfo      = "Error retrieve user information"
	ErrParseUserInfo    = "Error parse user information"
	ErrCreateGoogleUser = "Error create Google user"
	ErrUpdateGoogleID   = "Error link Google account"

	// PIN errors
	ErrCreateUser = "Error create user account"
	ErrSetPin     = "Error set login PIN"
)

// Success messages
const (
	SuccessLogin               = "Login successful"
	SuccessTokenRefreshed      = "Token refreshed successfully"
	SuccessProfileRetrieved    = "Profile retrieved successfully"
	SuccessLogoutSuccessful    = "Logout successful"
	SuccessDevicesRetrieved    = "Devices retrieved successfully"
	SuccessDeviceRevoked       = "Device revoked successfully"
	SuccessLoggedOutAllDevices = "Logged out from all devices"
	SuccessPinSent             = "PIN sent"
)
