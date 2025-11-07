package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/notification-service/internal"
	"github.com/instrlabs/shared/initx"
	"github.com/instrlabs/shared/middlewarex"
	natsgo "github.com/nats-io/nats.go"
)

func main() {
	cfg := internal.LoadConfig()

	natsSrv := initx.NewNats(cfg.NatsURI)
	defer natsSrv.Close()

	sseService := internal.NewSSEService(cfg)

	_, _ = natsSrv.Subscribe(cfg.NatsSubjectNotificationsSSE, func(m *natsgo.Msg) {
		sseService.NotificationUser(m.Data)
	})

	app := fiber.New(fiber.Config{})

	middlewarex.SetupPrometheus(app)
	middlewarex.SetupServiceHealth(app)
	middlewarex.SetupLogger(app)
	internal.SetupMiddleware(app, cfg)
	middlewarex.SetupAuthentication(app)
	internal.SetupNotificationRoutes(app, cfg, natsSrv)

	app.Get("/sse", sseService.HandleSSE)

	log.Fatal(app.Listen(cfg.Port))
}
