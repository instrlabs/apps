package internal

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupMiddleware(app *fiber.App) {
	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ALLOWED_ORIGINS"),
		AllowMethods:     "GET",
		AllowHeaders:     "Origin",
		AllowCredentials: true,
	}))

	// AUTH
	app.Use(func(c *fiber.Ctx) error {
		token := c.Cookies("AccessToken")

		c.Request().Header.Del("Cookie")
		c.Request().Header.Del("X-Authenticated")
		c.Request().Header.Del("X-User-Id")
		c.Request().Header.Del("X-User-Roles")

		if token != "" {
			if info, err := ExtractTokenInfo(token); err == nil {
				c.Locals("UserID", info.UserID)
				c.Locals("Roles", info.Roles)
			} else {
				if errors.Is(err, ErrTokenExpired) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"message": "ExpiredToken",
						"errors":  nil,
						"data":    nil,
					})
				}

				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
					"errors":  nil,
					"data":    nil,
				})
			}
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"errors":  nil,
				"data":    nil,
			})
		}

		return c.Next()
	})
}
