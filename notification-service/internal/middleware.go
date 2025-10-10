package internal

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func SetupMiddleware(app *fiber.App, cfg *Config) {
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Duration(60) * time.Second,
	}))
	app.Use(compress.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Cookie, Set-Cookie",
		AllowCredentials: true,
	}))

	// AUTH
	app.Use(func(c *fiber.Ctx) error {
		host := c.Get("X-Forwarded-Host")
		scheme := c.Get("X-Forwarded-Proto")
		origin := scheme + "://" + host
		accessToken := c.Cookies("AccessToken")

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

		return c.Next()
	})
}
