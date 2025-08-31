package main

import (
	"log"

	"github.com/arthadede/payment-service/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := internal.LoadConfig()

	mongo := internal.NewMongoDB(config)
	defer mongo.Close()
	natsService := internal.NewNatsService(config)
	defer natsService.Close()

	productRepo := internal.NewProductRepository(mongo)
	productHandler := internal.NewProductHandler(productRepo)

	app := fiber.New(fiber.Config{})

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Type("json").SendFile("./static/swagger.json")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	internal.SetupMiddleware(app)

	app.Post("/products", productHandler.CreateProduct)
	app.Get("/products", productHandler.ListProducts)
	app.Get("/products/:id", productHandler.GetProduct)
	app.Patch("/products/:id", productHandler.UpdateProduct)
	app.Delete("/products/:id", productHandler.DeleteProduct)

	log.Fatal(app.Listen(config.Port))
}
