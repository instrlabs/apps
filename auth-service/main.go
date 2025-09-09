package main

import (
	"log"

	"github.com/arthadede/auth-service/internal"
	initx "github.com/histweety-labs/shared/init"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := internal.LoadConfig()

	mongo := initx.NewMongo(&initx.MongoConfig{MongoURI: config.MongoURI, MongoDB: config.MongoDB})
	defer mongo.Close()

	userRepo := internal.NewUserRepository(mongo)
	userController := internal.NewUserController(userRepo, config)
	userHandler := internal.NewUserHandler(userController, config)

	app := fiber.New(fiber.Config{})

	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)

	internal.SetupMiddleware(app)

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Post("/logout", userHandler.Logout)
	app.Post("/refresh", userHandler.RefreshToken)
	app.Post("/forgot-password", userHandler.ForgotPassword)
	app.Post("/reset-password", userHandler.ResetPassword)

	app.Get("/profile", userHandler.GetProfile)
	app.Put("/profile", userHandler.UpdateProfile)
	app.Post("/change-password", userHandler.ChangePassword)

	app.Get("/google", userHandler.GoogleLogin)
	app.Get("/google/callback", userHandler.GoogleCallback)

	log.Fatal(app.Listen(config.Port))
}
