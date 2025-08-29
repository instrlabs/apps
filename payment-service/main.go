package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/arthadede/payment-service/internal"
)

func main() {
	// Initialize configuration
	cfg := internal.NewConfig()

	// Initialize MongoDB
	mongodb := internal.NewMongoDB(cfg)
	defer mongodb.Close()

	// Initialize payment repository
	paymentRepo := internal.NewPaymentRepository(mongodb)

	// Initialize Midtrans service
	midtransService := internal.NewMidtransService(cfg)

	// Initialize NATS service
	natsService, err := internal.NewNatsService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize NATS service: %v", err)
	}
	defer natsService.Close()

	// Initialize payment handler
	paymentHandler := internal.NewPaymentHandler(midtransService, paymentRepo, natsService, cfg)

	// Create Fiber app
	app := fiber.New()

	// Configure middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} | ${latency} | ${ip} | ${method} ${path}${query} | ${ua}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowMethods: "GET, POST, OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// Register routes
	app.Post("/payments", paymentHandler.CreatePayment)
	app.Get("/payments/:orderId", paymentHandler.GetPaymentStatus)
	app.Post("/payments/notification", paymentHandler.HandleNotification)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "payment",
		})
	})

	// Subscribe to payment requests from NATS
	err = natsService.SubscribeToPaymentRequests(paymentHandler.ProcessPaymentRequest)
	if err != nil {
		log.Fatalf("Failed to subscribe to payment requests: %v", err)
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting payment service on port %s", cfg.Port)
		if err := app.Listen(cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shutdown the server
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
