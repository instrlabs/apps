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

	// Global rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Duration(60) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			log.Warnf("Global rate limit exceeded for IP: %s", c.IP())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Rate limit exceeded",
				"errors":  nil,
				"data":    nil,
			})
		},
	}))

	// Enhanced rate limiting for auth endpoints
	app.Use("/auth", limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Duration(60) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			log.Warnf("Auth rate limit exceeded for IP: %s, Path: %s", c.IP(), c.Path())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many authentication attempts. Please try again later.",
				"errors":  nil,
				"data":    nil,
			})
		},
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "content-type, cookie",
		AllowCredentials: true,
	}))

	// CSRF protection
	app.Use(func(c *fiber.Ctx) error {
		method := c.Method()
		if method == fiber.MethodOptions || method == fiber.MethodGet {
			return c.Next()
		}

		if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch || method == fiber.MethodDelete {
			origin := c.Get("x-user-origin")
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

	// AUTH with enhanced security logging
	app.Use(func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		refreshToken := c.Cookies("refresh_token")
		requestPath := c.Path()
		clientIP := c.IP()
		userAgent := c.Get("User-Agent")
		userOrigin := c.Get("x-user-origin")

		c.Request().Header.Del("cookie")
		c.Request().Header.Del("x-authenticated")
		c.Request().Header.Del("x-user-id")
		c.Request().Header.Del("x-user-roles")

		if accessToken != "" {
			if info, err := ExtractTokenInfo(cfg.JWTSecret, accessToken); err == nil {
				c.Request().Header.Set("x-authenticated", "true")
				c.Request().Header.Set("x-user-id", info.UserID)
				c.Request().Header.Set("x-user-roles", strings.Join(info.Roles, ","))
			} else {
				log.Warnf("Authentication failed - Path: %s, IP: %s, Error: %v",
					requestPath, clientIP, err)
				if !errors.Is(err, ErrTokenEmpty) {
					log.Warnf("Invalid token attempt - Path: %s, IP: %s, User-Agent: %s, Origin: %s",
						requestPath, clientIP, userAgent, userOrigin)
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"message": err.Error(),
						"errors":  nil,
						"data":    nil,
					})
				}
				c.Request().Header.Set("x-authenticated", "false")
			}
		}

		if refreshToken != "" && requestPath == "/auth/refresh" {
			c.Request().Header.Set("x-user-refresh", refreshToken)
			log.Infof("Token refresh attempt - IP: %s, User-Agent: %s", clientIP, userAgent)
		}

		c.Request().Header.Set("x-gateway", "true")

		return c.Next()
	})
}
