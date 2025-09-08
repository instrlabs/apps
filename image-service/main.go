package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/arthadede/image-service/internal"
)

func main() {
	cfg := internal.LoadConfig()

	s3Service := internal.NewS3Service(cfg)
	mongo := internal.NewMongoDB(cfg)
	defer mongo.Close()
	natsSvc := internal.NewNatsService(cfg)
	defer natsSvc.Close()

	fileRepo := internal.NewFileRepository(mongo)
	instrRepo := internal.NewInstructionRepository(mongo)
	productSvc := internal.NewProductService()

	instructionHandler := internal.NewInstructionHandler(s3Service, fileRepo, instrRepo, productSvc, natsSvc)

	app := fiber.New(fiber.Config{})

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Type("json").SendFile("./static/swagger.json")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	internal.SetupMiddleware(app)

	app.Post("/compress", instructionHandler.ImageCompress)

	log.Fatal(app.Listen(cfg.Port))
}
