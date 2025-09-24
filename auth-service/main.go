package main

import (
	"log"

	"github.com/arthadede/auth-service/internal"
	initx "github.com/histweety-labs/shared/init"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := internal.LoadConfig()

	mongo := initx.NewMongo(&initx.MongoConfig{
		MongoURI: config.MongoURI,
		MongoDB:  config.MongoDB,
	})
	defer mongo.Close()

	userRepo := internal.NewUserRepository(mongo)
	userHandler := internal.NewUserHandler(config, userRepo)

	app := fiber.New(fiber.Config{})

	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{
		"/check",
		"/login",
		"/refresh",
		"/google",
		"/send-pin",
	})

	app.Post("/login", userHandler.Login)
	app.Post("/logout", userHandler.Logout)
	app.Post("/refresh", userHandler.RefreshToken)
	app.Post("/send-pin", userHandler.SendPin)
	app.Post("/check", userHandler.CheckEmail)

	app.Get("/profile", userHandler.GetProfile)

	app.Get("/google", userHandler.GoogleLogin)
	app.Get("/google/callback", userHandler.GoogleCallback)

	log.Fatal(app.Listen(config.Port))
}
