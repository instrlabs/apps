package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Key         string             `bson:"key" json:"key"`
	Name        string             `bson:"name" json:"name"`
	Price       float64            `bson:"price" json:"price"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
	ProductType string             `bson:"productType,omitempty" json:"productType,omitempty"`
	UserID      string             `bson:"userId,omitempty" json:"userId,omitempty"`
	Active      bool               `bson:"active" json:"active"`
	IsFree      bool               `bson:"isFree" json:"isFree"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
