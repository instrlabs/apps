package internal

import (
	"context"
	"time"

	"github.com/instrlabs/shared/modelx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstructionDetailRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewInstructionDetailRepository(db *mongo.Database) *InstructionDetailRepository {
	return &InstructionDetailRepository{
		db:         db,
		collection: db.Collection("image_instruction_details"),
	}
}

func (r *InstructionDetailRepository) CreateMany(details []*modelx.InstructionFile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs := make([]interface{}, 0, len(details))
	for _, detail := range details {
		docs = append(docs, detail)
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

func (r *InstructionDetailRepository) GetByID(id primitive.ObjectID, detail modelx.InstructionFile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&detail)
	return err
}

func (r *InstructionDetailRepository) UpdateStatus(id primitive.ObjectID, newStatus modelx.InstructionDetailStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now().UTC(),
		},
	})

	return err
}

func (r *InstructionDetailRepository) UpdateFileSize(id primitive.ObjectID, fileSize int64) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"file_size":  fileSize,
			"updated_at": time.Now().UTC(),
		},
	})

	return err
}

func (r *InstructionDetailRepository) DeleteMany(ids []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

	return err
}
