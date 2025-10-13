package internal

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *InstructionRepository) ListLatest(limit int64) ([]Instruction, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
	cursor, err := r.collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		log.Errorf("Failed to list instructions: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var instructions []Instruction
	for cursor.Next(context.Background()) {
		var instr Instruction
		if err := cursor.Decode(&instr); err != nil {
			log.Errorf("Failed to decode instruction: %v", err)
			continue
		}
		instructions = append(instructions, instr)
	}

	return instructions, nil
}
