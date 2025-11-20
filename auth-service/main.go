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

	"github.com/instrlabs/auth-service/internal/config"
	"github.com/instrlabs/auth-service/internal/handlers"
	"github.com/instrlabs/auth-service/internal/helpers"
	"github.com/instrlabs/auth-service/internal/repositories"
	"github.com/instrlabs/auth-service/internal/services"
	"github.com/instrlabs/auth-service/internal/validators"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

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
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:   60 * time.Second,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "auth-service",
		AppName:       "auth-service",
	})

	// Setup shared middleware
	middlewarex.SetupPrometheus(app)
	middlewarex.SetupLogger(app)
	middlewarex.SetupServiceSwagger(app, "/auth")
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupAuthentication(app)

	return app
}

// setupRoutes configures all application routes
func setupRoutes(app *fiber.App, db interface{}, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.(*mongo.Database))

	// Initialize helpers
	emailService := helpers.NewEmailService()

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.TokenExpiryHours)
	pinService := services.NewPinService(userRepo, emailService, cfg.PinEnabled)
	oauthService := services.NewOAuthService(
		userRepo,
		authService,
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectUrl,
		cfg.WebUrl,
	)
	userService := services.NewUserService(userRepo)

	// Initialize validators
	validator := validators.NewRequestValidator()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, validator)
	pinHandler := handlers.NewPinHandler(pinService, validator)
	oauthHandler := handlers.NewOAuthHandler(oauthService, authService)
	userHandler := handlers.NewUserHandler(userService)

	// Setup routes
	app.Post("/login", authHandler.Login)
	app.Post("/logout", authHandler.Logout)
	app.Post("/refresh", authHandler.RefreshToken)
	app.Post("/send-pin", pinHandler.SendPin)
	app.Get("/profile", userHandler.GetProfile)
	app.Get("/google", oauthHandler.GoogleLogin)
	app.Get("/google/callback", oauthHandler.GoogleCallback)

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":       "ok",
			"timestamp":    time.Now().UTC(),
			"service":      "auth-service",
			"architecture": "simplified",
		})
	})
}

// startServer starts the server with graceful shutdown
func startServer(app *fiber.App, cfg *config.Config) {
	// Start server in goroutine
	go func() {
		port := cfg.Port
		if port == "" {
			port = "3001"
		}
		fiberlog.Infof("Starting auth service on port %s", port)
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
