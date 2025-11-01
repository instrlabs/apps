package internal

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	config       *Config
	userRepo     *UserRepository
	tokenService *TokenService
}

func NewAuthHandler(config *Config, userRepo *UserRepository, tokenService *TokenService) *AuthHandler {
	return &AuthHandler{
		config:       config,
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// RefreshToken refreshes an access token using a refresh token
// @Summary Refresh access token
// @Description Issues a new access token using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// Get refresh token from body or cookie
	var req RefreshTokenRequest
	refreshToken := ""

	// Try to get from body first
	if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
		refreshToken = req.RefreshToken
	} else {
		// Try to get from cookie
		refreshToken = c.Cookies("refresh_token")
	}

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	ctx := context.Background()

	// Validate refresh token
	claims, err := h.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	// Parse user ID
	userID, err := ParseUserID(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Find user
	user, err := h.userRepo.FindByID(ctx, userID)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Verify refresh token matches stored token
	if user.RefreshToken == nil || *user.RefreshToken != refreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	// Check if refresh token has expired
	if user.RefreshTokenExpires == nil || time.Now().UTC().After(*user.RefreshTokenExpires) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Refresh token has expired",
		})
	}

	// Generate new access token
	accessToken, err := h.tokenService.GenerateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access token",
		})
	}

	// Optionally rotate refresh token (security best practice)
	newRefreshToken, err := h.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// Update refresh token in database
	refreshExpiry := time.Now().UTC().Add(time.Duration(h.config.RefreshTokenExpiry) * time.Hour)
	if err := h.userRepo.UpdateRefreshToken(ctx, user.ID, newRefreshToken, refreshExpiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update refresh token",
		})
	}

	// Set cookies
	h.setAuthCookies(c, accessToken, newRefreshToken)

	// Calculate expiry
	expiresAt := time.Now().UTC().Add(time.Duration(h.config.AccessTokenExpiry) * time.Hour)

	return c.JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		User:         user,
	})
}

// Logout logs out a user by revoking their refresh token
// @Summary Logout user
// @Description Revokes the user's refresh token
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get user ID from context (set by authentication middleware)
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := ParseUserID(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx := context.Background()

	// Clear refresh token from database
	if err := h.userRepo.ClearRefreshToken(ctx, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout",
		})
	}

	// Clear cookies
	h.clearAuthCookies(c)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// GetProfile returns the authenticated user's profile
// @Summary Get user profile
// @Description Returns the authenticated user's profile information
// @Tags User
// @Produce json
// @Success 200 {object} User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by authentication middleware)
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := ParseUserID(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx := context.Background()

	// Find user
	user, err := h.userRepo.FindByID(ctx, userID)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.JSON(user)
}

// setAuthCookies sets authentication cookies
func (h *AuthHandler) setAuthCookies(c *fiber.Ctx, accessToken, refreshToken string) {
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

// clearAuthCookies clears authentication cookies
func (h *AuthHandler) clearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Domain:   h.config.CookieDomain,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Domain:   h.config.CookieDomain,
	})
}
