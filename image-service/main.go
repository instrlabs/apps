package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	initx "github.com/histweety-labs/shared/init"

	"github.com/arthadede/image-service/internal"
)

func main() {
	cfg := internal.LoadConfig()

	s3 := initx.NewS3(&initx.S3Config{
		S3Endpoint:  cfg.S3Endpoint,
		S3AccessKey: cfg.S3AccessKey,
		S3SecretKey: cfg.S3SecretKey,
		S3UseSSL:    cfg.S3UseSSL,
		S3Region:    cfg.S3Region,
		S3Bucket:    cfg.S3Bucket,
	})
	mongo := initx.NewMongo(&initx.MongoConfig{
		MongoURI: cfg.MongoURI,
		MongoDB:  cfg.MongoDB,
	})
	defer mongo.Close()
	nats := initx.NewNats(cfg.NatsURL)
	defer nats.Close()

	app := fiber.New(fiber.Config{})
	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{})

	instrRepo := internal.NewInstructionRepository(mongo)

	paymentSvc := internal.NewPaymentService(cfg)

	instrHandler := internal.NewInstructionHandler(cfg, s3, nats, instrRepo, paymentSvc)

	app.Get("/instructions", instrHandler.ListInstructions)
	app.Get("/instructions/:id", instrHandler.GetInstructionByID)
	app.Get("/instructions/:id/:file_name", instrHandler.GetInstructionFile)
	app.Patch("/instructions/:id/status", instrHandler.UpdateInstructionStatus)
	app.Patch("/instructions/:id/outputs", instrHandler.UpdateInstructionOutputs)

	app.Post("/compress", instrHandler.ImageCompress)

	log.Fatal(app.Listen(cfg.Port))
}
