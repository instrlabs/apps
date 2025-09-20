package main

import (
	"log"

	"github.com/arthadede/notification-service/internal"
	"github.com/gofiber/fiber/v2"
	initx "github.com/histweety-labs/shared/init"
	natsgo "github.com/nats-io/nats.go"
)

func main() {
	cfg := internal.NewConfig()

	natsSrv := initx.NewNats(cfg.NatsURL)
	defer natsSrv.Close()

	sseService := internal.NewSSEService(cfg)

	_, _ = natsSrv.Conn.Subscribe(cfg.NatsSubjectNotificationsSSE, func(m *natsgo.Msg) {
		sseService.Broadcast(m.Data)
	})

	app := fiber.New(fiber.Config{})
	initx.SetupLogger(app)
	initx.SetupServiceHealth(app)
	internal.SetupMiddleware(app)

	app.Get("/sse", sseService.HandleSSE)

	log.Println(app.Listen(cfg.Port))
}
