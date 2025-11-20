package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/shared/middlewarex"
	natsgo "github.com/nats-io/nats.go"
)

func main() {
	app := fiber.New(fiber.Config{})

	middlewarex.SetupPrometheus(app)
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupServiceSwagger(app, "/payment")
	middlewarex.SetupAuthentication(app)

	app.Post("/transactions", instrHandler.CreateTransaction)
	app.Get("/transactions", instrHandler.GetInstructionByID)

	//log.Fatal(app.Listen(cfg.Port))
}
