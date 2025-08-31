package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	nats "github.com/nats-io/nats.go"
)

type NatsService struct {
	conn *nats.Conn
	cfg  *Config
}

// PaymentEventMessage represents a payment event message
type PaymentEventMessage struct {
	ID            string        `json:"id"`
	OrderID       string        `json:"orderId"`
	UserID        string        `json:"userId"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	PaymentMethod string        `json:"paymentMethod"`
	Status        PaymentStatus `json:"status"`
	RedirectURL   string        `json:"redirectUrl,omitempty"`
	Timestamp     time.Time     `json:"timestamp"`
	Type          PaymentType   `json:"type,omitempty"`
}

// PaymentRequestMessage represents a payment request message
type PaymentRequestMessage struct {
	OrderID       string      `json:"orderId"`
	UserID        string      `json:"userId"`
	Amount        float64     `json:"amount"`
	Currency      string      `json:"currency"`
	PaymentMethod string      `json:"paymentMethod,omitempty"`
	Description   string      `json:"description,omitempty"`
	CallbackURL   string      `json:"callbackUrl,omitempty"`
	Type          PaymentType `json:"type,omitempty"`
}

func NewNatsService(cfg *Config) *NatsService {
	conn, err := nats.Connect(cfg.NatsURL, nats.Timeout(10*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	log.Println("Connected to NATS successfully")

	return &NatsService{
		conn: conn,
		cfg:  cfg,
	}
}

func (n *NatsService) Close() {
	if n.conn != nil {
		n.conn.Close()
		log.Println("Disconnected from NATS")
	}
}

// PublishPaymentEvent publishes a payment event to NATS
func (n *NatsService) PublishPaymentEvent(event *PaymentEventMessage) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal payment event: %w", err)
	}

	err = n.conn.Publish(n.cfg.NatsSubjectPaymentEvents, data)
	if err != nil {
		return fmt.Errorf("failed to publish payment event: %w", err)
	}

	log.Printf("Published payment event for order %s with status %s", event.OrderID, event.Status)
	return nil
}

// PaymentRequestHandlerFunc is a function that handles payment requests
type PaymentRequestHandlerFunc func(ctx context.Context, request *PaymentRequestMessage) (*PaymentEventMessage, error)

// SubscribeToPaymentRequests subscribes to payment requests
func (n *NatsService) SubscribeToPaymentRequests(handler PaymentRequestHandlerFunc) error {
	_, err := n.conn.Subscribe(n.cfg.NatsSubjectPaymentRequests, func(msg *nats.Msg) {
		log.Printf("Received payment request from subject: %s", n.cfg.NatsSubjectPaymentRequests)

		var request PaymentRequestMessage
		err := json.Unmarshal(msg.Data, &request)
		if err != nil {
			log.Printf("Error unmarshaling payment request: %v", err)
			return
		}

		ctx := context.Background()
		response, err := handler(ctx, &request)
		if err != nil {
			log.Printf("Error processing payment request: %v", err)
			return
		}

		// Publish the payment event
		err = n.PublishPaymentEvent(response)
		if err != nil {
			log.Printf("Error publishing payment event: %v", err)
			return
		}

		log.Printf("Successfully processed payment request for order: %s", request.OrderID)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to payment requests: %w", err)
	}

	log.Printf("Subscribed to NATS subject: %s", n.cfg.NatsSubjectPaymentRequests)
	return nil
}
