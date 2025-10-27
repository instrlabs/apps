package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	initx "github.com/instrlabs/shared/init"
	natsgo "github.com/nats-io/nats.go"

	"github.com/instrlabs/image-service/internal"
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
	nats := initx.NewNats(cfg.NatsURI)
	defer nats.Close()

	app := fiber.New(fiber.Config{})

	initx.SetupPrometheus(app)
	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{
		"/products",
	})

	productRepo := internal.NewProductRepository(mongo)
	instrRepo := internal.NewInstructionRepository(mongo)
	detailRepo := internal.NewInstructionDetailRepository(mongo)

	imageSvc := internal.NewImageService()

	productHandler := internal.NewProductHandler(productRepo)
	instrHandler := internal.NewInstructionHandler(cfg, s3, nats, instrRepo, detailRepo, productRepo, imageSvc)

	_, _ = nats.Conn.Subscribe(cfg.NatsSubjectImageRequests, func(m *natsgo.Msg) {
		instrHandler.RunInstructionMessage(m.Data)
	})

	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			err := instrHandler.CleanInstruction()
			if err != nil {
				log.Errorf("CleanInstruction: %v", err)
			}
		}
	}()

	app.Post("/instructions", instrHandler.CreateInstruction)
	app.Post("/instructions/:id/details", instrHandler.CreateInstructionDetails)

	app.Get("/instructions/:id/details/:detailId", instrHandler.GetInstructionDetail)
	app.Get("/instructions/:id/details/:detailId/file", instrHandler.GetInstructionDetilFile)
	app.Get("/instructions", instrHandler.ListInstructions)
	app.Get("/instructions/:id", instrHandler.GetInstructionByID)
	app.Get("/instructions/:id/details", instrHandler.GetInstructionDetails)

	app.Get("/files", instrHandler.ListUncleanedFiles)

	app.Get("/products", productHandler.ListProducts)

	log.Fatal(app.Listen(cfg.Port))
}
