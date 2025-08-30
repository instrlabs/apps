package main

import (
	"time"

	"github.com/arthadede/gateway-service/internal"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := internal.LoadConfig()

	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	internal.SetupMiddleware(app)

	internal.SetupGatewayRoutes(app, config)

	log.Fatal(app.Listen(config.Port))
}
