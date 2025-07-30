package main

import (
	"github.com/arthadede/auth-service/constants"
	"github.com/arthadede/auth-service/controllers"
	"github.com/arthadede/auth-service/database"
	"github.com/arthadede/auth-service/handlers"
	"github.com/arthadede/auth-service/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

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

	if cfg.Environment == "production" {
		app.Use(limiter.New())
	}

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} | ${latency} | ${ip} | ${method} ${path}${query} | ${ua} | ${locals:UserID}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth",
		})
	})

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Post("/refresh", userHandler.RefreshToken)
	app.Post("/forgot-password", userHandler.ForgotPassword)
	app.Post("/reset-password", userHandler.ResetPassword)

	// Google OAuth routes
	app.Get("/google", userHandler.GoogleLogin)
	app.Post("/google/callback", userHandler.GoogleCallback)

	log.Fatal(app.Listen(cfg.Port))
}
