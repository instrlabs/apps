package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Key         string             `bson:"key" json:"key"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	ProductType string             `bson:"product_type" json:"product_type"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	IsFree      bool               `bson:"is_free" json:"is_free"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
