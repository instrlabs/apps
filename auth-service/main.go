// @title Auth Service API
// @version 1.0
// @description Authentication service API documentation
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /auth
package main

import (
	"log"

	"github.com/arthadede/auth-service/constants"
	"github.com/arthadede/auth-service/controllers"
	"github.com/arthadede/auth-service/database"
	"github.com/arthadede/auth-service/docs"
	"github.com/arthadede/auth-service/handlers"
	"github.com/arthadede/auth-service/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func main() {
	// Swagger documentation setup
	docs.SwaggerInfo.Title = "Auth Service API"
	docs.SwaggerInfo.Description = "Authentication service API documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/auth"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	cfg := constants.NewConfig()

	mongo := database.NewMongoDB(cfg)
	defer mongo.Close()

	userRepo := repositories.NewUserRepository(mongo)

	userController := controllers.NewUserController(userRepo, cfg)

	userHandler := handlers.NewUserHandler(userController)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(helmet.New())

	if cfg.Environment == "production" {
		app.Use(limiter.New())
	}

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} | ${latency} | ${ip} | ${method} ${path}${query} | ${ua} | ${locals:UserID}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	// Swagger documentation route
	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title:        "Auth Service API",
		DeepLinking:  false,
		DocExpansion: "list",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth",
		})
	})

	app.Use(AuthMiddleware())

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Post("/refresh", userHandler.RefreshToken)
	app.Post("/forgot-password", userHandler.ForgotPassword)
	app.Post("/reset-password", userHandler.ResetPassword)
	app.Post("/verify-token", userHandler.VerifyToken)

	// Google OAuth routes
	app.Get("/google", userHandler.GoogleLogin)
	app.Get("/google/callback", userHandler.GoogleCallback)

	log.Fatal(app.Listen(cfg.Port))
}
