package helpers

import (
	"github.com/gofiber/fiber/v2"
)

// ResponseHelper provides standardized response formatting
type ResponseHelper struct{}

// NewResponseHelper creates a new ResponseHelper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// SuccessResponse sends a successful response
func (h *ResponseHelper) SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    data,
	})
}

// SuccessResponseWithPagination sends a successful response with pagination metadata
func (h *ResponseHelper) SuccessResponseWithPagination(c *fiber.Ctx, message string, data interface{}, pagination interface{}) error {
	return c.JSON(fiber.Map{
		"message":    message,
		"errors":     nil,
		"data":       data,
		"pagination": pagination,
	})
}

// ErrorResponse sends an error response
func (h *ResponseHelper) ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"message": message,
		"errors":  errors,
		"data":    nil,
	})
}

// ValidationError sends a validation error response
func (h *ResponseHelper) ValidationError(c *fiber.Ctx, message string, fieldErrors map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": message,
		"errors":  fieldErrors,
		"data":    nil,
	})
}
