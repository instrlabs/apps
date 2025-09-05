package internal

import "go.mongodb.org/mongo-driver/mongo"

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
