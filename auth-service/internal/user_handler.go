package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
	"github.com/instrlabs/shared/functionx"
)

type UserHandler struct {
	cfg      *Config
	userRepo UserRepositoryInterface
}

func NewUserHandler(cfg *Config, userRepo UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (h *UserHandler) sendErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) sendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"errors":  nil,
		"data":    data,
	})
}

func (h *UserHandler) validateLoginInput(input struct {
	Email string `json:"email" validate:"required,email"`
	Pin   string `json:"pin" validate:"required"`
}) error {
	if input.Email == "" {
		return fmt.Errorf(ErrEmailRequired)
	}
	if input.Pin == "" {
		return fmt.Errorf(ErrPasswordRequired)
	}
	return nil
}

func (h *UserHandler) createTokensForUser(userID string) (string, string, error) {
	accessToken := GenerateAccessToken(userID)
	refreshToken := GenerateRefreshToken()

	err := h.userRepo.AddRefreshToken(userID, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to add refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (h *UserHandler) buildTokenResponse(accessToken, refreshToken string, expiresIn int) fiber.Map {
	return fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
	}
}

func (h *UserHandler) findUserByGoogleID(googleID string) (*User, error) {
	var user User
	err := h.userRepo.FindByGoogleID(googleID, &user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found by Google ID")
		}
		return nil, fmt.Errorf("failed to find user by Google ID: %w", err)
	}
	return &user, nil
}

func (h *UserHandler) findUserByEmailAndLinkGoogle(email, googleID string) (*User, error) {
	var user User
	err := h.userRepo.FindByEmail(email, &user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found by email")
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	err = h.userRepo.UpdateGoogleID(user.ID.Hex(), googleID)
	if err != nil {
		return nil, fmt.Errorf("failed to update Google ID: %w", err)
	}

	return &user, nil
}

func (h *UserHandler) createNewGoogleUser(email, googleID string) (*User, error) {
	newUser := NewGoogleUser(email, googleID)
	err := h.userRepo.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google user: %w", err)
	}
	return newUser, nil
}

func (h *UserHandler) findOrCreateGoogleUser(googleID, email string) (*User, error) {
	user, err := h.findUserByGoogleID(googleID)
	if err == nil {
		return user, nil
	}

	user, err = h.findUserByEmailAndLinkGoogle(email, googleID)
	if err == nil {
		return user, nil
	}

	return h.createNewGoogleUser(email, googleID)
}

func (h *UserHandler) parseLogoutInput(c *fiber.Ctx) (struct {
	RefreshToken string `json:"refresh_token"`
}, error) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := c.BodyParser(&input)
	return input, err
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	log.Info("Login: Processing login request")

	var input struct {
		Email string `json:"email" validate:"required,email"`
		Pin   string `json:"pin" validate:"required"`
	}

	// Parse and validate input with early returns
	if err := c.BodyParser(&input); err != nil {
		log.Warnf("Login: Invalid request body: %v", err)
		return h.sendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := h.validateLoginInput(input); err != nil {
		log.Infof("Login: Validation failed: %v", err)
		return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	log.Infof("Login: Attempting to login user with email: %s", input.Email)

	// Find user with early return on error
	user, err := h.findUserByEmail(input.Email)
	if err != nil {
		log.Infof("Login: Invalid credentials for email: %s", input.Email)
		return h.sendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidCredentials)
	}

	// Validate PIN with early return
	if !user.ComparePin(input.Pin) {
		log.Infof("Login: Invalid PIN for email: %s", input.Email)
		return h.sendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidCredentials)
	}

	// Handle post-login tasks
	h.handlePostLoginTasks(user)

	// Create tokens with early return on error
	accessToken, refreshToken, err := h.createTokensForUser(user.ID.Hex())
	if err != nil {
		log.Errorf("Login: Failed to create tokens: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrCreateSession)
	}

	log.Infof("Login: User logged in successfully: %s", input.Email)
	return h.sendSuccessResponse(c, SuccessLogin, h.buildTokenResponse(accessToken, refreshToken, h.cfg.TokenExpiryHours*3600))
}

// Helper method for Login
func (h *UserHandler) findUserByEmail(email string) (*User, error) {
	var user User
	err := h.userRepo.FindByEmail(email, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *UserHandler) handlePostLoginTasks(user *User) {
	_ = h.userRepo.ClearPin(user.ID.Hex())
	if user.RegisteredAt == nil {
		_ = h.userRepo.SetRegisteredAt(user.ID.Hex())
	}
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	log.Info("RefreshToken: Processing token refresh request")

	// Parse and validate input
	refreshToken, err := h.parseRefreshTokenInput(c)
	if err != nil {
		log.Warnf("RefreshToken: Invalid input: %v", err)
		return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Get user ID from middleware
	userId, _ := c.Locals("userId").(string)
	log.Infof("RefreshToken: Processing refresh for user %s", userId)

	// Validate refresh token
	if err := h.userRepo.ValidateRefreshToken(userId, refreshToken); err != nil {
		log.Warn("RefreshToken: Invalid refresh token")
		return h.sendErrorResponse(c, fiber.StatusUnauthorized, ErrInvalidToken)
	}

	// Get user for token generation
	user, err := h.findUserByID(userId)
	if err != nil {
		log.Warn("RefreshToken: User not found")
		return h.sendErrorResponse(c, fiber.StatusUnauthorized, ErrUserNotFound)
	}

	// Rotate tokens (remove old, add new)
	newAccessToken, newRefreshToken, err := h.rotateRefreshTokens(user.ID.Hex(), refreshToken)
	if err != nil {
		log.Errorf("RefreshToken: Failed to rotate tokens: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrUpdateSession)
	}

	log.Info("RefreshToken: Token refreshed successfully")
	return h.sendSuccessResponse(c, SuccessTokenRefreshed, h.buildTokenResponse(newAccessToken, newRefreshToken, h.cfg.TokenExpiryHours*3600))
}

// Helper methods for RefreshToken
func (h *UserHandler) parseRefreshTokenInput(c *fiber.Ctx) (string, error) {
	var input struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return "", fmt.Errorf(ErrInvalidRequestBody)
	}

	if input.RefreshToken == "" {
		return "", fmt.Errorf(ErrRefreshTokenRequired)
	}

	return input.RefreshToken, nil
}

func (h *UserHandler) findUserByID(userID string) (*User, error) {
	var user User
	err := h.userRepo.FindByID(userID, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *UserHandler) rotateRefreshTokens(userID, oldRefreshToken string) (string, string, error) {
	// Remove old refresh token
	if err := h.userRepo.RemoveRefreshToken(userID, oldRefreshToken); err != nil {
		return "", "", fmt.Errorf("failed to remove old refresh token: %w", err)
	}

	// Create new tokens
	newAccessToken := GenerateAccessToken(userID)
	newRefreshToken := GenerateRefreshToken()

	// Add new refresh token
	if err := h.userRepo.AddRefreshToken(userID, newRefreshToken); err != nil {
		return "", "", fmt.Errorf("failed to add new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (h *UserHandler) GoogleLogin(c *fiber.Ctx) error {
	log.Info("GoogleLogin: Initiating Google OAuth login")

	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)
	conf := &oauth2.Config{
		ClientID:     h.cfg.GoogleClientID,
		ClientSecret: h.cfg.GoogleClientSecret,
		RedirectURL:  h.cfg.GoogleRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	redirectUrl := conf.AuthCodeURL(state)
	log.Infof("GoogleLogin: Redirecting to Google OAuth URL: %s", redirectUrl)

	return c.Redirect(redirectUrl)
}

func (h *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	log.Info("GoogleCallback: Processing Google OAuth callback")

	// Validate authorization code
	code := c.Query("code")
	if code == "" {
		log.Info("GoogleCallback: Missing authorization code")
		return h.sendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidToken)
	}

	// Exchange code for token
	token, err := h.exchangeGoogleCode(code)
	if err != nil {
		log.Errorf("GoogleCallback: Failed to exchange token: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrExchangeToken)
	}

	// Get Google user info
	googleInfo, err := h.getGoogleUserInfo(token)
	if err != nil {
		log.Errorf("GoogleCallback: Failed to get user info: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrGetUserInfo)
	}

	// Find or create user
	user, err := h.findOrCreateGoogleUser(googleInfo.ID, googleInfo.Email)
	if err != nil {
		log.Errorf("GoogleCallback: Failed to find/create user: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrCreateGoogleUser)
	}

	// Create tokens for user
	accessToken, refreshToken, err := h.createTokensForUser(user.ID.Hex())
	if err != nil {
		log.Errorf("GoogleCallback: Failed to create tokens: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrCreateSession)
	}

	// Redirect with tokens
	redirectURL := h.buildGoogleRedirectURL(accessToken, refreshToken)
	log.Infof("GoogleCallback: User logged in successfully: %s", googleInfo.Email)
	return c.Redirect(redirectURL, fiber.StatusFound)
}

// Helper methods for Google OAuth
func (h *UserHandler) exchangeGoogleCode(code string) (*oauth2.Token, error) {
	conf := &oauth2.Config{
		ClientID:     h.cfg.GoogleClientID,
		ClientSecret: h.cfg.GoogleClientSecret,
		RedirectURL:  h.cfg.GoogleRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf.Exchange(context.Background(), code)
}

func (h *UserHandler) getGoogleUserInfo(token *oauth2.Token) (*struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}, error) {
	conf := &oauth2.Config{
		ClientID:     h.cfg.GoogleClientID,
		ClientSecret: h.cfg.GoogleClientSecret,
		RedirectURL:  h.cfg.GoogleRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	err = json.Unmarshal(data, &googleInfo)
	return &googleInfo, err
}

func (h *UserHandler) buildGoogleRedirectURL(accessToken, refreshToken string) string {
	expiresIn := h.cfg.TokenExpiryHours * 3600
	return fmt.Sprintf("%s?access_token=%s&refresh_token=%s&token_type=Bearer&expires_in=%s",
		h.cfg.WebUrl,
		url.QueryEscape(accessToken),
		url.QueryEscape(refreshToken),
		strconv.Itoa(expiresIn),
	)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	log.Info("GetProfile: Processing profile request using locals.userId")

	userId, _ := c.Locals("userId").(string)
	var user User
	err := h.userRepo.FindByID(userId, &user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Infof("GetProfile: User not found for userId %s", userId)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": ErrUserNotFound,
				"errors":  nil,
				"data":    nil,
			})
		}
		log.Errorf("GetProfile: Failed to find user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("GetProfile: Profile retrieved successfully for user: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessProfileRetrieved,
		"errors":  nil,
		"data":    fiber.Map{"user": user},
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	log.Info("Logout: Processing logout request")

	userId, _ := c.Locals("userId").(string)

	// Parse logout input (optional)
	input, err := h.parseLogoutInput(c)
	if err != nil {
		log.Warnf("Logout: Failed to parse logout input: %v", err)
	}

	// Logout based on whether specific refresh token was provided
	if input.RefreshToken != "" {
		err = h.logoutSpecificToken(userId, input.RefreshToken)
	} else {
		err = h.logoutAllTokens(userId)
	}

	if err != nil {
		log.Errorf("Logout: Failed to logout: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrFailedToLogout)
	}

	log.Infof("Logout: User logged out successfully: %s", userId)
	return h.sendSuccessResponse(c, SuccessLogoutSuccessful, nil)
}

// Helper methods for Logout
func (h *UserHandler) logoutSpecificToken(userID, refreshToken string) error {
	if err := h.userRepo.RemoveRefreshToken(userID, refreshToken); err != nil {
		return fmt.Errorf("failed to remove refresh token: %w", err)
	}
	return nil
}

func (h *UserHandler) logoutAllTokens(userID string) error {
	if err := h.userRepo.ClearAllRefreshTokens(userID); err != nil {
		return fmt.Errorf("failed to clear all refresh tokens: %w", err)
	}
	return nil
}

func (h *UserHandler) SendPin(c *fiber.Ctx) error {
	log.Info("SendPin: Processing send-pin request")

	// Parse and validate input
	email, err := h.parseSendPinInput(c)
	if err != nil {
		log.Warnf("SendPin: Invalid input: %v", err)
		return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	log.Infof("SendPin: Sending pin with email: %s", email)

	// Find or create user
	_, err = h.findOrCreateUserByEmail(email)
	if err != nil {
		log.Errorf("SendPin: Failed to find/create user: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrCreateUser)
	}

	// Generate and set PIN
	pin, err := h.generateAndSetPin(email)
	if err != nil {
		log.Errorf("SendPin: Failed to set PIN: %v", err)
		return h.sendErrorResponse(c, fiber.StatusInternalServerError, ErrSetPin)
	}

	// Handle non-PIN enabled case
	if !h.cfg.PinEnabled {
		log.Infof("SendPin: PIN_ENABLED disabled. Using fixed PIN for email: %s", email)
		return h.sendSuccessResponse(c, SuccessPinSent, nil)
	}

	// Send email with PIN
	h.sendPinEmail(email, pin)

	return h.sendSuccessResponse(c, SuccessPinSent, nil)
}

// Helper methods for SendPin
func (h *UserHandler) parseSendPinInput(c *fiber.Ctx) (string, error) {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return "", fmt.Errorf(ErrInvalidRequestBody)
	}

	if input.Email == "" {
		return "", fmt.Errorf(ErrEmailRequired)
	}

	return input.Email, nil
}

func (h *UserHandler) findOrCreateUserByEmail(email string) (*User, error) {
	user, err := h.findUserByEmail(email)
	if err == nil {
		return user, nil
	}

	if err == mongo.ErrNoDocuments {
		newUser := NewUser(email)
		if createErr := h.userRepo.Create(newUser); createErr != nil {
			return nil, fmt.Errorf("failed to create user: %w", createErr)
		}
		return newUser, nil
	}

	return nil, fmt.Errorf("failed to find user: %w", err)
}

func (h *UserHandler) generateAndSetPin(email string) (string, error) {
	pin := "000000"
	if h.cfg.PinEnabled {
		pin = GenerateSixDigitPIN()
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	err := h.userRepo.SetPinWithExpiry(email, string(hash))
	if err != nil {
		return "", fmt.Errorf("failed to set PIN: %w", err)
	}

	return pin, nil
}

func (h *UserHandler) sendPinEmail(email, pin string) {
	subject := "Your Login PIN"
	body := fmt.Sprintf("Your one-time PIN is: %s. It expires in 10 minutes.", pin)
	functionx.SendEmail(email, subject, body)
}
