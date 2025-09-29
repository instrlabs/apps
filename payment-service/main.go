package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/instrlabs/payment-service/internal"
	initx "github.com/instrlabs/shared/init"
)

func main() {
	config := internal.LoadConfig()

	mongo := initx.NewMongo(&initx.MongoConfig{
		MongoURI: config.MongoURI,
		MongoDB:  config.MongoDB,
	})
	defer mongo.Close()

	// Product endpoints have been moved to image-service. No local product repository/handler.

	app := fiber.New(fiber.Config{})

	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{})

	// Product routes moved to image-service. This service no longer exposes /products endpoints.

	log.Fatal(app.Listen(config.Port))
}
