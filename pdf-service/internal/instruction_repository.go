package internal

import (
	"context"
	"log"
	"time"

	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InstructionRepository struct {
	collection *mongo.Collection
}

func NewInstructionRepository(db *initx.Mongo) *InstructionRepository {
	return &InstructionRepository{
		collection: db.DB.Collection("pdf_instructions"),
	}
}

func (r *InstructionRepository) Create(instruction *Instruction) (*Instruction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instruction.ID = primitive.NewObjectID()
	instruction.CreatedAt = time.Now()
	instruction.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, instruction)
	if err != nil {
		log.Printf("Failed to create instruction: %v", err)
		return nil, err
	}

	return instruction, nil
}

func (r *InstructionRepository) GetByID(id primitive.ObjectID) (*Instruction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var instruction Instruction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&instruction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Failed to get instruction by ID %s: %v", id.Hex(), err)
		return nil, err
	}

	return &instruction, nil
}

func (r *InstructionRepository) ListLatest(userID string, limit int64) ([]Instruction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID}, findOptions)
	if err != nil {
		log.Printf("Failed to list instructions for user %s: %v", userID, err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var instructions []Instruction
	if err := cursor.All(ctx, &instructions); err != nil {
		log.Printf("Failed to decode instructions: %v", err)
		return nil, err
	}

	return instructions, nil
}
