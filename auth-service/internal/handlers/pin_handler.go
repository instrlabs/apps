package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/instrlabs/auth-service/internal/services"
	"github.com/instrlabs/auth-service/internal/validators"
)

// PinHandler handles PIN-related endpoints
type PinHandler struct {
	pinService *services.PinService
	validator  *validators.RequestValidator
}

// NewPinHandler creates a new PIN handler
func NewPinHandler(pinService *services.PinService, validator *validators.RequestValidator) *PinHandler {
	return &PinHandler{
		pinService: pinService,
		validator:  validator,
	}
}

// SendPin handles sending PIN to user email
func (h *PinHandler) SendPin(c *fiber.Ctx) error {
	// Parse and validate request
	req, err := h.validator.ParsePinRequest(c)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Call service
	err = h.pinService.GenerateAndSendPIN(req.Email)
	if err != nil {
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Send response
	return h.sendSuccessResponse(c, "PIN sent successfully", nil)
}

// sendErrorResponse sends a standardized error response
func (h *PinHandler) sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    nil,
	})
}

// sendSuccessResponse sends a standardized success response
func (h *PinHandler) sendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    data,
	})
}
