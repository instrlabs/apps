package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/auth-service-2/internal"
	initx "github.com/instrlabs/shared/init"
)

func main() {
	// Load configuration
	cfg := internal.LoadConfig()

	// Initialize MongoDB
	mongo := initx.NewMongo(&initx.MongoConfig{
		MongoURI: cfg.MongoURI,
		MongoDB:  cfg.MongoDB,
	})
	defer mongo.Close()

	// Initialize repositories
	userRepo := internal.NewUserRepository(mongo)

	// Initialize services
	tokenService := internal.NewTokenService(cfg)
	emailService := internal.NewEmailService(cfg)

	// Initialize handlers
	authHandler := internal.NewAuthHandler(cfg, userRepo, tokenService)
	oauthHandler := internal.NewOAuthHandler(cfg, userRepo, tokenService)
	pinHandler := internal.NewPinHandler(cfg, userRepo, tokenService, emailService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Auth Service 2",
	})

	// Setup middleware and utilities
	initx.SetupPrometheus(app)
	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app, cfg.APIBaseURL, "/auth")
	initx.SetupServiceHealth(app)

	// Setup authentication middleware (exclude public routes)
	initx.SetupAuthenticated(app, []string{
		"/oauth/google",
		"/oauth/google/callback",
		"/auth/pin/request",
		"/auth/pin/verify",
		"/auth/refresh",
	})

	// ==========================================
	// ROUTES - High-Level Flow
	// ==========================================

	// OAuth Routes (Google)
	app.Get("/oauth/google", oauthHandler.GoogleLogin)
	app.Get("/oauth/google/callback", oauthHandler.GoogleCallback)

	// PIN Authentication Routes
	app.Post("/auth/pin/request", pinHandler.RequestPin)
	app.Post("/auth/pin/verify", pinHandler.VerifyPin)

	// Token Management Routes
	app.Post("/auth/refresh", authHandler.RefreshToken)
	app.Post("/auth/logout", authHandler.Logout)

	// User Routes
	app.Get("/auth/profile", authHandler.GetProfile)

	// Start server
	log.Infof("🚀 Auth Service 2 starting on %s", cfg.Port)
	log.Infof("📝 Environment: %s", cfg.Environment)
	log.Infof("🔐 PIN Authentication: %v", cfg.PinEnabled)
	log.Fatal(app.Listen(cfg.Port))
}
