package internal

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type PinHandler struct {
	config       *Config
	userRepo     *UserRepository
	tokenService *TokenService
	emailService *EmailService
}

func NewPinHandler(config *Config, userRepo *UserRepository, tokenService *TokenService, emailService *EmailService) *PinHandler {
	return &PinHandler{
		config:       config,
		userRepo:     userRepo,
		tokenService: tokenService,
		emailService: emailService,
	}
}

// RequestPin sends a PIN to the user's email
// @Summary Request authentication PIN
// @Description Sends a PIN code to the user's email for authentication
// @Tags PIN Authentication
// @Accept json
// @Produce json
// @Param request body PinRequest true "PIN request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/pin/request [post]
func (h *PinHandler) RequestPin(c *fiber.Ctx) error {
	if !h.config.PinEnabled {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "PIN authentication is disabled",
		})
	}

	var req PinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	ctx := context.Background()

	// Generate PIN
	pin := generatePin(h.config.PinLength)

	// Hash PIN
	pinHash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PIN",
		})
	}

	// Set expiry
	expiresAt := time.Now().UTC().Add(time.Duration(h.config.PinExpiryMins) * time.Minute)

	// Check if user exists
	user, err := h.userRepo.FindByEmail(ctx, req.Email)
	if err == mongo.ErrNoDocuments {
		// Create new user with PIN
		now := time.Now().UTC()
		user = &User{
			Email:      req.Email,
			PinHash:    stringPtr(string(pinHash)),
			PinExpires: &expiresAt,
			IsVerified: false, // Will be verified after PIN verification
		}
		if err := h.userRepo.Create(ctx, user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	} else {
		// Update existing user's PIN
		if err := h.userRepo.UpdatePIN(ctx, req.Email, string(pinHash), expiresAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update PIN",
			})
		}
	}

	// Send PIN via email
	if err := h.emailService.SendPinEmail(req.Email, pin, h.config.PinExpiryMins); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "PIN sent to your email",
		"expires_in_minutes": h.config.PinExpiryMins,
	})
}

// VerifyPin verifies the PIN and issues tokens
// @Summary Verify authentication PIN
// @Description Verifies the PIN code and issues authentication tokens
// @Tags PIN Authentication
// @Accept json
// @Produce json
// @Param request body PinVerifyRequest true "PIN verify request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/pin/verify [post]
func (h *PinHandler) VerifyPin(c *fiber.Ctx) error {
	if !h.config.PinEnabled {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "PIN authentication is disabled",
		})
	}

	var req PinVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Pin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and PIN are required",
		})
	}

	ctx := context.Background()

	// Find user
	user, err := h.userRepo.FindByEmail(ctx, req.Email)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or PIN",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Verify PIN
	if user.PinHash == nil || *user.PinHash == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No PIN found for this user",
		})
	}

	// Check PIN expiry
	if user.PinExpires == nil || time.Now().UTC().After(*user.PinExpires) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "PIN has expired",
		})
	}

	// Compare PIN
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PinHash), []byte(req.Pin)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or PIN",
		})
	}

	// Mark user as verified and registered
	if !user.IsVerified {
		user.IsVerified = true
		now := time.Now().UTC()
		if user.RegisteredAt == nil {
			user.RegisteredAt = &now
		}
		if err := h.userRepo.Update(ctx, user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user",
			})
		}
	}

	// Update last login
	_ = h.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	accessToken, err := h.tokenService.GenerateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access token",
		})
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// Store refresh token
	refreshExpiry := time.Now().UTC().Add(time.Duration(h.config.RefreshTokenExpiry) * time.Hour)
	if err := h.userRepo.UpdateRefreshToken(ctx, user.ID, refreshToken, refreshExpiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store refresh token",
		})
	}

	// Set cookies
	h.setAuthCookies(c, accessToken, refreshToken)

	// Clear PIN hash (one-time use)
	user.PinHash = nil
	user.PinExpires = nil
	_ = h.userRepo.Update(ctx, user)

	// Calculate expiry
	expiresAt := time.Now().UTC().Add(time.Duration(h.config.AccessTokenExpiry) * time.Hour)

	// Return response
	return c.JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         user,
	})
}

// setAuthCookies sets authentication cookies
func (h *PinHandler) setAuthCookies(c *fiber.Ctx, accessToken, refreshToken string) {
	accessExpiry := time.Now().Add(time.Duration(h.config.AccessTokenExpiry) * time.Hour)
	refreshExpiry := time.Now().Add(time.Duration(h.config.RefreshTokenExpiry) * time.Hour)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  accessExpiry,
		HTTPOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: "Lax",
		Domain:   h.config.CookieDomain,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshExpiry,
		HTTPOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: "Lax",
		Domain:   h.config.CookieDomain,
	})
}

// generatePin generates a random PIN of specified length
func generatePin(length int) string {
	const digits = "0123456789"
	pin := make([]byte, length)

	for i := range pin {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		pin[i] = digits[num.Int64()]
	}

	return string(pin)
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}
