package internal

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstructionRepository struct {
	db         *MongoDB
	collection *mongo.Collection
}

func NewInstructionRepository(db *MongoDB) *InstructionRepository {
	return &InstructionRepository{
		db:         db,
		collection: db.DB.Collection("image_instructions"),
	}
}

func (r *InstructionRepository) Create(i *Instruction) interface{} {
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now().UTC()
	}
	if i.UpdatedAt.IsZero() {
		i.UpdatedAt = i.CreatedAt
	}

	res, err := r.collection.InsertOne(context.Background(), i)

	if err != nil {
		log.Errorf("Failed to insert instruction: %v", err)
		return nil
	}

	return res.InsertedID
}
