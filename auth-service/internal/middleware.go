package internal

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	publicPaths := []string{
		"/health",
		"/swagger",
		"/login",
		"/register",
		"/forgot-password",
		"/reset-password",
		"/google",
	}

	isPublic := func(path string) bool {
		for _, prefix := range publicPaths {
			if path == prefix || strings.HasPrefix(path, prefix) {
				return true
			}
		}

		return false
	}

	return func(c *fiber.Ctx) error {
		if c.Get("X-Authenticated") == "true" {
			userId := c.Get("X-User-Id")
			if userId != "" {
				c.Locals("UserID", userId)
			}
			roles := c.Get("X-User-Roles")
			if roles != "" {
				c.Locals("Roles", roles)
			}
		}

		if !isPublic(c.Path()) && c.Get("X-Authenticated") == "false" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"errors":  nil,
				"data":    nil,
			})
		}

		return c.Next()
	}
}
