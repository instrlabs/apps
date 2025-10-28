package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/instrlabs/shared/email"
)

type UserHandler struct {
	cfg      *Config
	userRepo *UserRepository
}

func NewUserHandler(cfg *Config, userRepo *UserRepository) *UserHandler {
	return &UserHandler{
		cfg:      cfg,
		userRepo: userRepo,
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

func (h *UserHandler) generateAccessToken(userID string, roles []string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
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

func (h *UserHandler) getCookieDomain(c *fiber.Ctx) string {
	if h.cfg.Environment == "production" {
		return ".arthadede.com"
	}

	userOrigin := c.Get("x-user-origin")
	if userOrigin != "" {
		if !strings.HasPrefix(userOrigin, "http://localhost") {
			domain := strings.Split(userOrigin, "//")[1]
			return "." + strings.Join(strings.Split(domain, ".")[1:], ".")
		}
	}

	return ".localhost"
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

	accessToken, err := h.generateAccessToken(user.ID.Hex(), []string{"user"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	refreshToken, err := h.generateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.userRepo.UpdateRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	now := time.Now().UTC()
	domain := h.getCookieDomain(c)

	log.Info("Login: Setting access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  now.Add(time.Duration(h.cfg.TokenExpiryHours) * time.Hour),
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	log.Info("Login: Setting refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  now.Add(time.Duration(h.cfg.RefreshExpiryHours) * time.Hour),
		MaxAge:   h.cfg.RefreshExpiryHours * 3600,
	})

	log.Infof("Login: User logged in successfully: %s", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	log.Info("RefreshToken: Processing token refresh request")

	refreshToken := c.Get("x-user-refresh")
	log.Infof("RefreshToken: Refresh token: %s", refreshToken)
	user := h.userRepo.FindByRefreshToken(refreshToken)
	if user == nil || user.ID.IsZero() {
		log.Warn("RefreshToken: Invalid refresh token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	newAccessToken, err := h.generateAccessToken(user.ID.Hex(), []string{"user"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	newRefreshToken, err := h.generateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.userRepo.UpdateRefreshToken(user.ID.Hex(), newRefreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	now := time.Now().UTC()
	domain := h.getCookieDomain(c)

	log.Info("RefreshToken: Update access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "access_token",
		Value:    newAccessToken,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  now.Add(time.Duration(h.cfg.TokenExpiryHours) * time.Hour),
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	log.Info("RefreshToken: Setting new refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  now.Add(time.Duration(h.cfg.RefreshExpiryHours) * time.Hour),
		MaxAge:   h.cfg.RefreshExpiryHours * 3600,
	})

	log.Info("RefreshToken: Token refreshed successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"errors":  nil,
		"data":    nil,
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
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Errorf("GoogleCallback: Failed to get user info: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("GoogleCallback: Failed to read response body: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
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
			"message": ErrInternalServer,
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
					"message": ErrInternalServer,
					"errors":  nil,
					"data":    nil,
				})
			}
			user = created
		} else {
			user = u2
			if err := h.userRepo.UpdateGoogleID(user.ID.Hex(), googleInfo.ID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": ErrInternalServer,
					"errors":  nil,
					"data":    nil,
				})
			}
		}
	}

	accessToken, err := h.generateAccessToken(user.ID.Hex(), []string{"user"})
	if err != nil {
		log.Errorf("GoogleCallback: Failed to generate access token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	refreshToken, err := h.generateRefreshToken()
	if err != nil {
		log.Errorf("GoogleCallback: Failed to generate refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.userRepo.UpdateRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		log.Errorf("GoogleCallback: Failed to update refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	now := time.Now().UTC()
	domain := h.getCookieDomain(c)

	log.Info("GoogleCallback: Setting access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  now.Add(time.Duration(h.cfg.TokenExpiryHours) * time.Hour),
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  now.Add(time.Duration(h.cfg.RefreshExpiryHours) * time.Hour),
		MaxAge:   h.cfg.RefreshExpiryHours * 3600,
	})

	log.Infof("GoogleCallback: User logged in successfully: %s", googleInfo.Email)
	return c.Redirect(h.cfg.WebUrl, fiber.StatusFound)
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
		"message": "Profile retrieved successfully",
		"errors":  nil,
		"data":    map[string]interface{}{"user": user},
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	log.Info("Logout: Processing logout request using Locals UserID")

	userId, _ := c.Locals("userId").(string)
	if err := h.userRepo.ClearRefreshToken(userId); err != nil {
		log.Errorf("Logout: Failed to logout user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to logout user",
			"errors":  nil,
			"data":    nil,
		})
	}

	domain := h.getCookieDomain(c)

	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
	})

	c.Cookie(&fiber.Cookie{
		Domain:   domain,
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
	})

	log.Infof("Logout: User logged out successfully: %s", userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
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
				"message": ErrInternalServer,
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
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	if !h.cfg.PinEnabled {
		log.Infof("SendPin: PIN_ENABLED enabled. Using fixed PIN 000000 for email: %s", input.Email)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "PIN sent",
			"errors":  nil,
			"data":    nil,
		})
	}

	subject := "Your Login PIN"
	body := fmt.Sprintf("Your one-time PIN is: %s. It expires in 10 minutes.", pin)
	if err := email.SendEmail(input.Email, subject, body); err != nil {
		log.Errorf("SendPin: Failed to send email: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "PIN sent",
		"errors":  nil,
		"data":    nil,
	})
}
