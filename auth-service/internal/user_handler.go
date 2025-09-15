package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	cfg      *Config
	config   *Config
	userRepo *UserRepository
}

func NewUserHandler(cfg *Config, userRepo *UserRepository) *UserHandler {
	return &UserHandler{
		cfg:      cfg,
		config:   cfg,
		userRepo: userRepo,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	log.Println("Register: Processing registration request")

	var input struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Printf("Register: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Name == "" {
		log.Println("Register: Name is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Name is required",
			"errors": []fiber.Map{
				{
					"fieldName":    "name",
					"errorMessage": "Name is required",
				},
			},
			"data": nil,
		})
	}

	if input.Email == "" {
		log.Println("Register: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrEmailRequired,
			"errors": []fiber.Map{
				{
					"fieldName":    "email",
					"errorMessage": ErrEmailRequired,
				},
			},
			"data": nil,
		})
	}

	if input.Password == "" || len(input.Password) < 6 {
		log.Println("Register: Password is required and must be at least 6 characters")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrPasswordRequired,
			"errors": []fiber.Map{
				{
					"fieldName":    "password",
					"errorMessage": ErrPasswordRequired,
				},
			},
			"data": nil,
		})
	}

	log.Printf("Register: Attempting to register user with name: %s, email: %s", input.Name, input.Email)
	user, err := NewUser(input.Name, input.Email, input.Password)
	if err == nil {
		// Try to create user in repository
		err = h.userRepo.Create(user)
	}
	if err != nil {
		if err.Error() == "user with this email already exists" {
			log.Printf("Register: Email already exists: %s", input.Email)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": ErrEmailAlreadyExists,
				"errors": []fiber.Map{
					{
						"fieldName":    "email",
						"errorMessage": ErrEmailAlreadyExists,
					},
				},
				"data": nil,
			})
		}
		log.Printf("Register: Internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Printf("Register: User registered successfully: %s", user.Email)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"errors":  nil,
		"data":    user,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	log.Println("Login: Processing login request")

	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Printf("Login: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Email == "" {
		log.Println("Login: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrEmailRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Password == "" {
		log.Println("Login: Password is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrPasswordRequired,
			"errors": []fiber.Map{
				{
					"fieldName":    "password",
					"errorMessage": ErrPasswordRequired,
				},
			},
			"data": nil,
		})
	}

	log.Printf("Login: Attempting to login user with email: %s", input.Email)
	// Inline LoginUser logic
	user, err := h.userRepo.FindByEmail(input.Email)
	if err == nil && user != nil {
		if user.Password == "" || !user.ComparePassword(input.Password) {
			err = errors.New("invalid email or password")
		}
	}
	if err != nil {
		log.Printf("Login: Invalid credentials for email: %s, error: %v", input.Email, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidCredentials,
			"errors":  nil,
			"data":    nil,
		})
	}
	// Generate access and refresh tokens
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

	log.Println("Login: Setting access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.CookieDomain,
		Name:     "AccessToken",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	log.Println("Login: Setting refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.CookieDomain,
		Name:     "RefreshToken",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	log.Printf("Login: User logged in successfully: %s", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	log.Println("RefreshToken: Processing token refresh request")

	var input struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Printf("RefreshToken: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.RefreshToken == "" {
		log.Println("RefreshToken: Refresh token cookie is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrRefreshTokenRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Println("RefreshToken: Attempting to refresh token")
	// Inline RefreshToken logic
	user, err := h.userRepo.FindByRefreshToken(input.RefreshToken)
	if err != nil {
		log.Printf("RefreshToken: Invalid token error: %v", err)
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

	log.Println("RefreshToken: Setting new access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.cfg.CookieDomain,
		Name:     "AccessToken",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.cfg.Environment == "production",
		Path:     "/",
		MaxAge:   h.cfg.TokenExpiryHours * 3600,
	})

	log.Println("RefreshToken: Setting new refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.config.CookieDomain,
		Name:     "RefreshToken",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	log.Println("RefreshToken: Token refreshed successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	log.Println("ForgotPassword: Processing forgot password request")

	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Printf("ForgotPassword: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Email == "" {
		log.Println("ForgotPassword: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrEmailRequired,
			"errors": []fiber.Map{
				{
					"fieldName":    "email",
					"errorMessage": ErrEmailRequired,
				},
			},
			"data": nil,
		})
	}

	log.Printf("ForgotPassword: Requesting password reset for email: %s", input.Email)
	// Inline RequestPasswordReset logic
	if _, err := h.userRepo.FindByEmail(input.Email); err != nil {
		// Do not reveal whether email exists
		log.Printf("ForgotPassword: Email not found or other error: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If your email is registered, you will receive a password reset link",
		})
	}
	// Generate reset token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	resetToken := base64.URLEncoding.EncodeToString(b)
	expiry := time.Now().UTC().Add(time.Hour * time.Duration(h.config.ResetTokenExpiryHours))
	if err := h.userRepo.SetResetToken(input.Email, resetToken, expiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	// Send email or log in dev
	if h.config.Environment == "development" {
		log.Printf("Password reset token for %s: %s", input.Email, resetToken)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If your email is registered, you will receive a password reset link",
		})
	}
	from := h.config.EmailFrom
	to := []string{input.Email}
	resetURL := fmt.Sprintf("%s?token=%s", h.config.FEResetPassword, resetToken)
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click the link below to reset your password:\n\n%s\n\nIf you did not request a password reset, please ignore this email.", resetURL)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, input.Email, subject, body)
	auth := smtp.PlainAuth("", h.config.SMTPUsername, h.config.SMTPPassword, h.config.SMTPHost)
	if err := smtp.SendMail(h.config.SMTPHost+":"+h.config.SMTPPort, auth, from, to, []byte(message)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	// fallthrough to success response
	// continue to success response

	log.Printf("ForgotPassword: Password reset requested for email: %s", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If your email is registered, you will receive a password reset link",
	})
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	log.Println("ResetPassword: Processing password reset request")

	var input struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Printf("ResetPassword: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Token == "" {
		log.Println("ResetPassword: Token is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.NewPassword == "" || len(input.NewPassword) < 6 {
		log.Println("ResetPassword: New password is required and must be at least 6 characters")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrPasswordRequired,
			"errors": []fiber.Map{
				{
					"fieldName":    "new_password",
					"errorMessage": ErrPasswordRequired,
				},
			},
			"data": nil,
		})
	}

	log.Println("ResetPassword: Attempting to reset password with token")
	// Inline ResetPassword logic
	user, err := h.userRepo.FindByResetToken(input.Token)
	if err != nil {
		log.Printf("ResetPassword: Invalid token error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.userRepo.UpdatePassword(user.ID.Hex(), string(hashedPassword)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Println("ResetPassword: Password has been reset successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) GoogleLogin(c *fiber.Ctx) error {
	log.Println("GoogleLogin: Initiating Google OAuth login")

	// Build oauth2 config locally (no field stored)
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
	log.Printf("GoogleLogin: Redirecting to Google OAuth URL: %s", url)

	return c.Redirect(url)
}

func (h *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	log.Println("GoogleCallback: Processing Google OAuth callback")

	code := c.Query("code")
	if code == "" {
		log.Println("GoogleCallback: Missing authorization code")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Println("GoogleCallback: Handling Google callback with authorization code")
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
		log.Printf("GoogleCallback: Exchange error: %v", err)
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
	user, err := h.userRepo.FindByGoogleID(userInfo.ID)
	if err != nil {
		user, err = h.userRepo.FindByEmail(userInfo.Email)
		if err != nil {
			user = NewGoogleUser(userInfo.Name, userInfo.Email, userInfo.ID)
			if err := h.userRepo.Create(user); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": ErrInternalServer,
					"errors":  nil,
					"data":    nil,
				})
			}
		} else {
			if err := h.userRepo.UpdateGoogleID(user.ID.Hex(), userInfo.ID); err != nil {
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
	accessToken, err := tok.SignedString([]byte(h.config.JWTSecret))
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

	log.Println("GoogleCallback: Setting access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.config.CookieDomain,
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   h.config.TokenExpiryHours * 3600,
	})

	log.Println("GoogleCallback: Setting refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.config.CookieDomain,
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	redirectURL := h.config.FEOAuthRedirect
	if redirectURL == "" {
		redirectURL = "/"
	}

	log.Printf("GoogleCallback: Redirecting to frontend: %s", redirectURL)
	return c.Redirect(redirectURL)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	log.Println("GetProfile: Processing profile request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		log.Println("GetProfile: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("GetProfile: User not found for UserID %s: %v", userID, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUserNotFound,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Printf("GetProfile: Profile retrieved successfully for user: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile retrieved successfully",
		"errors":  nil,
		"data":    user,
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	log.Println("UpdateProfile: Processing profile update request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		log.Println("UpdateProfile: UserID not found in context")
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
		log.Printf("UpdateProfile: Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Update profile (inline)
	if request.Name == "" {
		log.Printf("UpdateProfile: name cannot be empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update profile",
			"errors":  fiber.Map{"general": "name cannot be empty"},
			"data":    nil,
		})
	}
	if err := h.userRepo.UpdateProfile(userID, request.Name); err != nil {
		log.Printf("UpdateProfile: Failed to update profile: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update profile",
			"errors": fiber.Map{
				"general": err.Error(),
			},
			"data": nil,
		})
	}

	// Get updated user
	updatedUser, err := h.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("UpdateProfile: Failed to get updated user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	log.Printf("UpdateProfile: Profile updated successfully for user: %s", updatedUser.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"errors":  nil,
		"data": fiber.Map{
			"user": updatedUser,
		},
	})
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	log.Println("ChangePassword: Processing password change request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		log.Println("ChangePassword: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Parse request body
	var request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&request); err != nil {
		log.Printf("ChangePassword: Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Change password (inline)
	if request.CurrentPassword == "" || request.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to change password",
			"errors":  fiber.Map{"general": "passwords cannot be empty"},
			"data":    nil,
		})
	}
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to change password",
			"errors":  fiber.Map{"general": err.Error()},
			"data":    nil,
		})
	}
	if !user.ComparePassword(request.CurrentPassword) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to change password",
			"errors":  fiber.Map{"general": "current password is incorrect"},
			"data":    nil,
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to change password",
			"errors":  fiber.Map{"general": err.Error()},
			"data":    nil,
		})
	}

	log.Printf("ChangePassword: Password changed successfully for user: %s", userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	log.Println("Logout: Processing logout request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		log.Println("Logout: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Logout user (inline)
	if err := h.userRepo.ClearRefreshToken(userID); err != nil {
		log.Printf("Logout: Failed to logout user: %v", err)
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
		Domain:   h.config.CookieDomain,
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   -1,
	})

	c.Cookie(&fiber.Cookie{
		Domain:   h.config.CookieDomain,
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   -1,
	})

	log.Printf("Logout: User logged out successfully: %s", userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
		"errors":  nil,
		"data":    nil,
	})
}
