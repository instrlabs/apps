package main

import (
	"log"

	"github.com/instrlabs/auth-service/internal"
	initx "github.com/instrlabs/shared/init"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := internal.LoadConfig()

	mongo := initx.NewMongo(&initx.MongoConfig{MongoURI: cfg.MongoURI, MongoDB: cfg.MongoDB})
	defer mongo.Close()

	userRepo := internal.NewUserRepository(mongo)
	userHandler := internal.NewUserHandler(cfg, userRepo)

	app := fiber.New(fiber.Config{})

	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app)
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{
		"/login",
		"/refresh",
		"/send-pin",
		"/check",
		"/google",
		"/google/callback",
	})

	app.Post("/login", userHandler.Login)
	app.Post("/logout", userHandler.Logout)
	app.Post("/refresh", userHandler.RefreshToken)
	app.Post("/send-pin", userHandler.SendPin)
	app.Post("/check", userHandler.CheckEmail)

	app.Get("/profile", userHandler.GetProfile)

	app.Get("/google", userHandler.GoogleLogin)
	app.Get("/google/callback", userHandler.GoogleCallback)

	log.Fatal(app.Listen(cfg.Port))
}
