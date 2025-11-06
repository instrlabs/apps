package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/gateway-service/internal"
	"github.com/instrlabs/shared/middlewarex"
)

func main() {
	cfg := internal.LoadConfig()

	app := fiber.New(fiber.Config{})

	middlewarex.SetupPrometheus(app)
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupLogger(app)
	internal.SetupGatewaySwaggerUI(app)
	internal.SetupMiddleware(app, cfg)
	internal.SetupGatewayRoutes(app, cfg)

	log.Fatal(app.Listen(cfg.Port))
}
