package internal

import (
	"context"
	"time"

	"github.com/instrlabs/shared/modelx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstructionRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewInstructionRepository(db *mongo.Database) *InstructionRepository {
	return &InstructionRepository{
		db:         db,
		collection: db.Collection("image_instructions"),
	}
}

func (r *InstructionRepository) Create(instr *modelx.Instruction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, instr)

	return err
}

func (r *InstructionRepository) GetByID(id primitive.ObjectID, instr modelx.Instruction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&instr)

	return err
}
