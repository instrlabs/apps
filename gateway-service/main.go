package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/gateway-service/internal"
	initx "github.com/instrlabs/shared/init"
)

func main() {
	cfg := internal.LoadConfig()

	app := fiber.New(fiber.Config{})

	initx.SetupPrometheus(app)
	initx.SetupServiceHealth(app)
	initx.SetupLogger(app)
	internal.SetupMiddleware(app, cfg)
	internal.SetupGatewaySwaggerUI(app)
	internal.SetupGatewayRoutes(app, cfg)

	log.Fatal(app.Listen(cfg.Port))
}
