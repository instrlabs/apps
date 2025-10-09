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

var ErrForbiddenOrigin = errors.New("FORBIDDEN_ORIGIN")

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
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Duration(60) * time.Second,
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Content-Type, Cookie",
		AllowCredentials: true,
	}))

	// CSRF protection
	app.Use(func(c *fiber.Ctx) error {
		method := c.Method()
		if method == fiber.MethodOptions || method == fiber.MethodGet {
			return c.Next()
		}

		if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch || method == fiber.MethodDelete {
			host := c.Get("X-Forwarded-Host")
			scheme := c.Get("X-Forwarded-Proto")
			origin := scheme + "://" + host
			if !isAllowedOrigin(origin, cfg.Origins) && cfg.CSRFEnabled {
				log.Warnf("CSRF protection: Forbidden origin: %s", origin)
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"message": ErrForbiddenOrigin.Error(),
					"errors":  nil,
					"data":    nil,
				})
			}
		}

		return c.Next()
	})

	// AUTH
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		accessToken := c.Cookies("AccessToken")
		refreshToken := c.Cookies("RefreshToken")

		c.Request().Header.Del("Cookie")
		c.Request().Header.Del("X-Authenticated")
		c.Request().Header.Del("X-User-Id")
		c.Request().Header.Del("X-User-Roles")
		c.Request().Header.Del("X-Origin")

		if accessToken != "" {
			if info, err := ExtractTokenInfo(cfg.JWTSecret, accessToken); err == nil {
				c.Request().Header.Set("X-Authenticated", "true")
				c.Request().Header.Set("X-User-Id", info.UserID)
				c.Request().Header.Set("X-User-Roles", strings.Join(info.Roles, ","))
				c.Request().Header.Set("X-Origin", origin)
			} else {
				log.Warnf("ExtractTokenInfo: Failed to extract token info: %v", err)
				if !errors.Is(err, ErrTokenEmpty) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"message": err.Error(),
						"errors":  nil,
						"data":    nil,
					})
				}

				c.Request().Header.Set("X-Authenticated", "false")
			}
		}

		if refreshToken != "" && c.Path() == "/auth/refresh" {
			c.Request().Header.Set("X-User-Refresh", refreshToken)
		}

		c.Request().Header.Set("X-Gateway", "true")

		return c.Next()
	})
}
