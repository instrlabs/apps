package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arthadede/notification-service/internal"
	"github.com/gofiber/fiber/v2"
	initx "github.com/histweety-labs/shared/init"
)

func main() {
	cfg := internal.NewConfig()

	nats := initx.NewNats(cfg.NatsURL)
	defer nats.Close()

	sseService := internal.NewSSEService(cfg)

	app := fiber.New(fiber.Config{})
	initx.SetupLogger(app)
	initx.SetupServiceHealth(app)
	internal.SetupMiddleware(app)

	app.Get("/sse", sseService.HandleSSE)

	go func() {
		log.Printf("Starting notification service on %s", cfg.Port)
		if err := app.Listen(cfg.Port); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server exiting")
}
