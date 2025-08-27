package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"

	"github.com/arthadede/pdf-service/internal"
)

func main() {
	cfg := internal.NewConfig()

	mongo, err := internal.NewMongoService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB service: %v", err)
	}
	defer mongo.Close()

	s3Service, err := internal.NewS3Service(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}

	natsService, err := internal.NewNatsService(cfg)
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

	jobRepo := internal.NewJobRepository(mongo.Database)
	pdfJobRepo := internal.NewPDFJobRepository(mongo.Database)
	pdfJobHandler := internal.NewPDFJobHandler(pdfJobRepo, s3Service, natsService)
	pdfJobProcessor := internal.NewPDFJobProcessor(jobRepo, pdfJobRepo, s3Service, natsService)

	err = natsService.SubscribeToPDFJob(pdfJobProcessor.ProcessJob)
	if err != nil {
		log.Fatalf("Failed to subscribe to PDF jobs: %v", err)
	}

	app.Get("/jobs/:id", pdfJobHandler.GetFile)
	app.Get("/jobs", pdfJobHandler.GetJobs)
	app.Post("/jobs", pdfJobHandler.CreateJob)

	log.Fatal(app.Listen(cfg.Port))
}
