package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"pdf-service/constants"
	"pdf-service/models"
)

type NatsService struct {
	conn *nats.Conn
	cfg  *constants.Config
}

type JobMessage struct {
	JobID string `json:"job_id"`
}

func NewNatsService(cfg *constants.Config) (*NatsService, error) {
	// Connect to NATS
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

func (n *NatsService) PublishPDFJob(job interface{}) error {
	msgBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal PDF job: %w", err)
	}

	err = n.conn.Publish(models.PDFJobSubject, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish PDF job: %w", err)
	}

	log.Printf("Published PDF job to NATS subject: %s", models.PDFJobSubject)
	return nil
}

type PDFJobHandler func(ctx context.Context, job *models.PDFJob) error

func (n *NatsService) SubscribeToPDFJobs(handler PDFJobHandler) error {
	_, err := n.conn.Subscribe(models.PDFJobSubject, func(msg *nats.Msg) {
		log.Printf("Received PDF job message from subject: %s", models.PDFJobSubject)

		var job models.PDFJob
		err := json.Unmarshal(msg.Data, &job)
		if err != nil {
			log.Printf("Error unmarshaling PDF job: %v", err)
			return
		}

		ctx := context.Background()
		err = handler(ctx, &job)
		if err != nil {
			log.Printf("Error processing PDF job: %v", err)
			return
		}

		log.Printf("Successfully processed PDF job: %s", job.JobID)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to PDF jobs: %w", err)
	}

	log.Printf("Subscribed to NATS subject: %s", models.PDFJobSubject)
	return nil
}
