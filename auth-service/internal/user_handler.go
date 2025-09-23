package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/smtp"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func (h *UserHandler) Login(c *fiber.Ctx) error {
	log.Info("Login: Processing login request")

	var input struct {
		Email string `json:"email" validate:"required,email"`
		Pin   string `json:"pin" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Errorf("Login: Invalid request body: %v", err)
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
			"errors": []fiber.Map{
				{
					"fieldName":    "pin",
					"errorMessage": ErrPasswordRequired,
				},
			},
			"data": nil,
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
	if user.RegisteredAt.IsZero() {
		if err := h.userRepo.SetRegisteredAt(user.ID.Hex()); err != nil {
			log.Errorf("Login: Failed to set RegisteredAt for user %s: %v", user.ID.Hex(), err)
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"roles":   []string{"user"},
	})
	accessToken, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	refreshToken := base64.StdEncoding.EncodeToString(b)
	if err := h.userRepo.UpdateRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	log.Info("Login: Setting access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "AccessToken",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	log.Info("Login: Setting refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "RefreshToken",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
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

	var input struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Errorf("RefreshToken: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.RefreshToken == "" {
		log.Info("RefreshToken: Refresh token cookie is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrRefreshTokenRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Info("RefreshToken: Attempting to refresh token")
	// Inline RefreshToken logic
	user := h.userRepo.FindByRefreshToken(input.RefreshToken)
	if user == nil || user.ID.IsZero() {
		log.Infof("RefreshToken: Invalid refresh token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}
	// Generate new tokens
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"roles":   []string{"user"},
	})
	accessToken, err := tok.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	newRefreshToken := base64.StdEncoding.EncodeToString(b)
	if err := h.userRepo.UpdateRefreshToken(user.ID.Hex(), newRefreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}

	log.Info("RefreshToken: Setting new access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "AccessToken",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	log.Info("RefreshToken: Setting new refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "RefreshToken",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
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
	rand.Read(b)
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
	url := conf.AuthCodeURL(state)
	log.Infof("GoogleLogin: Redirecting to Google OAuth URL: %s", url)

	return c.Redirect(url)
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

	log.Info("GoogleCallback: Handling Google callback with authorization code")
	// Build oauth2 config locally (no field stored)
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
		log.Errorf("GoogleCallback: Exchange error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	user := h.userRepo.FindByGoogleID(userInfo.ID)
	if user == nil || user.ID.IsZero() {
		u2 := h.userRepo.FindByEmail(userInfo.Email)
		if u2 == nil || u2.ID.IsZero() {
			newUser := NewGoogleUser(userInfo.Email, userInfo.ID, userInfo.Name)
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
			if err := h.userRepo.UpdateGoogleID(user.ID.Hex(), userInfo.ID, userInfo.Name); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": ErrInternalServer,
					"errors":  nil,
					"data":    nil,
				})
			}
		}
	}
	// Generate tokens
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"roles":   []string{"user"},
	})
	accessToken, err := tok.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	b2 := make([]byte, 32)
	if _, err := rand.Read(b2); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	refreshToken := base64.StdEncoding.EncodeToString(b2)
	if err := h.userRepo.UpdateRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	log.Info("GoogleCallback: Setting access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	log.Info("GoogleCallback: Setting refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	redirectURL := h.cfg.FEOAuthRedirect
	if redirectURL == "" {
		redirectURL = "/"
	}

	log.Infof("GoogleCallback: Redirecting to frontend: %s", redirectURL)
	return c.Redirect(redirectURL)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	log.Info("GetProfile: Processing profile request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		log.Info("GetProfile: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	user := h.userRepo.FindByID(userID)
	if user == nil || user.ID.IsZero() {
		log.Infof("GetProfile: User not found for UserID %s", userID)
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
		"data":    user,
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	log.Info("UpdateProfile: Processing profile update request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		log.Info("UpdateProfile: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Parse request body
	var request struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&request); err != nil {
		log.Errorf("UpdateProfile: Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Update profile (inline)
	if request.Name == "" {
		log.Warnf("UpdateProfile: name cannot be empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update profile",
			"errors":  fiber.Map{"general": "name cannot be empty"},
			"data":    nil,
		})
	}
	if err := h.userRepo.UpdateProfile(userID, request.Name); err != nil {
		log.Errorf("UpdateProfile: Failed to update profile: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update profile",
			"errors": fiber.Map{
				"general": err.Error(),
			},
			"data": nil,
		})
	}

	// Get updated user
	updatedUser := h.userRepo.FindByID(userID)
	if updatedUser == nil || updatedUser.ID.IsZero() {
		log.Errorf("UpdateProfile: Failed to get updated user after update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Infof("UpdateProfile: Profile updated successfully for user: %s", updatedUser.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"errors":  nil,
		"data": fiber.Map{
			"user": updatedUser,
		},
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	log.Info("Logout: Processing logout request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		log.Info("Logout: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Logout user (inline)
	if err := h.userRepo.ClearRefreshToken(userID); err != nil {
		log.Errorf("Logout: Failed to logout user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to logout user",
			"errors": fiber.Map{
				"general": err.Error(),
			},
			"data": nil,
		})
	}

	// Clear cookies (if any were set previously)
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   -1,
	})

	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.Domain,
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   -1,
	})

	log.Infof("Logout: User logged out successfully: %s", userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) SendPin(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&input); err != nil || input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
		})
	}

	u := h.userRepo.FindByEmail(input.Email)
	if u == nil || u.ID.IsZero() {
		newUser := NewUser(input.Email)
		if created := h.userRepo.Create(newUser); created == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": ErrInternalServer})
		}
	}

	pin := generateSixDigitPIN()
	if pin == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": ErrInternalServer})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": ErrInternalServer})
	}

	if err := h.userRepo.SetPinWithExpiry(input.Email, string(hash)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": ErrInternalServer})
	}

	if h.cfg.Environment == "development" {
		log.Infof("Login PIN for %s: %s", input.Email, pin)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "PIN sent"})
	}

	from := h.cfg.EmailFrom
	to := []string{input.Email}
	subject := "Your Login PIN"
	body := fmt.Sprintf("Your one-time PIN is: %s. It expires in 10 minutes.", pin)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, input.Email, subject, body)
	auth := smtp.PlainAuth("", h.cfg.SMTPUsername, h.cfg.SMTPPassword, h.cfg.SMTPHost)
	if err := smtp.SendMail(h.cfg.SMTPHost+":"+h.cfg.SMTPPort, auth, from, to, []byte(message)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": ErrInternalServer})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "PIN sent"})
}

func (h *UserHandler) CheckEmail(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&input); err != nil || input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
		})
	}
	user := h.userRepo.FindByEmail(input.Email)
	exists := user != nil && !user.ID.IsZero()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email checked successfully",
		"errors":  nil,
		"data":    map[string]interface{}{"exists": exists},
	})
}
