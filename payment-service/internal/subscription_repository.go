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
	subscriptionCollectionName = "subscriptions"
)

type SubscriptionStatus string

const (
	SubscriptionPending   SubscriptionStatus = "pending"
	SubscriptionActivated SubscriptionStatus = "activated"
	SubscriptionCancelled SubscriptionStatus = "cancelled"
	SubscriptionSuspended SubscriptionStatus = "suspended"
)

type Subscription struct {
	ID             string             `bson:"_id,omitempty"`
	UserID         string             `bson:"userId"`
	PaymentID      string             `bson:"paymentId"`
	SubscribedAt   time.Time          `bson:"subscribedAt"`
	UnsubscribedAt time.Time          `bson:"unsubscribedAt"`
	NextPaymentAt  time.Time          `bson:"nextPaymentAt"`
	Status         SubscriptionStatus `bson:"status"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

type SubscriptionRepository struct {
	db         *MongoDB
	collection *mongo.Collection
}

func NewSubscriptionRepository(db *MongoDB) *SubscriptionRepository {
	collection := db.GetCollection(subscriptionCollectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "subscriptionId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("Error creating index on subscriptions: %v\n", err)
	}

	return &SubscriptionRepository{db: db, collection: collection}
}

// CreateSubscription creates a new subscription record
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, s *Subscription) error {
	if s.ID == "" {
		s.ID = primitive.NewObjectID().Hex()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, s)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

// GetBySubscriptionID retrieves a subscription by its provider ID
func (r *SubscriptionRepository) GetBySubscriptionID(ctx context.Context, id string) (*Subscription, error) {
	var s Subscription
	err := r.collection.FindOne(ctx, bson.M{"subscriptionId": id}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription by id: %w", err)
	}
	return &s, nil
}

// UpdateSubscriptionStatus updates the status of a subscription
func (r *SubscriptionRepository) UpdateSubscriptionStatus(ctx context.Context, id string, status PaymentStatus) error {
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"subscriptionId": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}
	return nil
}

// ListSubscriptionsByUserID lists subscriptions associated with a user
func (r *SubscriptionRepository) ListSubscriptionsByUserID(ctx context.Context, userID string) ([]*Subscription, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var subs []*Subscription
	if err := cursor.All(ctx, &subs); err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions: %w", err)
	}
	return subs, nil
}
