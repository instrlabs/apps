package main

import (
	"log"

	"github.com/arthadede/auth-service/internal"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := internal.LoadConfig()

	mongo := internal.NewMongoDB(config)
	defer mongo.Close()

	userRepo := internal.NewUserRepository(mongo)
	userController := internal.NewUserController(userRepo, config)
	userHandler := internal.NewUserHandler(userController, config)

	app := fiber.New(fiber.Config{})

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Type("json").SendFile("./static/swagger.json")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

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
