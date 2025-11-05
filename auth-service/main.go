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
	sessionRepo := internal.NewUserSessionRepository(db)
	userHandler := internal.NewUserHandler(cfg, userRepo, sessionRepo)

	app := fiber.New(fiber.Config{})

	middlewarex.SetupPrometheus(app)
	middlewarex.SetupLogger(app)
	middlewarex.SetupServiceSwagger(app, "/auth")
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupAuthentication(app, []string{
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

	app.Get("/profile", userHandler.GetProfile)

	app.Get("/google", userHandler.GoogleLogin)
	app.Get("/google/callback", userHandler.GoogleCallback)

	app.Get("/devices", userHandler.GetDevices)
	app.Post("/devices/:sessionId/revoke", userHandler.RevokeDevice)
	app.Post("/devices/revoke-all", userHandler.LogoutAllDevices)

	log.Fatal(app.Listen(cfg.Port))
}
