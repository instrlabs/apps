package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/shared/initx"
	"github.com/instrlabs/shared/middlewarex"
	natsgo "github.com/nats-io/nats.go"

	"github.com/instrlabs/image-service/internal"
)

func main() {
	cfg := internal.LoadConfig()
	s3 := initx.NewS3()

	mongoClient, mongoDB := initx.NewMongo()
	defer initx.CloseMongo(mongoClient)

	nats := initx.NewNats(cfg.NatsURI)
	defer initx.CloseNats(nats)

	app := fiber.New(fiber.Config{})

	middlewarex.SetupPrometheus(app)
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupServiceSwagger(app, "/images")
	middlewarex.SetupAuthentication(app)

	instrRepo := internal.NewInstructionRepository(mongoDB)
	detailRepo := internal.NewInstructionDetailRepository(mongoDB)

	imageSvc := internal.NewImageService()
	productClient := internal.NewProductClient(cfg.ProductServiceURL)

	instrHandler := internal.NewInstructionHandler(cfg, s3, nats, instrRepo, detailRepo, productClient, imageSvc)
	instrProcessor := internal.NewInstructionProcessor(cfg, s3, nats, instrRepo, detailRepo, productClient, imageSvc)

	_, _ = nats.Subscribe(cfg.NatsSubjectImageRequests, func(m *natsgo.Msg) {
		instrProcessor.RunInstructionMessage(m.Data)
	})

	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			instrProcessor.CleanInstruction()
		}
	}()

	app.Post("/instructions", instrHandler.CreateInstruction)
	app.Post("/instructions/:id", instrHandler.CreateInstructionDetails)
	app.Get("/instructions/:id", instrHandler.GetInstructionByID)
	app.Get("/instructions/:id/:detail_id", instrHandler.GetInstructionDetail)
	app.Get("/instructions/:id/:detail_id/file", instrHandler.GetInstructionDetailFile)

	log.Fatal(app.Listen(cfg.Port))
}
