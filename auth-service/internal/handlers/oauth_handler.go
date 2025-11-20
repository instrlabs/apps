package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/instrlabs/auth-service/internal/services"
)

// OAuthHandler handles OAuth endpoints
type OAuthHandler struct {
	oauthService *services.OAuthService
	authService  *services.AuthService
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthService *services.OAuthService, authService *services.AuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		authService:  authService,
	}
}

// GoogleLogin initiates Google OAuth login
func (h *OAuthHandler) GoogleLogin(c *fiber.Ctx) error {
	// Call service
	url, err := h.oauthService.InitiateGoogleLogin()
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, "Failed to initiate Google login")
	}

	// Redirect to Google
	return c.Redirect(url, fiber.StatusFound)
}

// GoogleCallback handles Google OAuth callback
func (h *OAuthHandler) GoogleCallback(c *fiber.Ctx) error {
	// Get authorization code
	code := c.Query("code")
	if code == "" {
		return h.sendErrorResponse(c, fiber.StatusBadRequest, "Missing authorization code")
	}

	// Call service
	result, err := h.oauthService.HandleGoogleCallback(code)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Build redirect URL with tokens
	redirectURL := h.oauthService.BuildGoogleRedirectURL(
		result.AccessToken,
		result.RefreshToken,
		result.ExpiresIn,
	)

	// Redirect with tokens
	return c.Redirect(redirectURL, fiber.StatusFound)
}

// sendErrorResponse sends a standardized error response
func (h *OAuthHandler) sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    nil,
	})
}

// sendSuccessResponse sends a standardized success response
func (h *OAuthHandler) sendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    data,
	})
}
