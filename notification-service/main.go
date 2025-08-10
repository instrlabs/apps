package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthadede/notification-service/internal"
)

func main() {
	cfg := internal.NewConfig()

	natsService, err := internal.NewNatsService(cfg)
	if err != nil {
		log.Fatalf("Failed to create NATS service: %v", err)
	}
	defer natsService.Close()

	wsService := internal.NewWebSocketService(cfg)

	mux := http.NewServeMux()

	mux.HandleFunc(cfg.WebSocketPath, wsService.HandleWebSocket)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}

	err = natsService.SubscribeToJobNotifications(func(ctx context.Context, notification *internal.JobNotificationMessage) error {
		log.Printf("Received job notification: %+v", notification)
		return wsService.BroadcastNotification(ctx, notification)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to job notifications: %v", err)
	}

	go func() {
		log.Printf("Starting notification service on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
