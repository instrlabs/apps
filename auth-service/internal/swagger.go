package internal

import "github.com/gofiber/fiber/v2"

func SetupSwagger(app *fiber.App) {
	app.Static("/swagger", "./docs")

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/swagger.yaml", fiber.StatusTemporaryRedirect)
	})
}
