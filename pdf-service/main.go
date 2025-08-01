package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"

	"pdf-service/constants"
	"pdf-service/services"
)

func main() {
	cfg := constants.NewConfig()

	mongo, err := services.NewMongoService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB service: %v", err)
	}
	defer mongo.Close()

	s3Service, err := services.NewS3Service(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}

	natsService, err := services.NewNatsService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize NATS service: %v", err)
	}
	defer natsService.Close()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(helmet.New())

	if cfg.Environment == "production" {
		app.Use(limiter.New())
	}

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} | ${latency} | ${ip} | ${method} ${path}${query} | ${ua}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "labs",
		})
	})

	_ = s3Service

	log.Fatal(app.Listen(cfg.Port))
}
