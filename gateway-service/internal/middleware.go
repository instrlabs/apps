package internal

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
)

func SetupMiddleware(app *fiber.App) {
	// INBOUND LOG
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.WithFields(log.Fields{
			"method":      c.Method(),
			"path":        c.Path(),
			"remote_addr": c.IP(),
			"user_agent":  c.Get("User-Agent"),
			"duration_ms": time.Since(start).Milliseconds(),
		}).Info("Request processed")
		return err
	})

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ALLOWED_ORIGINS"),
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Content-Type, Authorization, X-Authenticated, X-Auth-Token, X-User-Id, X-User-Roles, X-Token-Expires-At",
		AllowCredentials: true,
	}))

	// AUTH
	app.Use(func(c *fiber.Ctx) error {
		token := c.Cookies("access_token")

		c.Request().Header.Del("Cookie")
		c.Request().Header.Del("X-Authenticated")
		c.Request().Header.Del("X-User-Id")
		c.Request().Header.Del("X-User-Roles")

		if token != "" {
			c.Request().Header.Set("X-Auth-Token", token)
			if info, err := ExtractTokenInfo(token); err == nil {
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
