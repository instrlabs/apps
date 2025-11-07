package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/auth-service/internal"
	"github.com/instrlabs/shared/initx"
	"github.com/instrlabs/shared/middlewarex"
)

func main() {
	cfg := internal.LoadConfig()

	client, db := initx.NewMongo()
	defer initx.CloseMongo(client)

	userRepo := internal.NewUserRepository(db)
	userHandler := internal.NewUserHandler(cfg, userRepo)

	app := fiber.New(fiber.Config{})

	middlewarex.SetupPrometheus(app)
	middlewarex.SetupLogger(app)
	middlewarex.SetupServiceSwagger(app, "/auth")
	middlewarex.SetupServiceHealth(app)

	app.Post("/login", userHandler.Login)
	app.Post("/logout", userHandler.Logout)
	app.Post("/refresh", userHandler.RefreshToken)
	app.Post("/send-pin", userHandler.SendPin)

	app.Get("/profile", userHandler.GetProfile)

	app.Get("/google", userHandler.GoogleLogin)
	app.Get("/google/callback", userHandler.GoogleCallback)

	log.Fatal(app.Listen(cfg.Port))
}
