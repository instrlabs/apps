package internal

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userController *UserController
	logger         *log.Logger
}

func NewUserHandler(userController *UserController) *UserHandler {
	return &UserHandler{
		userController: userController,
		logger:         log.New(log.Writer(), "[UserHandler] ", log.LstdFlags|log.Lshortfile),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{name=string,email=string,password=string} true "User registration details"
// @Success 201 {object} object{message=string,data=object{name=string,email=string}} "User registered successfully"
// @Failure 400 {object} object{message=string} "Invalid request body or validation error"
// @Failure 409 {object} object{message=string} "Email already exists"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /register [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticate a user with email and password and set HTTP-only cookies for tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "User login credentials"
// @Success 200 {object} object{message=string} "Login successful with tokens set as HTTP-only cookies"
// @Failure 400 {object} object{message=string} "Invalid request body or validation error"
// @Failure 401 {object} object{message=string} "Invalid credentials"
// @Router /login [post]
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
		Domain:   h.userController.GetCookieDomain(),
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   h.userController.GetTokenExpiryHours() * 3600,
	})

	h.logger.Println("Login: Setting refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.userController.GetCookieDomain(),
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using the refresh token from HTTP-only cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string} "Token refreshed successfully with new tokens set as HTTP-only cookies"
// @Failure 400 {object} object{message=string} "Missing refresh token cookie"
// @Failure 401 {object} object{message=string} "Invalid token"
// @Router /refresh [post]
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
		Domain:   h.userController.GetCookieDomain(),
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   h.userController.GetTokenExpiryHours() * 3600,
	})

	h.logger.Println("RefreshToken: Setting new refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.userController.GetCookieDomain(),
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
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

// ForgotPassword godoc
// @Summary Request password reset
// @Description Request a password reset link for a registered email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string} true "User email"
// @Success 200 {object} object{message=string} "Password reset link sent if email exists"
// @Failure 400 {object} object{message=string} "Invalid request body or validation error"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /forgot-password [post]
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

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset a user's password using a reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{token=string,new_password=string} true "Reset token and new password"
// @Success 200 {object} object{message=string} "Password has been reset successfully"
// @Failure 400 {object} object{message=string} "Invalid request body, token, or password validation error"
// @Router /reset-password [post]
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

// GoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirect user to Google OAuth consent screen
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirect to Google OAuth consent screen"
// @Router /google [get]
func (h *UserHandler) GoogleLogin(c *fiber.Ctx) error {
	h.logger.Println("GoogleLogin: Initiating Google OAuth login")

	url := h.userController.GetGoogleAuthURL()
	h.logger.Printf("GoogleLogin: Redirecting to Google OAuth URL: %s", url)

	return c.Redirect(url)
}

// GoogleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Process the OAuth code from Google, set HTTP-only cookies with tokens, and redirect to frontend
// @Tags auth
// @Produce json
// @Param code query string true "OAuth authorization code"
// @Success 302 {string} string "Redirect to frontend with tokens set as HTTP-only cookies"
// @Failure 400 {object} object{message=string} "Invalid or missing code parameter"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /google/callback [get]
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
		Domain:   h.userController.GetCookieDomain(),
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   h.userController.GetTokenExpiryHours() * 3600,
	})

	h.logger.Println("GoogleCallback: Setting refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Domain:   h.userController.GetCookieDomain(),
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	redirectURL := h.userController.GetOAuthRedirectURL()
	if redirectURL == "" {
		redirectURL = "/"
	}

	h.logger.Printf("GoogleCallback: Redirecting to frontend: %s", redirectURL)
	return c.Redirect(redirectURL)
}

// VerifyToken godoc
// @Summary Verify authentication token
// @Description Verify the validity of the access token from context (set by middleware), X-Auth-Token header (set by gateway), or HTTP-only cookie and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string,data=object{user=object}} "Token verified successfully"
// @Failure 401 {object} object{message=string} "Missing or invalid access token"
// @Router /verify-token [post]
func (h *UserHandler) VerifyToken(c *fiber.Ctx) error {
	h.logger.Println("VerifyToken: Processing token verification request")

	var accessToken string
	if token, ok := c.Locals("token").(string); ok && token != "" {
		h.logger.Println("VerifyToken: Using token from context")
		accessToken = token
	} else {
		h.logger.Println("VerifyToken: Access token not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access token is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("VerifyToken: Attempting to verify token")
	user, err := h.userController.VerifyToken(accessToken)
	if err != nil {
		h.logger.Printf("VerifyToken: Invalid token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Printf("VerifyToken: Token verified successfully for user: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token verified successfully",
		"errors":  nil,
		"data":    nil,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the user profile information based on the authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string,data=object{user=User}} "Profile retrieved successfully"
// @Failure 401 {object} object{message=string} "Missing or invalid access token"
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	h.logger.Println("GetProfile: Processing profile request")

	var accessToken string
	if token, ok := c.Locals("token").(string); ok && token != "" {
		h.logger.Println("GetProfile: Using token from context")
		accessToken = token
	} else {
		h.logger.Println("GetProfile: Access token not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access token is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	h.logger.Println("GetProfile: Attempting to verify token and get user profile")
	user, err := h.userController.VerifyToken(accessToken)
	if err != nil {
		h.logger.Printf("GetProfile: Invalid token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
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

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the user profile information (currently only name)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{name=string} true "Profile update request"
// @Success 200 {object} object{message=string,data=object{user=User}} "Profile updated successfully"
// @Failure 400 {object} object{message=string,errors=object} "Invalid request"
// @Failure 401 {object} object{message=string} "Missing or invalid access token"
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	h.logger.Println("UpdateProfile: Processing profile update request")

	// Get token from context (set by middleware)
	var accessToken string
	if token, ok := c.Locals("token").(string); ok && token != "" {
		h.logger.Println("UpdateProfile: Using token from context")
		accessToken = token
	} else {
		h.logger.Println("UpdateProfile: Access token not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access token is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Verify token and get user
	user, err := h.userController.VerifyToken(accessToken)
	if err != nil {
		h.logger.Printf("UpdateProfile: Invalid token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
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
			"message": "Invalid request body",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Update profile
	err = h.userController.UpdateProfile(user.ID.Hex(), request.Name)
	if err != nil {
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
	updatedUser, err := h.userController.VerifyToken(accessToken)
	if err != nil {
		h.logger.Printf("UpdateProfile: Failed to get updated user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Profile updated but failed to retrieve updated data",
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

// ChangePassword godoc
// @Summary Change user password
// @Description Change the user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{current_password=string,new_password=string} true "Password change request"
// @Success 200 {object} object{message=string} "Password changed successfully"
// @Failure 400 {object} object{message=string,errors=object} "Invalid request"
// @Failure 401 {object} object{message=string} "Missing or invalid access token"
// @Router /change-password [post]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	h.logger.Println("ChangePassword: Processing password change request")

	// Get token from context (set by middleware)
	var accessToken string
	if token, ok := c.Locals("token").(string); ok && token != "" {
		h.logger.Println("ChangePassword: Using token from context")
		accessToken = token
	} else {
		h.logger.Println("ChangePassword: Access token not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access token is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Verify token and get user
	user, err := h.userController.VerifyToken(accessToken)
	if err != nil {
		h.logger.Printf("ChangePassword: Invalid token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
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
			"message": "Invalid request body",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Change password
	err = h.userController.ChangePassword(user.ID.Hex(), request.CurrentPassword, request.NewPassword)
	if err != nil {
		h.logger.Printf("ChangePassword: Failed to change password: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to change password",
			"errors": fiber.Map{
				"general": err.Error(),
			},
			"data": nil,
		})
	}

	h.logger.Printf("ChangePassword: Password changed successfully for user: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
		"errors":  nil,
		"data":    nil,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout a user by clearing their refresh token and cookies
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string} "Logout successful"
// @Failure 401 {object} object{message=string} "Missing or invalid access token"
// @Router /logout [post]
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	h.logger.Println("Logout: Processing logout request")

	// Get token from context (set by middleware)
	var accessToken string
	if token, ok := c.Locals("token").(string); ok && token != "" {
		h.logger.Println("Logout: Using token from context")
		accessToken = token
	} else {
		h.logger.Println("Logout: Access token not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access token is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Verify token and get user
	user, err := h.userController.VerifyToken(accessToken)
	if err != nil {
		h.logger.Printf("Logout: Invalid token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Logout user
	err = h.userController.LogoutUser(user.ID.Hex())
	if err != nil {
		h.logger.Printf("Logout: Failed to logout user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to logout user",
			"errors": fiber.Map{
				"general": err.Error(),
			},
			"data": nil,
		})
	}

	// Clear cookies
	c.Cookie(&fiber.Cookie{
		Domain:   h.userController.GetCookieDomain(),
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   -1,
	})

	c.Cookie(&fiber.Cookie{
		Domain:   h.userController.GetCookieDomain(),
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   -1,
	})

	h.logger.Printf("Logout: User logged out successfully: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
		"errors":  nil,
		"data":    nil,
	})
}
