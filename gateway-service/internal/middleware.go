package internal

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func isAllowedOrigin(origin, allowlist string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}
	for _, a := range strings.Split(allowlist, ",") {
		if strings.TrimSpace(a) == origin {
			return true
		}
	}
	return false
}

func SetupMiddleware(app *fiber.App, cfg *Config) {
	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Accept, Content-Type, X-Authenticated, X-User-Id, X-User-Roles",
		AllowCredentials: true,
	}))

	// CSRF protection
	app.Use(func(c *fiber.Ctx) error {
		method := c.Method()
		if method == fiber.MethodOptions || method == fiber.MethodGet {
			return c.Next()
		}

		if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch || method == fiber.MethodDelete {
			origin := c.Get("Origin")
			if !isAllowedOrigin(origin, cfg.Origins) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "FORBIDDEN_ORIGIN",
				})
			}
		}

		return c.Next()
	})

	// AUTH
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		token := c.Cookies("AccessToken")

		c.Request().Header.Del("Cookie")
		c.Request().Header.Del("X-Authenticated")
		c.Request().Header.Del("X-User-Id")
		c.Request().Header.Del("X-User-Roles")
		c.Request().Header.Del("X-Origin")

		if token != "" {
			if info, err := ExtractTokenInfo(cfg.JWTSecret, token); err == nil {
				c.Request().Header.Set("X-Authenticated", "true")
				c.Request().Header.Set("X-User-Id", info.UserID)
				c.Request().Header.Set("X-User-Roles", strings.Join(info.Roles, ","))
				c.Request().Header.Set("X-Origin", origin)
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
