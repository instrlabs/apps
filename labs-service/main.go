package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/arthadede/labs-service/internal"
)

func main() {
	config := internal.LoadConfig()

	mongo := internal.NewMongoDB(config)
	defer mongo.Close()
	s3Service := internal.NewS3Service(config)

	app := fiber.New(fiber.Config{})

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Type("json").SendFile("./static/swagger.json")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	internal.SetupMiddleware(app)

	app.Post("/image/compress", imageHandler.ImageCompress)

	log.Fatal(app.Listen(config.Port))
}
