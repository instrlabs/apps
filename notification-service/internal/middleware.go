package internal

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// SetupMiddleware sets up middleware with cookie-to-header conversion for shared RefreshTokenIfNeeded
func SetupMiddleware(app *fiber.App, cfg *Config) {
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Duration(60) * time.Second,
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, OPTIONS",
		AllowHeaders:     "content-type, cookie, authorization",
		AllowCredentials: true,
	}))

	// Refreshed token
	app.Use(func(c *fiber.Ctx) error {
		var accessToken, refreshToken string

		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				accessToken = parts[1]
			}
		}

		refreshToken = c.Get("x-user-refresh")

		needsRefresh := false
		if accessToken == "" {
			needsRefresh = true
		} else {
			_, err := ExtractTokenInfo(cfg.JWTSecret, accessToken)
			if err != nil && (errors.Is(err, ErrTokenExpired) || errors.Is(err, ErrTokenInvalid)) {
				needsRefresh = true
			}
		}

		c.Request().Header.Set("x-user-id", "")
		c.Request().Header.Set("x-user-refresh-token", "")

		if accessToken != "" {
			if info, err := ExtractTokenInfo(cfg.JWTSecret, accessToken); err == nil {
				c.Request().Header.Set("x-user-id", info.UserID)
			}
		}

		if needsRefresh && refreshToken != "" {
			c.Request().Header.Set("x-user-refresh-token", refreshToken)
		}

		return c.Next()
	})
}
