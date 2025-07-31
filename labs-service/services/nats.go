package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"labs-service/constants"
)

// NatsService handles interactions with NATS
type NatsService struct {
	conn *nats.Conn
	cfg  *constants.Config
}

// JobMessage represents a message sent to NATS
type JobMessage struct {
	JobID string `json:"job_id"`
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

// PublishJobID publishes a job ID to NATS
func (n *NatsService) PublishJobID(ctx context.Context, jobID string) error {
	// Create the message
	msg := JobMessage{
		JobID: jobID,
	}

	// Convert to JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal job message: %w", err)
	}

	// Publish the message
	err = n.conn.Publish(n.cfg.NatsSubject, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish job message: %w", err)
	}

	return nil
}
