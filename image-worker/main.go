package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthadede/image-worker/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := internal.NewConfig()

	mongo := internal.NewMongoDB(cfg)
	defer mongo.Close()
	s3 := internal.NewS3Service(cfg)
	natsSvc := internal.NewNatsService(cfg)
	defer natsSvc.Close()

	processor := internal.NewProcessor(mongo, s3, natsSvc)

	if err := natsSvc.Subscribe(processor.Handle); err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	app := fiber.New(fiber.Config{})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

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
