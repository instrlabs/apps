package internal

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetupMiddleware(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		publicPaths := []string{
			"/login",
			"/refresh",
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

		if c.Get("X-Authenticated") == "true" {
			userId := c.Get("X-User-Id")
			c.Locals("UserID", userId)
			roles := c.Get("X-User-Roles")
			c.Locals("Roles", roles)
		}

		if !isPublic(c.Path()) && c.Get("X-Authenticated") == "false" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"errors":  nil,
				"data":    nil,
			})
		}

		return c.Next()
	})
}
