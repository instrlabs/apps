package main

import (
	"github.com/gofiber/fiber/v2"
	initx "github.com/instrlabs/shared/init"
	"github.com/instrlabs/gateway-service/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := internal.LoadConfig()

	app := fiber.New(fiber.Config{})

	initx.SetupLogger(app)
	internal.SetupMiddleware(app, cfg)
	internal.SetupGatewaySwaggerUI(app, cfg)
	internal.SetupGatewayRoutes(app, cfg)

	log.Fatal(app.Listen(cfg.Port))
}
