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
		AllowHeaders:     "content-type, cookie",
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

	// Refreshed token
	app.Use(func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		refreshToken := c.Cookies("refresh_token")
		clientIP := c.Get("x-user-ip")
		requestPath := c.Path()

		c.Request().Header.Del("cookie")
		c.Request().Header.Del("x-authenticated")
		c.Request().Header.Del("x-user-id")
		c.Request().Header.Del("x-user-roles")

		if accessToken == "" && refreshToken != "" {
			newTokens, err := RefreshAccessToken(refreshToken)
			if err != nil {
				log.Warnf("Token refresh failed - IP: %s, Path: %s, Error: %v", clientIP, requestPath, err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
					"errors":  nil,
					"data":    nil,
				})
			}

			c.Cookie(&fiber.Cookie{
				Name:     "access_token",
				Value:    newTokens.AccessToken,
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Path:     "/",
			})
			c.Cookie(&fiber.Cookie{
				Name:     "refresh_token",
				Value:    newTokens.RefreshToken,
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Path:     "/",
			})

			accessToken = newTokens.AccessToken
			log.Infof("Token refreshed - IP: %s, Path: %s", clientIP, requestPath)
		}

		if accessToken != "" {
			info, err := ExtractTokenInfo(cfg.JWTSecret, accessToken)
			if err != nil {
				c.Request().Header.Set("x-authenticated", "false")
				c.Request().Header.Set("x-user-id", "")
			} else {
				c.Request().Header.Set("x-authenticated", "true")
				c.Request().Header.Set("x-user-id", info.UserID)
			}
		}

		return c.Next()
	})
}
