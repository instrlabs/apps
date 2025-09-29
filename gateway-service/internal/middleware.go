package internal

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupMiddleware(app *fiber.App, cfg *Config) {
	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Accept, Content-Type, X-Authenticated, X-User-Id, X-User-Roles",
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
			if info, err := ExtractTokenInfo(cfg.JWTSecret, token); err == nil {
				c.Request().Header.Set("X-Authenticated", "true")
				c.Request().Header.Set("X-User-Id", info.UserID)
				c.Request().Header.Set("X-User-Roles", strings.Join(info.Roles, ","))
			} else {
				if errors.Is(err, ErrTokenExpired) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "EXPIRED_TOKEN",
					})
				}

				c.Request().Header.Set("X-Authenticated", "false")
			}
		}

		c.Request().Header.Set("X-Gateway", "true")

		return c.Next()
	})
}
