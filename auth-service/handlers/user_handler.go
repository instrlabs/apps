package handlers

import (
	"log"

	"github.com/arthadede/auth-service/constants"
	"github.com/arthadede/auth-service/controllers"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userController *controllers.UserController
	logger         *log.Logger
}

func NewUserHandler(userController *controllers.UserController) *UserHandler {
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
// @Param request body object{email=string,password=string} true "User registration details"
// @Success 201 {object} object{message=string,data=object{email=string}} "User registered successfully"
// @Failure 400 {object} object{message=string} "Invalid request body or validation error"
// @Failure 409 {object} object{message=string} "Email already exists"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	h.logger.Println("Register: Processing registration request")

	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&input); err != nil {
		h.logger.Printf("Register: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Email == "" {
		h.logger.Println("Register: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrEmailRequired,
		})
	}

	if input.Password == "" || len(input.Password) < 6 {
		h.logger.Println("Register: Password is required and must be at least 6 characters")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrPasswordRequired,
		})
	}

	h.logger.Printf("Register: Attempting to register user with email: %s", input.Email)
	user, err := h.userController.RegisterUser(input.Email, input.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			h.logger.Printf("Register: Email already exists: %s", input.Email)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": constants.ErrEmailAlreadyExists,
			})
		}
		h.logger.Printf("Register: Internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	h.logger.Printf("Register: User registered successfully: %s", user.Email)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"email": user.Email,
		},
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
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Email == "" {
		h.logger.Println("Login: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrEmailRequired,
		})
	}

	if input.Password == "" {
		h.logger.Println("Login: Password is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrPasswordRequired,
		})
	}

	h.logger.Printf("Login: Attempting to login user with email: %s", input.Email)
	tokens, err := h.userController.LoginUser(input.Email, input.Password)
	if err != nil {
		h.logger.Printf("Login: Invalid credentials for email: %s, error: %v", input.Email, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constants.ErrInvalidCredentials,
		})
	}

	h.logger.Println("Login: Setting access token cookie")
	// Set access token as HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   h.userController.GetTokenExpiryHours() * 3600,
	})

	h.logger.Println("Login: Setting refresh token cookie")
	// Set refresh token as HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	h.logger.Printf("Login: User logged in successfully: %s", input.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
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
			"message": constants.ErrRefreshTokenRequired,
		})
	}

	h.logger.Println("RefreshToken: Attempting to refresh token")
	tokens, err := h.userController.RefreshToken(refreshToken)
	if err != nil {
		h.logger.Printf("RefreshToken: Invalid token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constants.ErrInvalidToken,
		})
	}

	h.logger.Println("RefreshToken: Setting new access token cookie")
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   h.userController.GetTokenExpiryHours() * 3600,
	})

	h.logger.Println("RefreshToken: Setting new refresh token cookie")
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
	})

	h.logger.Println("RefreshToken: Token refreshed successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed successfully",
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
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Email == "" {
		h.logger.Println("ForgotPassword: Email is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrEmailRequired,
		})
	}

	h.logger.Printf("ForgotPassword: Requesting password reset for email: %s", input.Email)
	err := h.userController.RequestPasswordReset(input.Email)
	if err != nil {
		h.logger.Printf("ForgotPassword: Error requesting password reset: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
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
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Token == "" {
		h.logger.Println("ResetPassword: Token is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidToken,
		})
	}

	if input.NewPassword == "" || len(input.NewPassword) < 6 {
		h.logger.Println("ResetPassword: New password is required and must be at least 6 characters")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrPasswordRequired,
		})
	}

	h.logger.Println("ResetPassword: Attempting to reset password with token")
	err := h.userController.ResetPassword(input.Token, input.NewPassword)
	if err != nil {
		h.logger.Printf("ResetPassword: Invalid token error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidToken,
		})
	}

	h.logger.Println("ResetPassword: Password has been reset successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully",
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
			"message": constants.ErrInvalidToken,
		})
	}

	h.logger.Println("GoogleCallback: Handling Google callback with authorization code")
	tokens, err := h.userController.HandleGoogleCallback(code)
	if err != nil {
		h.logger.Printf("GoogleCallback: Error handling callback: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	h.logger.Println("GoogleCallback: Setting access token cookie")
	// Set access token as HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokens["access_token"],
		HTTPOnly: true,
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		// Expires in 1 hour (or based on your token expiry configuration)
		MaxAge: h.userController.GetTokenExpiryHours() * 3600,
	})

	h.logger.Println("GoogleCallback: Setting refresh token cookie")
	// Set refresh token as HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens["refresh_token"],
		HTTPOnly: true,
		Secure:   h.userController.GetEnvironment() == "production",
		Path:     "/",
		// Refresh tokens typically have longer expiry
		MaxAge: 30 * 24 * 3600, // 30 days
	})

	// Redirect to frontend
	redirectURL := h.userController.GetOAuthRedirectURL()
	if redirectURL == "" {
		redirectURL = "/"
	}

	h.logger.Printf("GoogleCallback: Redirecting to frontend: %s", redirectURL)
	return c.Redirect(redirectURL)
}

// VerifyToken godoc
// @Summary Verify authentication token
// @Description Verify the validity of the access token from HTTP-only cookie and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string,data=object{user=object}} "Token verified successfully"
// @Failure 401 {object} object{message=string} "Missing or invalid access token cookie"
// @Router /verify-token [post]
func (h *UserHandler) VerifyToken(c *fiber.Ctx) error {
	h.logger.Println("VerifyToken: Processing token verification request")

	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		h.logger.Println("VerifyToken: Access token cookie is missing")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access token is required",
		})
	}

	h.logger.Println("VerifyToken: Attempting to verify token")
	user, err := h.userController.VerifyToken(accessToken)
	if err != nil {
		h.logger.Printf("VerifyToken: Invalid token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	h.logger.Printf("VerifyToken: Token verified successfully for user: %s", user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token verified successfully",
		"data": fiber.Map{
			"user": user,
		},
	})
}
