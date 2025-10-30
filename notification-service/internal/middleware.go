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
	app.Use(compress.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Duration(60) * time.Second,
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     "GET, OPTIONS",
		AllowHeaders:     "content-type, cookie",
		AllowCredentials: true,
	}))

	// AUTH
	app.Use(func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		refreshToken := c.Cookies("refresh_token")
		clientIP := c.IP()
		requestPath := c.Path()

		c.Request().Header.Del("cookie")
		c.Request().Header.Del("x-authenticated")
		c.Request().Header.Del("x-user-id")
		c.Request().Header.Del("x-user-roles")

		if accessToken == "" && refreshToken != "" {
			newTokens, err := RefreshAccessToken(refreshToken)
			if err == nil {
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
			} else {
				log.Warnf("Token refresh failed - IP: %s, Path: %s, Error: %v", clientIP, requestPath, err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized",
					"errors":  nil,
					"data":    nil,
				})
			}
		}

		if accessToken != "" {
			if info, err := ExtractTokenInfo(cfg.JWTSecret, accessToken); err == nil {
				c.Request().Header.Set("x-authenticated", "true")
				c.Request().Header.Set("x-user-id", info.UserID)
				c.Request().Header.Set("x-user-roles", strings.Join(info.Roles, ","))
			} else {
				log.Warnf("ExtractTokenInfo: Failed to extract token info: %v", err)
				if !errors.Is(err, ErrTokenEmpty) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"message": err.Error(),
						"errors":  nil,
						"data":    nil,
					})
				}

				c.Request().Header.Set("x-authenticated", "false")
			}
		}

		return c.Next()
	})
}
