package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	conn *nats.Conn
	cfg  *Config
}

type JobMessage struct {
	JobID string `json:"job_id"`
}

func NewNatsService(cfg *Config) (*NatsService, error) {
	conn, err := nats.Connect(cfg.NatsURL, nats.Timeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Println("Connected to NATS successfully")

	return &NatsService{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (n *NatsService) Close() {
	if n.conn != nil {
		n.conn.Close()
		log.Println("Disconnected from NATS")
	}
}

type PDFNotificationHandler func(ctx context.Context, job *JobNotificationMessage) error

func (n *NatsService) SubscribeToPDFNotification(handler PDFNotificationHandler) error {
	_, err := n.conn.Subscribe(n.cfg.NatsSubjectJobNotifications, func(msg *nats.Msg) {
		log.Printf("Received job message from subject: %s", n.cfg.NatsSubjectJobNotifications)

		var job JobNotificationMessage
		err := json.Unmarshal(msg.Data, &job)
		if err != nil {
			log.Printf("Error unmarshaling job: %v", err)
			return
		}

		ctx := context.Background()
		err = handler(ctx, &job)
		if err != nil {
			log.Printf("Error processing job: %v", err)
			return
		}

		log.Printf("Successfully processed job: %s", job.ID)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to jobs: %w", err)
	}

	log.Printf("Subscribed to NATS subject: %s", n.cfg.NatsSubjectJobNotifications)
	return nil
}
