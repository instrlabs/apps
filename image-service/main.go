package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	initx "github.com/histweety-labs/shared/init"
	natsgo "github.com/nats-io/nats.go"

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
	fileRepo := internal.NewFileRepository(mongo)

	paymentSvc := internal.NewPaymentService(cfg)

	imageSvc := internal.NewImageService()
	instrHandler := internal.NewInstructionHandler(cfg, s3, nats, instrRepo, fileRepo, paymentSvc, imageSvc)

	_, _ = nats.Conn.Subscribe(cfg.NatsSubjectImagesRequests, func(m *natsgo.Msg) {
		instrHandler.RunInstructionMessage(m.Data)
	})

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			_ = instrHandler.CleanInstruction(nil)
		}
	}()

	app.Post("/instructions/:productKey", instrHandler.CreateInstruction)
	app.Post("/instructions/:id/details", instrHandler.CreateInstructionDetails)

	app.Get("/instructions", instrHandler.ListInstructions)
	app.Get("/instructions/:id", instrHandler.GetInstructionByID)
	app.Get("/instructions/:id/details", instrHandler.GetInstructionDetails)
	app.Get("/instructions/:id/details/:fileId", instrHandler.GetInstructionFileBytes)

	log.Fatal(app.Listen(cfg.Port))
}
