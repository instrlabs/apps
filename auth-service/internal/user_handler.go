package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
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

func (h *UserHandler) Login(c *fiber.Ctx) error {
	log.Info("Login: Processing login request")

	var input struct {
		Email string `json:"email" validate:"required,email"`
		Pin   string `json:"pin" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Warnf("Login: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Email == "" {
		log.Info("Login: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrEmailRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Pin == "" {
		log.Info("Login: Pin is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrPasswordRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("Login: Attempting to login user with email: %s", input.Email)
	user := h.userRepo.FindByEmail(input.Email)
	if user == nil || user.ID.IsZero() {
		log.Infof("Login: Invalid credentials for email: %s", input.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidCredentials,
			"errors":  nil,
			"data":    nil,
		})
	}
	if !user.ComparePin(input.Pin) {
		log.Infof("Login: Invalid PIN for email: %s", input.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidCredentials,
			"errors":  nil,
			"data":    nil,
		})
	}

	_ = h.userRepo.ClearPin(user.ID.Hex())
	if user.RegisteredAt == nil || user.RegisteredAt.IsZero() {
		if err := h.userRepo.SetRegisteredAt(user.ID.Hex()); err != nil {
			log.Errorf("Login: Failed to set RegisteredAt for user %s: %v", user.ID.Hex(), err)
		}
	}

	accessToken := GenerateAccessToken(user.ID.Hex(), "")
	refreshToken := GenerateRefreshToken()

	if err := h.userRepo.AddRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		log.Errorf("Login: Failed to add refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrCreateSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("Login: User logged in successfully: %s", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessLogin,
		"errors":  nil,
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    h.cfg.TokenExpiryHours * 3600,
		},
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	log.Info("RefreshToken: Processing token refresh request")

	var input struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Warnf("RefreshToken: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.RefreshToken == "" {
		log.Info("RefreshToken: Refresh token is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrRefreshTokenRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	userId, _ := c.Locals("userId").(string)
	refreshToken := input.RefreshToken
	log.Infof("RefreshToken: Refresh token received from request body")

	if !h.userRepo.ValidateRefreshToken(userId, refreshToken) {
		log.Warn("RefreshToken: Invalid refresh token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	user := h.userRepo.FindByID(userId)
	if user == nil || user.ID.IsZero() {
		log.Warn("RefreshToken: User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUserNotFound,
			"errors":  nil,
			"data":    nil,
		})
	}

	newAccessToken := GenerateAccessToken(user.ID.Hex(), "")
	newRefreshToken := GenerateRefreshToken()

	if err := h.userRepo.RemoveRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		log.Errorf("RefreshToken: Failed to remove old refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrUpdateSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	if err := h.userRepo.AddRefreshToken(user.ID.Hex(), newRefreshToken); err != nil {
		log.Errorf("RefreshToken: Failed to add new refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrUpdateSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Info("RefreshToken: Token refreshed successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessTokenRefreshed,
		"errors":  nil,
		"data": fiber.Map{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"token_type":    "Bearer",
			"expires_in":    h.cfg.TokenExpiryHours * 3600,
		},
	})
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

	code := c.Query("code")
	if code == "" {
		log.Info("GoogleCallback: Missing authorization code")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

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
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		log.Errorf("GoogleCallback: Failed to exchange token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrExchangeToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Errorf("GoogleCallback: Failed to get user info: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGetUserInfo,
			"errors":  nil,
			"data":    nil,
		})
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("GoogleCallback: Failed to read response body: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGetUserInfo,
			"errors":  nil,
			"data":    nil,
		})
	}

	var googleInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(data, &googleInfo); err != nil {
		log.Errorf("GoogleCallback: Failed to unmarshal user info: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrParseUserInfo,
			"errors":  nil,
			"data":    nil,
		})
	}
	user := h.userRepo.FindByGoogleID(googleInfo.ID)
	if user == nil || user.ID.IsZero() {
		u2 := h.userRepo.FindByEmail(googleInfo.Email)
		if u2 == nil || u2.ID.IsZero() {
			newUser := NewGoogleUser(googleInfo.Email, googleInfo.ID)
			created := h.userRepo.Create(newUser)
			if created == nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": ErrCreateGoogleUser,
					"errors":  nil,
					"data":    nil,
				})
			}
			user = created
		} else {
			user = u2
			if err := h.userRepo.UpdateGoogleID(user.ID.Hex(), googleInfo.ID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": ErrUpdateGoogleID,
					"errors":  nil,
					"data":    nil,
				})
			}
		}
	}

	accessToken := GenerateAccessToken(user.ID.Hex(), "")
	refreshToken := GenerateRefreshToken()

	if err := h.userRepo.AddRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		log.Errorf("GoogleCallback: Failed to add refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrCreateSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	expiresIn := h.cfg.TokenExpiryHours * 3600
	redirectURL := fmt.Sprintf("%s?access_token=%s&refresh_token=%s&token_type=Bearer&expires_in=%s",
		h.cfg.WebUrl,
		url.QueryEscape(accessToken),
		url.QueryEscape(refreshToken),
		strconv.Itoa(expiresIn),
	)

	log.Infof("GoogleCallback: User logged in successfully: %s", googleInfo.Email)
	return c.Redirect(redirectURL, fiber.StatusFound)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	log.Info("GetProfile: Processing profile request using locals.userId")

	userId, _ := c.Locals("userId").(string)
	user := h.userRepo.FindByID(userId)
	if user == nil || user.ID.IsZero() {
		log.Infof("GetProfile: User not found for userId %s", userId)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUserNotFound,
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

	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&input); err == nil && input.RefreshToken != "" {
		if err := h.userRepo.RemoveRefreshToken(userId, input.RefreshToken); err != nil {
			log.Errorf("Logout: Failed to remove refresh token: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": ErrFailedToLogout,
				"errors":  nil,
				"data":    nil,
			})
		}
	} else {
		if err := h.userRepo.ClearAllRefreshTokens(userId); err != nil {
			log.Errorf("Logout: Failed to clear all refresh tokens: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": ErrFailedToLogout,
				"errors":  nil,
				"data":    nil,
			})
		}
	}

	log.Infof("Logout: User logged out successfully: %s", userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessLogoutSuccessful,
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) SendPin(c *fiber.Ctx) error {
	log.Info("SendPin: Processing send-pin request")

	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Errorf("SendPin: Failed to parse body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Email == "" {
		log.Warnf("SendPin: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrEmailRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("SendPin: Sending pin with email: %s", input.Email)
	user := h.userRepo.FindByEmail(input.Email)
	if user == nil || user.ID.IsZero() {
		newUser := NewUser(input.Email)
		if created := h.userRepo.Create(newUser); created == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": ErrCreateUser,
				"errors":  nil,
				"data":    nil,
			})
		}
	}

	pin := "000000"
	if h.cfg.PinEnabled {
		pin = GenerateSixDigitPIN()
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err := h.userRepo.SetPinWithExpiry(input.Email, string(hash)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrSetPin,
			"errors":  nil,
			"data":    nil,
		})
	}

	if !h.cfg.PinEnabled {
		log.Infof("SendPin: PIN_ENABLED enabled. Using fixed PIN 000000 for email: %s", input.Email)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": SuccessPinSent,
			"errors":  nil,
			"data":    nil,
		})
	}

	subject := "Your Login PIN"
	body := fmt.Sprintf("Your one-time PIN is: %s. It expires in 10 minutes.", pin)
	functionx.SendEmail(input.Email, subject, body)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessPinSent,
		"errors":  nil,
		"data":    nil,
	})
}
