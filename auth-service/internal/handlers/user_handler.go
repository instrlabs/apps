package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/instrlabs/auth-service/internal/services"
)

// UserHandler handles user endpoints
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile retrieves user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from middleware
	userID, _ := c.Locals("userID").(string)
	if userID == "" {
		return h.sendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	// Call service
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Send response
	return h.sendSuccessResponse(c, "Profile retrieved successfully", fiber.Map{"user": user})
}

// sendErrorResponse sends a standardized error response
func (h *UserHandler) sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    nil,
	})
}

// sendSuccessResponse sends a standardized success response
func (h *UserHandler) sendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    data,
	})
}
