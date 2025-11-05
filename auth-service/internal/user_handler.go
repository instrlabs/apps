package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/instrlabs/shared/functionx"
)

type UserHandler struct {
	cfg         *Config
	userRepo    *UserRepository
	sessionRepo *UserSessionRepository
}

func NewUserHandler(cfg *Config, userRepo *UserRepository, sessionRepo *UserSessionRepository) *UserHandler {
	return &UserHandler{
		cfg:         cfg,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func generateSixDigitPIN() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Errorf("generateSixDigitPIN: Failed to generate six digit PIN: %v", err)
		return ""
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func (h *UserHandler) generateAccessToken(userID, sessionID string) (string, error) {
	now := time.Now().UTC()
	expirationTime := now.Add(time.Duration(h.cfg.TokenExpiryHours) * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"session_id": sessionID,
		"iat":        now.Unix(),
		"exp":        expirationTime.Unix(),
	})

	return token.SignedString([]byte(h.cfg.JWTSecret))
}

func (h *UserHandler) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
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

	userIP, _ := c.Locals("userIP").(string)
	userAgent, _ := c.Locals("userAgent").(string)

	session, err := h.sessionRepo.CreateUserSession(user.ID.Hex(), userIP, userAgent)
	if err != nil {
		log.Errorf("Login: Failed to create session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrCreateSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	accessToken, err := h.generateAccessToken(user.ID.Hex(), session.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGenerateAccessToken,
			"errors":  nil,
			"data":    nil,
		})
	}
	refreshToken, err := h.generateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGenerateRefreshToken,
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.sessionRepo.UpdateUserSessionRefreshToken(session.ID, refreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrUpdateSession,
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

	refreshToken := input.RefreshToken
	log.Infof("RefreshToken: Refresh token received from request body")

	session, err := h.sessionRepo.FindUserSessionByRefreshToken(refreshToken)
	if session == nil || err != nil {
		log.Warn("RefreshToken: Invalid refresh token or session")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	userIP, _ := c.Locals("userIP").(string)
	userAgent, _ := c.Locals("userAgent").(string)
	deviceHash := GenerateDeviceHash(userIP, userAgent)

	if !h.sessionRepo.ValidateUserSession(session.ID, session.UserID, deviceHash) {
		log.Warn("RefreshToken: Device hash mismatch - possible token theft")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	user := h.userRepo.FindByID(session.UserID)
	if user == nil || user.ID.IsZero() {
		log.Warn("RefreshToken: User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUserNotFound,
			"errors":  nil,
			"data":    nil,
		})
	}

	newAccessToken, err := h.generateAccessToken(user.ID.Hex(), session.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGenerateAccessToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	newRefreshToken, err := h.generateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGenerateRefreshToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	if err := h.sessionRepo.UpdateUserSessionRefreshToken(session.ID, newRefreshToken); err != nil {
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

	userIP, _ := c.Locals("userIP").(string)
	userAgent, _ := c.Locals("userAgent").(string)

	session, err := h.sessionRepo.CreateUserSession(user.ID.Hex(), userIP, userAgent)
	if err != nil {
		log.Errorf("GoogleCallback: Failed to create session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrCreateSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	accessToken, err := h.generateAccessToken(user.ID.Hex(), session.ID.Hex())
	if err != nil {
		log.Errorf("GoogleCallback: Failed to generate access token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGenerateAccessToken,
			"errors":  nil,
			"data":    nil,
		})
	}
	refreshToken, err := h.generateRefreshToken()
	if err != nil {
		log.Errorf("GoogleCallback: Failed to generate refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGenerateRefreshToken,
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.sessionRepo.UpdateUserSessionRefreshToken(session.ID, refreshToken); err != nil {
		log.Errorf("GoogleCallback: Failed to update refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrUpdateSession,
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
	log.Info("Logout: Processing logout request using Locals UserID")

	userId, _ := c.Locals("userId").(string)
	sessionId, _ := c.Locals("sessionId").(string)

	if sessionId != "" {
		sessionObjID, err := primitive.ObjectIDFromHex(sessionId)
		if err != nil {
			log.Errorf("Logout: Invalid session ID format: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": ErrInvalidSessionID,
				"errors":  nil,
				"data":    nil,
			})
		}
		if err := h.sessionRepo.DeactivateUserSession(sessionObjID); err != nil {
			log.Errorf("Logout: Failed to deactivate session: %v", err)
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

func (h *UserHandler) GetDevices(c *fiber.Ctx) error {
	log.Info("GetDevices: Retrieving user devices")

	userId, _ := c.Locals("userId").(string)
	sessions, err := h.sessionRepo.GetUserSessions(userId)
	if err != nil {
		log.Errorf("GetDevices: Failed to get sessions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrGetUserSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("GetDevices: Retrieved %d devices for user %s", len(sessions), userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessDevicesRetrieved,
		"errors":  nil,
		"data":    fiber.Map{"sessions": sessions},
	})
}

func (h *UserHandler) RevokeDevice(c *fiber.Ctx) error {
	log.Info("RevokeDevice: Revoking device access")

	userId, _ := c.Locals("userId").(string)
	sessionId := c.Params("sessionId")

	sessionObjID, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		log.Errorf("RevokeDevice: Invalid session ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidSessionID,
			"errors":  nil,
			"data":    nil,
		})
	}

	session, err := h.sessionRepo.FindUserSessionByID(sessionObjID, userId)
	if err != nil {
		log.Errorf("RevokeDevice: Failed to find session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrFindUserSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	if session == nil {
		log.Warnf("RevokeDevice: Session %s not found for user %s", sessionId, userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": ErrDeviceNotFound,
			"errors":  nil,
			"data":    nil,
		})
	}

	if err := h.sessionRepo.DeactivateUserSession(sessionObjID); err != nil {
		log.Errorf("RevokeDevice: Failed to deactivate session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrDeactivateSession,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("RevokeDevice: Device %s revoked for user %s", sessionId, userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessDeviceRevoked,
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) LogoutAllDevices(c *fiber.Ctx) error {
	log.Info("LogoutAllDevices: Logging out from all devices")

	userId, _ := c.Locals("userId").(string)

	if err := h.sessionRepo.ClearAllUserSessions(userId); err != nil {
		log.Errorf("LogoutAllDevices: Failed to deactivate sessions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrClearAllSessions,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("LogoutAllDevices: User %s logged out from all devices", userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": SuccessLoggedOutAllDevices,
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
		pin = generateSixDigitPIN()
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
