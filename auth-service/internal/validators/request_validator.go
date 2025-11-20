package validators

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// LoginRequest represents login request data
type LoginRequest struct {
	Email string `json:"email"`
	Pin   string `json:"pin"`
}

// RefreshTokenRequest represents refresh token request data
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// PinRequest represents PIN request data
type PinRequest struct {
	Email string `json:"email"`
}

// RequestValidator handles request validation
type RequestValidator struct{}

// NewRequestValidator creates a new request validator
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

// ParseLoginRequest parses and validates login request
func (v *RequestValidator) ParseLoginRequest(c *fiber.Ctx) (*LoginRequest, error) {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}

	if err := v.validateLoginRequest(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

// ParseRefreshTokenRequest parses and validates refresh token request
func (v *RequestValidator) ParseRefreshTokenRequest(c *fiber.Ctx) (*RefreshTokenRequest, error) {
	var req RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}

	if err := v.validateRefreshTokenRequest(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

// ParsePinRequest parses and validates PIN request
func (v *RequestValidator) ParsePinRequest(c *fiber.Ctx) (*PinRequest, error) {
	var req PinRequest

	if err := c.BodyParser(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}

	if err := v.validatePinRequest(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

// validateLoginRequest validates login request data
func (v *RequestValidator) validateLoginRequest(req *LoginRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Pin == "" {
		return fmt.Errorf("PIN is required")
	}
	return nil
}

// validateRefreshTokenRequest validates refresh token request data
func (v *RequestValidator) validateRefreshTokenRequest(req *RefreshTokenRequest) error {
	if req.RefreshToken == "" {
		return fmt.Errorf("refresh token is required")
	}
	return nil
}

// validatePinRequest validates PIN request data
func (v *RequestValidator) validatePinRequest(req *PinRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}
