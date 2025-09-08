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

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type JobRequestMessage struct {
	ID string `json:"id"`
}

type JobNotificationMessage struct {
	ID     string    `json:"id"`
	Status JobStatus `json:"status"`
	UserID string    `json:"userId,omitempty"`
}

func NewNatsService(cfg *Config) *NatsService {
	conn, err := nats.Connect(cfg.NatsURL, nats.Timeout(10*time.Second))
	if err != nil {
		_ = fmt.Errorf("failed to connect to NATS: %w", err)
		return nil
	}

	log.Println("Connected to NATS successfully")
	return &NatsService{conn: conn, cfg: cfg}
}

func (n *NatsService) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

func (n *NatsService) PublishJobRequest(jobID string, userID string) error {
	msg := &JobNotificationMessage{ID: jobID, UserID: userID}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return n.conn.Publish(n.cfg.NatsSubjectRequests, b)
}

type ImageJobHandler func(ctx context.Context, job *JobNotificationMessage) error

func (n *NatsService) Subscribe(handler ImageJobHandler) error {
	_, err := n.conn.Subscribe(n.cfg.NatsSubjectNotifications, func(msg *nats.Msg) {
		log.Printf("Received image job on subject: %s", n.cfg.NatsSubjectNotifications)
		var job JobNotificationMessage
		if err := json.Unmarshal(msg.Data, &job); err != nil {
			log.Printf("failed to unmarshal image job: %v", err)
			return
		}
		if err := handler(context.Background(), &job); err != nil {
			log.Printf("handler failed for job %s: %v", job.ID, err)
			//_ = n.PublishJobNotification(job.ID, JobStatusFailed, job.ID)
			return
		}
		log.Printf("processed image job %s", job.ID)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	log.Printf("Subscribed to NATS subject: %s", n.cfg.NatsSubjectRequests)
	return nil
}
