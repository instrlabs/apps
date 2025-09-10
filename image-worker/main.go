package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthadede/image-worker/internal"
	"github.com/gofiber/fiber/v2"
	initx "github.com/histweety-labs/shared/init"
	"github.com/nats-io/nats.go"
)

func main() {
	cfg := internal.NewConfig()

	s3Srv := initx.NewS3(&initx.S3Config{
		S3Endpoint:  cfg.S3Endpoint,
		S3AccessKey: cfg.S3AccessKey,
		S3SecretKey: cfg.S3SecretKey,
		S3UseSSL:    cfg.S3UseSSL,
		S3Region:    cfg.S3Region,
		S3Bucket:    cfg.S3Bucket,
	})
	mongoSrv := initx.NewMongo(&initx.MongoConfig{MongoURI: cfg.MongoURI, MongoDB: cfg.MongoDB})
	defer mongoSrv.Close()
	natsSrv := initx.NewNats(cfg.NatsURL)
	defer natsSrv.Close()

	imageServ := internal.NewImageService()
	instrRepo := internal.NewInstructionRepository(mongoSrv)
	processor := internal.NewProcessor(mongoSrv, s3Srv, natsSrv, imageServ, instrRepo)

	app := fiber.New(fiber.Config{})
	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{})

	if natsSrv != nil {
		_, err := natsSrv.Conn.Subscribe(cfg.NatsSubjectRequests, func(msg *nats.Msg) {
			var job internal.JobMessage
			if err := json.Unmarshal(msg.Data, &job); err != nil {
				log.Printf("failed to unmarshal job: %v", err)
				return
			}
			if err := processor.Handle(context.Background(), &job); err != nil {
				log.Printf("job handler error: %v", err)
			}
		})
		if err != nil {
			log.Fatalf("failed to subscribe: %v", err)
		}
	}

	sched := internal.NewScheduler(instrRepo, natsSrv, cfg)
	sched.Start()
	defer sched.Stop()

	go func() {
		log.Printf("image-worker listening on %s", cfg.Port)
		if err := app.Listen(cfg.Port); err != nil {
			log.Printf("fiber server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(ctx)
}
