package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/instr-labs/gateway-service/internal"
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
	internal.SetupGatewaySwaggerUI(app, config)
	internal.SetupGatewayRoutes(app, config)

	log.Fatal(app.Listen(config.Port))
}
