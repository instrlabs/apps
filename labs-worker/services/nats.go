package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"labs-worker/constants"
)

type JobMessage struct {
	JobID string `json:"job_id"`
}

type CompletionMessage struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type NatsService struct {
	conn *nats.Conn
	cfg  *constants.Config
}

func NewNatsService(cfg *constants.Config) (*NatsService, error) {
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

func (n *NatsService) SubscribeToPDFJobs(handler func(jobID string)) error {
	_, err := n.conn.Subscribe(n.cfg.NatsSubjectPDFJobs, func(msg *nats.Msg) {
		var jobMsg JobMessage
		err := json.Unmarshal(msg.Data, &jobMsg)
		if err != nil {
			log.Printf("Failed to unmarshal job message: %v", err)
			return
		}

		handler(jobMsg.JobID)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to jobs: %w", err)
	}

	log.Printf("Subscribed to %s", n.cfg.NatsSubjectPDFJobs)
	return nil
}

func (n *NatsService) PublishCompletionToPDFResults(jobID string, status string, err error) error {
	msg := CompletionMessage{
		JobID:  jobID,
		Status: status,
	}

	if err != nil {
		msg.Error = err.Error()
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal completion message: %w", err)
	}

	err = n.conn.Publish(n.cfg.NatsSubjectPDFResults, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish completion message: %w", err)
	}

	return nil
}
