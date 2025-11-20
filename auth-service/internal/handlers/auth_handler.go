package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/instrlabs/auth-service/internal/services"
	"github.com/instrlabs/auth-service/internal/validators"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
	validator   *validators.RequestValidator
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService, validator *validators.RequestValidator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse and validate request
	req, err := h.validator.ParseLoginRequest(c)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Call service
	result, err := h.authService.Login(req.Email, req.Pin)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Send response
	return h.sendSuccessResponse(c, "Login successful", result)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// Get user ID from middleware
	userID, _ := c.Locals("userId").(string)
	if userID == "" {
		return h.sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	// Parse and validate request
	req, err := h.validator.ParseRefreshTokenRequest(c)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Call service
	result, err := h.authService.RefreshToken(userID, req.RefreshToken)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	// Send response
	return h.sendSuccessResponse(c, "Token refreshed successfully", result)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get user ID from middleware
	userID, _ := c.Locals("userId").(string)
	if userID == "" {
		return h.sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	// Parse optional refresh token
	req, err := h.validator.ParseRefreshTokenRequest(c)
	if err != nil {
		// Continue with empty refresh token if parsing fails
		req = &validators.RefreshTokenRequest{}
	}

	// Call service
	err = h.authService.Logout(userID, req.RefreshToken)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout")
	}

	// Send response
	return h.sendSuccessResponse(c, "Logout successful", nil)
}

// sendErrorResponse sends a standardized error response
func (h *AuthHandler) sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    nil,
	})
}

// sendSuccessResponse sends a standardized success response
func (h *AuthHandler) sendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    data,
	})
}
