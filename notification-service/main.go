package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/instrlabs/notification-service/internal"
	initx "github.com/instrlabs/shared/init"
	natsgo "github.com/nats-io/nats.go"
)

func main() {
	cfg := internal.NewConfig()

	natsSrv := initx.NewNats(cfg.NatsURL)
	defer natsSrv.Close()

	sseService := internal.NewSSEService(cfg)

	_, _ = natsSrv.Conn.Subscribe(cfg.NatsSubjectNotificationsSSE, func(m *natsgo.Msg) {
		sseService.NotificationUser(m.Data)
	})

	app := fiber.New(fiber.Config{})
	initx.SetupLogger(app)
	initx.SetupServiceHealth(app)
	internal.SetupMiddleware(app, cfg)

	app.Get("/sse", sseService.HandleSSE)

	log.Println(app.Listen(cfg.Port))
}
