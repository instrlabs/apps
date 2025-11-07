package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/shared/initx"
	"github.com/instrlabs/shared/middlewarex"

	"github.com/instrlabs/product-service/internal"
)

func main() {
	cfg := internal.LoadConfig()

	mongoClient, mongoDB := initx.NewMongo()
	defer initx.CloseMongo(mongoClient)

	app := fiber.New(fiber.Config{})

	middlewarex.SetupPrometheus(app)
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupServiceSwagger(app, "/products")
	middlewarex.SetupAuthentication(app, []string{})

	productRepo := internal.NewProductRepository(mongoDB)
	productHandler := internal.NewProductHandler(productRepo)

	app.Get("/", productHandler.ListProducts)

	log.Fatal(app.Listen(cfg.Port))
}
