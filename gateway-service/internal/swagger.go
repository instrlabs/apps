package internal

import (
	"github.com/gofiber/fiber/v2"
)

func SetupGatewaySwaggerUI(app *fiber.App) {
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Type("html").SendString(`
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Swagger UI - Multiple APIs</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
</head>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"></script>
<script>
	const ui = SwaggerUIBundle({
      urls: [
        { url: "/auth/swagger", name: "Auth API v1" },
        { url: "/products/swagger", name: "Products API v1" },
        { url: "/images/swagger", name: "Image API v1" },
        { url: "/pdfs/swagger", name: "PDF API v1" }
      ],
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      layout: "StandaloneLayout"
    });
</script>
</body>
</html>
		`)
	})
}
