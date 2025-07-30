package main

import (
	"github.com/arthadede/auth-service/constants"
	"github.com/arthadede/auth-service/controllers"
	"github.com/arthadede/auth-service/database"
	_ "github.com/arthadede/auth-service/docs"
	"github.com/arthadede/auth-service/handlers"
	"github.com/arthadede/auth-service/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"log"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	cfg := constants.NewConfig()

	mongoDB := database.NewMongoDB(cfg)
	defer mongoDB.Close()

	userRepo := repositories.NewUserRepository(mongoDB)

	userController := controllers.NewUserController(userRepo, cfg)

	userHandler := handlers.NewUserHandler(userController)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(helmet.New())

	//if cfg.Environment == "production" {
	//	app.Use(limiter.New())
	//}

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} | ${latency} | ${ip} | ${method} ${path}${query} | ${ua} | ${locals:UserID}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth",
		})
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)
	auth.Post("/forgot-password", userHandler.ForgotPassword)
	auth.Post("/reset-password", userHandler.ResetPassword)

	// Google OAuth routes
	auth.Get("/google", userHandler.GoogleLogin)
	auth.Post("/google/callback", userHandler.GoogleCallback)

	log.Fatal(app.Listen(cfg.Port))
}
