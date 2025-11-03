package internal

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	ErrForbiddenOrigin = errors.New("FORBIDDEN_ORIGIN")
)

func isAllowedOrigin(origin, allowlist string) bool {
	origin = strings.TrimSpace(origin)
	for _, a := range strings.Split(allowlist, ",") {
		if strings.TrimSpace(a) == origin {
			return true
		}
	}
	return false
}

func SetupMiddleware(app *fiber.App, cfg *Config) {
	app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(etag.New())
	app.Use(compress.New())

	// Rate limiter
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Duration(60) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			log.Warnf("Rate Limit exceeded for IP: %s", c.IP())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Rate limit exceeded",
				"errors":  nil,
				"data":    nil,
			})
		},
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "content-type, cookie, authorization",
		AllowCredentials: true,
	}))

	// CSRF protection
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("x-user-origin")
		if !isAllowedOrigin(origin, cfg.Origins) && cfg.CSRFEnabled {
			log.Warnf("CSRF protection: Forbidden origin: %s", origin)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": ErrForbiddenOrigin.Error(),
				"errors":  nil,
				"data":    nil,
			})
		}

		return c.Next()
	})

	// JWT token extraction and user authentication
	app.Use(func(c *fiber.Ctx) error {
		var accessToken string

		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				accessToken = parts[1]
			}
		}

		c.Request().Header.Set("x-user-id", "")

		if accessToken != "" {
			if info, err := ExtractTokenInfo(cfg.JWTSecret, accessToken); err == nil {
				log.Infof("Authenticated user %s for %s %s", info.UserID, c.Method(), c.Path())
				c.Request().Header.Set("x-user-id", info.UserID)
			} else {
				log.Warnf("Token extraction failed for %s %s: %v", c.Method(), c.Path(), err)
			}
		}

		return c.Next()
	})
}
