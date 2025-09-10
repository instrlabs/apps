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
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		log.Errorf("Failed to update instruction outputs: %v", err)
		return err
	}
	return nil
}

// ClaimNextPendingForRetry atomically claims one pending instruction for retry processing.
func (r *InstructionRepository) ClaimNextPendingForRetry(now time.Time, maxAgeMin, retryMax, lockTTLMin int) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"status":     "PENDING",
		"updated_at": bson.M{"$lt": now.Add(-time.Duration(maxAgeMin) * time.Minute)},
		"$or": []bson.M{
			{"retry_lock_until": bson.M{"$exists": false}},
			{"retry_lock_until": bson.M{"$eq": nil}},
			{"retry_lock_until": bson.M{"$lt": now}},
		},
		"$expr": bson.M{
			"$lt": []any{bson.M{"$ifNull": []any{"$retry_count", 0}}, retryMax},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"retry_lock_until": now.Add(time.Duration(lockTTLMin) * time.Minute),
			"updated_at":       now,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc bson.M
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	return doc, err
}

func (r *InstructionRepository) ReleaseRetryLock(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$unset": bson.M{"retry_lock_until": ""}})
	return err
}

func (r *InstructionRepository) MarkRetried(id primitive.ObjectID, now time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": id, "status": "PENDING"}, bson.M{
		"$inc": bson.M{"retry_count": 1},
		"$set": bson.M{
			"last_retry_at": now,
			"updated_at":    now,
		},
		"$unset": bson.M{"retry_lock_until": ""},
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		_, _ = r.collection.UpdateByID(ctx, id, bson.M{"$unset": bson.M{"retry_lock_until": ""}})
	}
	return nil
}
