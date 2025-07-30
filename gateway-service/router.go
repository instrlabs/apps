package main

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	log "github.com/sirupsen/logrus"
)

func SetupGatewayRoutes(app *fiber.App, config Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		health := map[string]interface{}{
			"status":   "ok",
			"services": map[string]string{},
		}

		servicesStatus := health["services"].(map[string]string)

		for _, service := range config.Services {
			if service.URL == "" {
				servicesStatus[service.Name] = "not configured"
				continue
			}

			agent := fiber.AcquireAgent()
			req := agent.Request()
			req.SetRequestURI(service.URL + "/health")
			req.Header.SetMethod(fiber.MethodGet)

			if err := agent.Parse(); err != nil {
				servicesStatus[service.Name] = "error"
				fiber.ReleaseAgent(agent)
				continue
			}

			code, _, errs := agent.Bytes()
			fiber.ReleaseAgent(agent)

			if len(errs) > 0 || code != fiber.StatusOK {
				servicesStatus[service.Name] = "error"
			} else {
				servicesStatus[service.Name] = "ok"
			}
		}

		for _, status := range servicesStatus {
			if status == "error" {
				health["status"] = "partial"
				break
			}
		}

		return c.JSON(health)
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

		prefix := service.Prefix
		app.All(prefix+"/*", func(c *fiber.Ctx) error {
			forwardPath := c.Path()[len(prefix):]
			log.WithFields(log.Fields{
				"service":      service.Name,
				"method":       c.Method(),
				"path":         forwardPath,
				"forwarded_to": targetURL.String(),
			}).Info("Forwarding request")

			c.Request().Header.Set("X-Gateway", "true")

			url := targetURL.String() + forwardPath
			if err := proxy.Do(c, url); err != nil {
				log.WithFields(log.Fields{
					"service": service.Name,
					"method":  c.Method(),
					"path":    forwardPath,
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
