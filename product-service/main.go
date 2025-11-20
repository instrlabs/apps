package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/shared/initx"
	"github.com/instrlabs/shared/middlewarex"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/instrlabs/product-service/internal/config"
	"github.com/instrlabs/product-service/internal/handlers"
	"github.com/instrlabs/product-service/internal/repositories"
	"github.com/instrlabs/product-service/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	client, db := initx.NewMongo()
	defer initx.CloseMongo(client)

	// Create Fiber app
	app := setupFiberApp(cfg)

	// Setup routes
	setupRoutes(app, db, cfg)

	// Start server
	startServer(app, cfg)
}

// setupFiberApp configures the Fiber application with middleware
func setupFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:   time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:  time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:   time.Duration(cfg.IdleTimeout) * time.Second,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "product-service",
		AppName:       "product-service",
	})

	// Setup shared middleware
	middlewarex.SetupPrometheus(app)
	middlewarex.SetupLogger(app)
	middlewarex.SetupServiceSwagger(app, "/products")
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupAuthentication(app)

	return app
}

// setupRoutes configures all application routes
func setupRoutes(app *fiber.App, db interface{}, cfg *config.Config) {
	// Initialize repositories
	productRepo := repositories.NewProductRepository(db.(*mongo.Database))

	// Initialize services
	productService := services.NewProductService(productRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)

	// Health check route (must come before wildcard routes)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":       "ok",
			"timestamp":    time.Now().UTC(),
			"service":      "product-service",
			"architecture": "simplified",
		})
	})

	// Type-based routes
	app.Get("/:type/:id", productHandler.GetProductByID)
	app.Get("/:type", productHandler.ListProductsByType)
}

// startServer starts the server with graceful shutdown
func startServer(app *fiber.App, cfg *config.Config) {
	// Start server in goroutine
	go func() {
		port := cfg.Port
		if port == "" {
			port = "3005"
		}
		fiberlog.Infof("Starting product service on port %s", port)
		if err := app.Listen(port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fiberlog.Info("Shutting down server...")

	// Shutdown Fiber app
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fiberlog.Info("Server exited")
}
