package internal

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	log "github.com/sirupsen/logrus"
)

func SetupGatewayRoutes(app *fiber.App, config *Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		health := map[string]interface{}{
			"status":   "ok",
			"services": map[string]string{},
		}

		servicesStatus := health["services"].(map[string]string)

		for _, service := range config.Services {
			status := "ok"
			resp, err := http.Get(service.URL + "/health")
			if err != nil || resp.StatusCode != http.StatusOK {
				status = "error"
			}

			servicesStatus[service.Name] = status
		}

		return c.JSON(health)
	})

	for _, service := range config.Services {
		prefix := service.Prefix
		targetURL, _ := url.Parse(service.URL)

		app.All(prefix+"/*", func(c *fiber.Ctx) error {
			forwardPath := c.Path()[len(prefix):]
			queryString := string(c.Request().URI().QueryString())

			log.WithFields(log.Fields{
				"service":      service.Name,
				"method":       c.Method(),
				"path":         forwardPath,
				"query":        queryString,
				"forwarded_to": targetURL.String(),
			}).Info("Forwarding request")

			if token, ok := c.Locals("token").(string); ok && token != "" {
				c.Request().Header.Set("X-Auth-Token", token)
				log.WithFields(log.Fields{
					"service": service.Name,
					"path":    forwardPath,
				}).Info("Forwarding token in X-Auth-Token header")
			}

			url := targetURL.String() + forwardPath
			if queryString != "" {
				url += "?" + queryString
			}

			if err := proxy.DoTimeout(c, url, 30*time.Second); err != nil {
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
