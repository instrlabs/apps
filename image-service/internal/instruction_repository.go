package internal

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/histweety-labs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstructionRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

func NewInstructionRepository(db *initx.Mongo) *InstructionRepository {
	return &InstructionRepository{
		db:         db,
		collection: db.DB.Collection("image_instructions"),
	}
}

func (r *InstructionRepository) Create(i *Instruction) error {
	_, err := r.collection.InsertOne(context.Background(), i)
	if err != nil {
		log.Errorf("Failed to insert instruction: %v", err)
		return err
	}
	return nil
}

func (r *InstructionRepository) GetByID(id primitive.ObjectID) *Instruction {
	var instruction Instruction
	_ = r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&instruction)
	return &instruction
}

func (r *InstructionRepository) UpdateStatus(id primitive.ObjectID, status InstructionStatus) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		log.Errorf("Failed to update instruction status: %v", err)
		return err
	}
	return nil
}
