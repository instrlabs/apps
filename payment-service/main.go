package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	initx "github.com/instrlabs/shared/init"
	"github.com/instrlabs/payment-service/internal"
)

func main() {
	config := internal.LoadConfig()

	mongo := initx.NewMongo(&initx.MongoConfig{
		MongoURI: config.MongoURI,
		MongoDB:  config.MongoDB,
	})
	defer mongo.Close()

	productRepo := internal.NewProductRepository(mongo)
	productHandler := internal.NewProductHandler(productRepo)

	app := fiber.New(fiber.Config{})

	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{})

	app.Post("/products", productHandler.CreateProduct)
	app.Get("/products", productHandler.ListProducts)
	app.Get("/products/:id", productHandler.GetProduct)
	app.Patch("/products/:id", productHandler.UpdateProduct)
	app.Delete("/products/:id", productHandler.DeleteProduct)

	log.Fatal(app.Listen(config.Port))
}
