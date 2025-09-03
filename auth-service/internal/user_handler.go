package internal

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userController *UserController
	config         *Config
	logger         *log.Logger
}

func NewUserHandler(userController *UserController, config *Config) *UserHandler {
	return &UserHandler{
		userController: userController,
		config:         config,
		logger:         log.New(log.Writer(), "[UserHandler] ", log.LstdFlags|log.Lshortfile),
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	h.logger.Println("Register: Processing registration request")

	var input struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&input); err != nil {
		h.logger.Printf("Register: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Name == "" {
		h.logger.Println("Register: Name is required")
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
		h.logger.Println("Register: Email is required")
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
		h.logger.Println("Register: Password is required and must be at least 6 characters")
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

	h.logger.Printf("Register: Attempting to register user with name: %s, email: %s", input.Name, input.Email)
	user, err := h.userController.RegisterUser(input.Name, input.Email, input.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			h.logger.Printf("Register: Email already exists: %s", input.Email)
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
		h.logger.Printf("Register: Internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Printf("Register: User registered successfully: %s", user.Email)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"errors":  nil,
		"data":    user,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	h.logger.Println("Login: Processing login request")

	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		h.logger.Printf("Login: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Email == "" {
		h.logger.Println("Login: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrEmailRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Password == "" {
		h.logger.Println("Login: Password is required")
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

	h.logger.Printf("Login: Attempting to login user with email: %s", input.Email)
	tokens, err := h.userController.LoginUser(input.Email, input.Password)
	if err != nil {
		h.logger.Printf("Login: Invalid credentials for email: %s, error: %v", input.Email, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidCredentials,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("Login: Setting access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.config.CookieDomain,
		Name:     "AccessToken",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   h.config.TokenExpiryHours * 3600,
	})

	h.logger.Println("Login: Setting refresh token cookie")
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

	h.logger.Printf("Login: User logged in successfully: %s", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	h.logger.Println("RefreshToken: Processing token refresh request")

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		h.logger.Println("RefreshToken: Refresh token cookie is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrRefreshTokenRequired,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("RefreshToken: Attempting to refresh token")
	tokens, err := h.userController.RefreshToken(refreshToken)
	if err != nil {
		h.logger.Printf("RefreshToken: Invalid token error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("RefreshToken: Setting new access token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.config.CookieDomain,
		Name:     "AccessToken",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.config.Environment == "production",
		Path:     "/",
		MaxAge:   h.config.TokenExpiryHours * 3600,
	})

	h.logger.Println("RefreshToken: Setting new refresh token cookie")
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

	h.logger.Println("RefreshToken: Token refreshed successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	h.logger.Println("ForgotPassword: Processing forgot password request")

	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&input); err != nil {
		h.logger.Printf("ForgotPassword: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Email == "" {
		h.logger.Println("ForgotPassword: Email is required")
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

	h.logger.Printf("ForgotPassword: Requesting password reset for email: %s", input.Email)
	err := h.userController.RequestPasswordReset(input.Email)
	if err != nil {
		h.logger.Printf("ForgotPassword: Error requesting password reset: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Printf("ForgotPassword: Password reset requested for email: %s", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If your email is registered, you will receive a password reset link",
	})
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	h.logger.Println("ResetPassword: Processing password reset request")

	var input struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&input); err != nil {
		h.logger.Printf("ResetPassword: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.Token == "" {
		h.logger.Println("ResetPassword: Token is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	if input.NewPassword == "" || len(input.NewPassword) < 6 {
		h.logger.Println("ResetPassword: New password is required and must be at least 6 characters")
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

	h.logger.Println("ResetPassword: Attempting to reset password with token")
	err := h.userController.ResetPassword(input.Token, input.NewPassword)
	if err != nil {
		h.logger.Printf("ResetPassword: Invalid token error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("ResetPassword: Password has been reset successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) GoogleLogin(c *fiber.Ctx) error {
	h.logger.Println("GoogleLogin: Initiating Google OAuth login")

	url := h.userController.GetGoogleAuthURL()
	h.logger.Printf("GoogleLogin: Redirecting to Google OAuth URL: %s", url)

	return c.Redirect(url)
}

func (h *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	h.logger.Println("GoogleCallback: Processing Google OAuth callback")

	code := c.Query("code")
	if code == "" {
		h.logger.Println("GoogleCallback: Missing authorization code")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidToken,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("GoogleCallback: Handling Google callback with authorization code")
	tokens, err := h.userController.HandleGoogleCallback(code)
	if err != nil {
		h.logger.Printf("GoogleCallback: Error handling callback: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("GoogleCallback: Setting access token cookie")
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

	h.logger.Println("GoogleCallback: Setting refresh token cookie")
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

	h.logger.Printf("GoogleCallback: Redirecting to frontend: %s", redirectURL)
	return c.Redirect(redirectURL)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	h.logger.Println("GetProfile: Processing profile request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		h.logger.Println("GetProfile: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	user, err := h.userController.userRepo.FindByID(userID)
	if err != nil {
		h.logger.Printf("GetProfile: User not found for UserID %s: %v", userID, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUserNotFound,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Printf("GetProfile: Profile retrieved successfully for user: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile retrieved successfully",
		"errors":  nil,
		"data": fiber.Map{
			"user": user,
		},
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	h.logger.Println("UpdateProfile: Processing profile update request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		h.logger.Println("UpdateProfile: UserID not found in context")
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
		h.logger.Printf("UpdateProfile: Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Update profile
	if err := h.userController.UpdateProfile(userID, request.Name); err != nil {
		h.logger.Printf("UpdateProfile: Failed to update profile: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update profile",
			"errors": fiber.Map{
				"general": err.Error(),
			},
			"data": nil,
		})
	}

	// Get updated user
	updatedUser, err := h.userController.userRepo.FindByID(userID)
	if err != nil {
		h.logger.Printf("UpdateProfile: Failed to get updated user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServer,
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Printf("UpdateProfile: Profile updated successfully for user: %s", updatedUser.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"errors":  nil,
		"data": fiber.Map{
			"user": updatedUser,
		},
	})
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	h.logger.Println("ChangePassword: Processing password change request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		h.logger.Println("ChangePassword: UserID not found in context")
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
		h.logger.Printf("ChangePassword: Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidRequestBody,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Change password
	if err := h.userController.ChangePassword(userID, request.CurrentPassword, request.NewPassword); err != nil {
		h.logger.Printf("ChangePassword: Failed to change password: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to change password",
			"errors": fiber.Map{
				"general": err.Error(),
			},
			"data": nil,
		})
	}

	h.logger.Printf("ChangePassword: Password changed successfully for user: %s", userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	h.logger.Println("Logout: Processing logout request using Locals UserID")

	userID, _ := c.Locals("UserID").(string)
	if userID == "" {
		h.logger.Println("Logout: UserID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized,
			"errors":  nil,
			"data":    nil,
		})
	}

	// Logout user
	if err := h.userController.LogoutUser(userID); err != nil {
		h.logger.Printf("Logout: Failed to logout user: %v", err)
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

	h.logger.Printf("Logout: User logged out successfully: %s", userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
		"errors":  nil,
		"data":    nil,
	})
}
