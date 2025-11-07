package internal

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	natsgo "github.com/nats-io/nats.go"
)

func SetupNotificationRoutes(app *fiber.App, cfg *Config, natsSrv *natsgo.Conn) {
	app.Get("/health", func(c *fiber.Ctx) error {
		health := map[string]interface{}{
			"status":   "ok",
			"services": map[string]string{},
		}

		servicesStatus := health["services"].(map[string]string)

		// Check NATS connection
		if natsSrv != nil && natsSrv.IsConnected() {
			servicesStatus["nats"] = "ok"
		} else {
			servicesStatus["nats"] = "error"
		}

		// Check Auth Service
		resp, err := http.Get(cfg.AuthService + "/health")
		if err != nil || resp.StatusCode != http.StatusOK {
			servicesStatus["auth-service"] = "error"
		} else {
			servicesStatus["auth-service"] = "ok"
		}

		// Overall service status
		allOk := true
		for _, status := range servicesStatus {
			if status != "ok" {
				allOk = false
				break
			}
		}

		if !allOk {
			health["status"] = "degraded"
		}

		return c.JSON(health)
	})
}
