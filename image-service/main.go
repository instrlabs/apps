package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/arthadede/image-service/internal"
)

func main() {
	config := internal.LoadConfig()

	s3Service := internal.NewS3Service(config)
	mongo := internal.NewMongoDB(config)
	defer mongo.Close()

	fileRepo := internal.NewFileRepository(mongo)
	instrRepo := internal.NewInstructionRepository(mongo)
	productSvc := internal.NewProductService()

	workers := internal.NewWorkerPool(2, s3Service, instrRepo, fileRepo, productSvc)
	defer workers.Stop()

	instructionHandler := internal.NewInstructionHandler(s3Service, fileRepo, instrRepo, productSvc, workers)

	app := fiber.New(fiber.Config{})

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Type("json").SendFile("./static/swagger.json")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	internal.SetupMiddleware(app)

	app.Post("/compress", instructionHandler.ImageCompress)

	log.Fatal(app.Listen(config.Port))
}
