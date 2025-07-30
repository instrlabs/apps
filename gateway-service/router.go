package main

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	log "github.com/sirupsen/logrus"
)

func SetupGatewayRoutes(app *fiber.App, config Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "ok"})
	})

	app.Get("/check-auth", func(c *fiber.Ctx) error {
		// Find auth-service configuration
		var authServiceURL string
		for _, service := range config.Services {
			if service.Name == "auth-service" {
				authServiceURL = service.URL
				break
			}
		}

		if authServiceURL == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{
				"status":  "error",
				"message": "Auth service not configured",
			})
		}

		// Create a new Fiber agent for making HTTP requests
		agent := fiber.AcquireAgent()
		defer fiber.ReleaseAgent(agent)

		// Set up the request to auth service health endpoint
		req := agent.Request()
		req.SetRequestURI(authServiceURL + "/health")
		req.Header.SetMethod(fiber.MethodGet)

		// Send the request with a timeout
		if err := agent.Parse(); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to parse request to auth service")

			return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{
				"status":  "error",
				"message": "Failed to connect to auth service",
				"error":   err.Error(),
			})
		}

		// Check the response
		code, body, errs := agent.Bytes()
		if len(errs) > 0 {
			log.WithFields(log.Fields{
				"errors": errs,
			}).Error("Failed to connect to auth service")

			return c.Status(fiber.StatusServiceUnavailable).JSON(map[string]string{
				"status":  "error",
				"message": "Auth service is unavailable",
				"error":   errs[0].Error(),
			})
		}

		if code != fiber.StatusOK {
			return c.Status(fiber.StatusServiceUnavailable).JSON(map[string]string{
				"status":  "error",
				"message": "Auth service returned non-OK status",
				"code":    string(code),
			})
		}

		// Return the response from auth service
		return c.Status(code).Send(body)
	})

	for _, service := range config.Services {
		targetURL, err := url.Parse(service.URL)
		if err != nil {
			log.WithFields(log.Fields{
				"service": service.Name,
				"url":     service.URL,
				"error":   err.Error(),
			}).Fatal("Failed to parse service URL")
		}

		for _, prefix := range service.Prefixes {
			app.All(prefix+"*", func(c *fiber.Ctx) error {
				log.WithFields(log.Fields{
					"service":      service.Name,
					"method":       c.Method(),
					"path":         c.Path(),
					"forwarded_to": targetURL.String() + c.Path(),
				}).Info("Forwarding request")

				c.Request().Header.Set("X-Gateway", "true")

				url := targetURL.String() + c.Path()
				if err := proxy.Do(c, url); err != nil {
					log.WithFields(log.Fields{
						"service": service.Name,
						"method":  c.Method(),
						"path":    c.Path(),
						"error":   err.Error(),
					}).Error("Proxy error")

					return c.Status(fiber.StatusBadGateway).JSON(map[string]string{
						"error":   "Bad Gateway",
						"message": "The service is currently unavailable",
					})
				}

				return nil
			})

			log.WithFields(log.Fields{
				"service": service.Name,
				"prefix":  prefix,
				"target":  service.URL,
			}).Info("Registered route")
		}
	}

	app.Use(func(c *fiber.Ctx) error {
		log.WithFields(log.Fields{
			"method": c.Method(),
			"path":   c.Path(),
		}).Warn("No route matched")

		return c.Status(fiber.StatusNotFound).JSON(map[string]string{
			"error":   "Not Found",
			"message": "The requested resource does not exist",
		})
	})
}
