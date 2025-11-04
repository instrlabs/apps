package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Key       string             `json:"key" bson:"key"`
	Name      string             `json:"name" bson:"name"`
	Type      string             `json:"type" bson:"type"`
	Price     float64            `json:"price" bson:"price"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
