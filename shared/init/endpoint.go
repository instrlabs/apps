package initx

import "github.com/gofiber/fiber/v2"

func SetupServiceSwagger(app *fiber.App) {
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Type("json").SendFile("./static/swagger.json")
	})
}

func SetupServiceHealth(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
