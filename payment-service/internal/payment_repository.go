package internal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	paymentCollectionName = "payments"
)

// Payment represents a payment document in MongoDB
type Payment struct {
	ID            string        `bson:"_id,omitempty"`
	OrderID       string        `bson:"orderId"`
	UserID        string        `bson:"userId"`
	Amount        float64       `bson:"amount"`
	Currency      string        `bson:"currency"`
	PaymentMethod string        `bson:"paymentMethod"`
	Status        PaymentStatus `bson:"status"`
	RedirectURL   string        `bson:"redirectUrl,omitempty"`
	CreatedAt     time.Time     `bson:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt"`
}

// PaymentRepository handles payment data operations
type PaymentRepository struct {
	db         *MongoDB
	collection *mongo.Collection
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *MongoDB) *PaymentRepository {
	collection := db.GetCollection(paymentCollectionName)

	// Create indexes
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "orderId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("Error creating index: %v\n", err)
	}

	return &PaymentRepository{
		db:         db,
		collection: collection,
	}
}

// CreatePayment creates a new payment record
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *Payment) error {
	if payment.ID == "" {
		payment.ID = primitive.NewObjectID().Hex()
	}

	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID string) (*Payment, error) {
	var payment Payment

	err := r.collection.FindOne(ctx, bson.M{"orderId": orderID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by order ID: %w", err)
	}

	return &payment, nil
}

// UpdatePaymentStatus updates the status of a payment
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, orderID string, status PaymentStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"orderId": orderID}, update)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// ListPaymentsByUserID lists payments by user ID
func (r *PaymentRepository) ListPaymentsByUserID(ctx context.Context, userID string) ([]*Payment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var payments []*Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode payments: %w", err)
	}

	return payments, nil
}
