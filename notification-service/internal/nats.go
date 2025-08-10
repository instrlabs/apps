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

// JobStatus represents the status of a job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type JobNotificationMessage struct {
	ID     string    `json:"id"`
	Status JobStatus `json:"status"`
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

type NotificationHandlerFunc func(ctx context.Context, notification *JobNotificationMessage) error

func (n *NatsService) SubscribeToJobNotifications(handler NotificationHandlerFunc) error {
	_, err := n.conn.Subscribe(n.cfg.NatsSubjectJobNotifications, func(msg *nats.Msg) {
		log.Printf("Received notification message from subject: %s", n.cfg.NatsSubjectJobNotifications)

		var notification JobNotificationMessage
		err := json.Unmarshal(msg.Data, &notification)
		if err != nil {
			log.Printf("Error unmarshaling notification: %v", err)
			return
		}

		ctx := context.Background()
		err = handler(ctx, &notification)
		if err != nil {
			log.Printf("Error processing notification: %v", err)
			return
		}

		log.Printf("Successfully processed notification for job: %s with status: %s",
			notification.ID, notification.Status)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to job notifications: %w", err)
	}

	log.Printf("Subscribed to NATS subject: %s", n.cfg.NatsSubjectJobNotifications)
	return nil
}
