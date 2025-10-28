package internal

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/proxy"
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

	for _, srv := range config.Services {
		prefix := srv.Prefix
		targetURL, _ := url.Parse(srv.URL)

		app.All(prefix+"/*", func(c *fiber.Ctx) error {
			forwardPath := c.Path()[len(prefix):]
			queryString := string(c.Request().URI().QueryString())

			log.Infof("Forwarding request to %s%s", srv.URL, forwardPath)

			parsedUrl := targetURL.String() + forwardPath
			if queryString != "" {
				parsedUrl += "?" + queryString
			}

			if err := proxy.DoTimeout(c, parsedUrl, 30*time.Second); err != nil {
				log.Errorf("proxy error: service=%s method=%s path=%s query=%s err=%v", srv.Name, c.Method(), forwardPath, queryString, err)

				return c.Status(fiber.StatusBadGateway).JSON(map[string]string{
					"error":   "Bad Gateway",
					"message": "The srv is currently unavailable",
				})
			}

			return nil
		})

		log.Infof("Service registered: service=%s target=%s", srv.Name, srv.URL)
	}

	app.Use(func(c *fiber.Ctx) error {
		log.Warnf("No route matched: method=%s path=%s", c.Method(), c.Path())

		return c.Status(fiber.StatusNotFound).JSON(map[string]string{
			"error":   "Not Found",
			"message": "The requested resource does not exist",
		})
	})
}
