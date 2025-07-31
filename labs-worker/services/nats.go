package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"labs-worker/constants"
	"labs-worker/models"
)

// NatsService handles interactions with NATS
type NatsService struct {
	conn *nats.Conn
	cfg  *constants.Config
}

// NewNatsService creates a new NatsService
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

// Close closes the NATS connection
func (n *NatsService) Close() {
	if n.conn != nil {
		n.conn.Close()
		log.Println("Disconnected from NATS")
	}
}

func (n *NatsService) SubscribeToJobs(handler func(jobID string)) error {
	_, err := n.conn.Subscribe(n.cfg.NatsSubjectPDF, func(msg *nats.Msg) {
		// Parse the message
		var jobMsg models.JobMessage
		err := json.Unmarshal(msg.Data, &jobMsg)
		if err != nil {
			log.Printf("Failed to unmarshal job message: %v", err)
			return
		}

		// Call the handler
		handler(jobMsg.JobID)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to jobs: %w", err)
	}

	log.Printf("Subscribed to %s", n.cfg.NatsSubjectPDF)
	return nil
}

// PublishCompletion publishes a completion event to NATS
func (n *NatsService) PublishCompletion(jobID string, status string, err error) error {
	// Create the message
	msg := models.CompletionMessage{
		JobID:  jobID,
		Status: status,
	}

	if err != nil {
		msg.Error = err.Error()
	}

	// Convert to JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal completion message: %w", err)
	}

	// Publish the message
	err = n.conn.Publish(n.cfg.NatsSubjectResults, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish completion message: %w", err)
	}

	return nil
}
