package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	config       *Config
	userRepo     *UserRepository
	tokenService *TokenService
	googleConfig *oauth2.Config
}

func NewOAuthHandler(config *Config, userRepo *UserRepository, tokenService *TokenService) *OAuthHandler {
	googleConfig := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthHandler{
		config:       config,
		userRepo:     userRepo,
		tokenService: tokenService,
		googleConfig: googleConfig,
	}
}

// GoogleUserInfo represents the user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// GoogleLogin initiates the Google OAuth flow
// @Summary Initiate Google OAuth login
// @Description Redirects to Google for authentication
// @Tags OAuth
// @Produce json
// @Success 302 {string} string "Redirect to Google"
// @Router /oauth/google [get]
func (h *OAuthHandler) GoogleLogin(c *fiber.Ctx) error {
	// Generate random state for CSRF protection
	state := generateRandomState()

	// Store state in cookie for verification in callback
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: "Lax",
		Domain:   h.config.CookieDomain,
	})

	// Redirect to Google
	url := h.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from Google OAuth
// @Summary Handle Google OAuth callback
// @Description Processes the OAuth callback, creates/updates user, and issues tokens
// @Tags OAuth
// @Produce json
// @Param state query string true "OAuth state"
// @Param code query string true "OAuth code"
// @Success 302 {string} string "Redirect to web app"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /oauth/google/callback [get]
func (h *OAuthHandler) GoogleCallback(c *fiber.Ctx) error {
	ctx := context.Background()

	// Verify state
	state := c.Query("state")
	cookieState := c.Cookies("oauth_state")

	if state == "" || state != cookieState {
		return h.redirectWithError(c, "Invalid state parameter")
	}

	// Clear state cookie
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Domain:   h.config.CookieDomain,
	})

	// Exchange code for token
	code := c.Query("code")
	if code == "" {
		return h.redirectWithError(c, "Missing authorization code")
	}

	token, err := h.googleConfig.Exchange(ctx, code)
	if err != nil {
		return h.redirectWithError(c, "Failed to exchange token")
	}

	// Get user info from Google
	userInfo, err := h.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return h.redirectWithError(c, "Failed to get user info")
	}

	if !userInfo.VerifiedEmail {
		return h.redirectWithError(c, "Email not verified")
	}

	// Find or create user
	user, err := h.userRepo.FindByGoogleID(ctx, userInfo.ID)
	if err == mongo.ErrNoDocuments {
		// Check if user exists with this email
		user, err = h.userRepo.FindByEmail(ctx, userInfo.Email)
		if err == mongo.ErrNoDocuments {
			// Create new user
			now := time.Now().UTC()
			user = &User{
				Email:        userInfo.Email,
				Username:     userInfo.Name,
				GoogleID:     &userInfo.ID,
				IsVerified:   true,
				RegisteredAt: &now,
			}
			if err := h.userRepo.Create(ctx, user); err != nil {
				return h.redirectWithError(c, "Failed to create user")
			}
		} else if err != nil {
			return h.redirectWithError(c, "Database error")
		} else {
			// Update existing user with Google ID
			user.GoogleID = &userInfo.ID
			user.IsVerified = true
			if err := h.userRepo.Update(ctx, user); err != nil {
				return h.redirectWithError(c, "Failed to update user")
			}
		}
	} else if err != nil {
		return h.redirectWithError(c, "Database error")
	}

	// Update last login
	_ = h.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	accessToken, err := h.tokenService.GenerateAccessToken(user)
	if err != nil {
		return h.redirectWithError(c, "Failed to generate access token")
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return h.redirectWithError(c, "Failed to generate refresh token")
	}

	// Store refresh token
	refreshExpiry := time.Now().UTC().Add(time.Duration(h.config.RefreshTokenExpiry) * time.Hour)
	if err := h.userRepo.UpdateRefreshToken(ctx, user.ID, refreshToken, refreshExpiry); err != nil {
		return h.redirectWithError(c, "Failed to store refresh token")
	}

	// Set cookies
	h.setAuthCookies(c, accessToken, refreshToken)

	// Redirect to web app
	return c.Redirect(h.config.WebURL+"/auth/callback?success=true", fiber.StatusTemporaryRedirect)
}

// getGoogleUserInfo fetches user info from Google
func (h *OAuthHandler) getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// redirectWithError redirects to web app with error
func (h *OAuthHandler) redirectWithError(c *fiber.Ctx, message string) error {
	return c.Redirect(h.config.WebURL+"/auth/callback?error="+message, fiber.StatusTemporaryRedirect)
}

// setAuthCookies sets authentication cookies
func (h *OAuthHandler) setAuthCookies(c *fiber.Ctx, accessToken, refreshToken string) {
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

// generateRandomState generates a random state for CSRF protection
func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
