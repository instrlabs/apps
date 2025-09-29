package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/instr-labs/gateway-service/internal"
	initx "github.com/instr-labs/shared/init"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := internal.LoadConfig()

	app := fiber.New(fiber.Config{})

	initx.SetupLogger(app)
	internal.SetupMiddleware(app, config)
	internal.SetupGatewaySwaggerUI(app, config)
	internal.SetupGatewayRoutes(app, config)

	log.Fatal(app.Listen(config.Port))
}
