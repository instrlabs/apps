package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Key       string             `json:"key" bson:"key"`
	Name      string             `json:"name" bson:"name"`
	Desc      string             `json:"desc" bson:"desc"`
	Type      string             `json:"type" bson:"type"`
	Price     float64            `json:"price" bson:"price"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
