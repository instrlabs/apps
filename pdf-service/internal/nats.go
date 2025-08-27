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

func (n *NatsService) PublishJobID(subject string, jobID string) error {
	msg := JobMessage{
		JobID: jobID,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal job message: %w", err)
	}

	err = n.conn.Publish(subject, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish job message: %w", err)
	}

	return nil
}

func (n *NatsService) PublishPDFJob(pdfJobID string) error {
	msgBytes, err := json.Marshal(&PDFJobMessage{ID: pdfJobID})
	if err != nil {
		return fmt.Errorf("failed to marshal PDF job: %w", err)
	}

	err = n.conn.Publish(n.cfg.NatsSubjectPDFJobs, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish PDF job: %w", err)
	}

	log.Printf("Published PDF job to NATS subject: %s", n.cfg.NatsSubjectPDFJobs)
	return nil
}

func (n *NatsService) PublishJobNotification(jobID string, status JobStatus) error {
	msgBytes, err := json.Marshal(&JobNotificationMessage{ID: jobID, Status: status})
	if err != nil {
		return fmt.Errorf("failed to marshal job notification: %w", err)
	}

	err = n.conn.Publish(n.cfg.NatsSubjectJobNotifications, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish job notification: %w", err)
	}

	log.Printf("Published job to NATS subject: %s", n.cfg.NatsSubjectJobNotifications)
	return nil
}

type PDFJobHandlerFunc func(ctx context.Context, job *PDFJobMessage) error

func (n *NatsService) SubscribeToPDFJob(handler PDFJobHandlerFunc) error {
	_, err := n.conn.Subscribe(n.cfg.NatsSubjectPDFJobs, func(msg *nats.Msg) {
		log.Printf("Received job message from subject: %s", n.cfg.NatsSubjectPDFJobs)

		var job PDFJobMessage
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

	log.Printf("Subscribed to NATS subject: %s", n.cfg.NatsSubjectPDFJobs)
	return nil
}
