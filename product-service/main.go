package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	initx "github.com/instrlabs/shared/init"

	"github.com/instrlabs/product-service/internal"
)

func main() {
	cfg := internal.LoadConfig()

	mongo := initx.NewMongo(&initx.MongoConfig{
		MongoURI: cfg.MongoURI,
		MongoDB:  cfg.MongoDB,
	})
	defer mongo.Close()

	app := fiber.New(fiber.Config{})

	initx.SetupPrometheus(app)
	initx.SetupLogger(app)
	initx.SetupServiceSwagger(app, cfg.ApiUrl, "/products")
	initx.SetupServiceHealth(app)
	initx.SetupAuthenticated(app, []string{
		"/",
	})

	productRepo := internal.NewProductRepository(mongo)
	productHandler := internal.NewProductHandler(productRepo)

	app.Get("/", productHandler.ListProducts)

	log.Fatal(app.Listen(cfg.Port))
}
