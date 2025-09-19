package internal

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/histweety-labs/shared/init"
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

func (r *InstructionRepository) UpdateOutputs(id primitive.ObjectID, outputs []File) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"outputs":    outputs,
			"status":     InstructionStatusCompleted,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		log.Errorf("Failed to update instruction outputs: %v", err)
		return err
	}
	return nil
}

func (r *InstructionRepository) ListByUser(userID primitive.ObjectID) []Instruction {
	ctx := context.Background()
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cur, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		log.Errorf("Failed to list instructions: %v", err)
		return []Instruction{}
	}
	defer cur.Close(ctx)

	var res []Instruction
	for cur.Next(ctx) {
		var ins Instruction
		if err := cur.Decode(&ins); err != nil {
			log.Errorf("Failed to decode instruction: %v", err)
			continue
		}
		res = append(res, ins)
	}
	return res
}

// ListOldCompleted returns completed instructions created before a cutoff time, limited by `limit`.
func (r *InstructionRepository) ListOldCompleted(before time.Time, limit int64) []Instruction {
	ctx := context.Background()
	filter := bson.M{
		"status":     InstructionStatusCompleted,
		"created_at": bson.M{"$lt": before},
	}
	opts := options.Find().SetLimit(limit)
	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Errorf("Failed to list old completed instructions: %v", err)
		return []Instruction{}
	}
	defer cur.Close(ctx)

	var res []Instruction
	for cur.Next(ctx) {
		var ins Instruction
		if err := cur.Decode(&ins); err != nil {
			log.Errorf("Failed to decode instruction: %v", err)
			continue
		}
		res = append(res, ins)
	}
	return res
}

// DeleteByID deletes an instruction document by id.
func (r *InstructionRepository) DeleteByID(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Errorf("Failed to delete instruction: %v", err)
		return err
	}
	return nil
}
