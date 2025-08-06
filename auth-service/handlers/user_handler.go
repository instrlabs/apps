package handlers

import (
	"github.com/arthadede/auth-service/constants"
	"github.com/arthadede/auth-service/controllers"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userController *controllers.UserController
}

func NewUserHandler(userController *controllers.UserController) *UserHandler {
	return &UserHandler{
		userController: userController,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrEmailRequired,
		})
	}

	if input.Password == "" || len(input.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrPasswordRequired,
		})
	}

	user, err := h.userController.RegisterUser(input.Email, input.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": constants.ErrEmailAlreadyExists,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"email": user.Email,
		},
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrEmailRequired,
		})
	}

	if input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrPasswordRequired,
		})
	}

	tokens, err := h.userController.LoginUser(input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constants.ErrInvalidCredentials,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"access_token":  tokens["access_token"],
			"refresh_token": tokens["refresh_token"],
		},
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrRefreshTokenRequired,
		})
	}

	tokens, err := h.userController.RefreshToken(input.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constants.ErrInvalidToken,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"data": fiber.Map{
			"access_token":  tokens["access_token"],
			"refresh_token": tokens["refresh_token"],
		},
	})
}

func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrEmailRequired,
		})
	}

	err := h.userController.RequestPasswordReset(input.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If your email is registered, you will receive a password reset link",
	})
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var input struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidToken,
		})
	}

	if input.NewPassword == "" || len(input.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrPasswordRequired,
		})
	}

	err := h.userController.ResetPassword(input.Token, input.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidToken,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully",
	})
}

func (h *UserHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.userController.GetGoogleAuthURL()

	return c.Redirect(url)
}

func (h *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	var input struct {
		Code string `json:"code" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidToken,
		})
	}

	tokens, err := h.userController.HandleGoogleCallback(input.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Google login successful",
		"data": fiber.Map{
			"access_token":  tokens["access_token"],
			"refresh_token": tokens["refresh_token"],
		},
	})
}

func (h *UserHandler) VerifyToken(c *fiber.Ctx) error {
	var input struct {
		Token string `json:"token" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidRequestBody,
		})
	}

	if input.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Token is required",
		})
	}

	user, err := h.userController.VerifyToken(input.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token verified successfully",
		"data": fiber.Map{
			"user": user,
		},
	})
}
