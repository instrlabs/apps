package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	db         *MongoDB
	collection *mongo.Collection
}

func NewFileRepository(db *MongoDB) *FileRepository {
	return &FileRepository{
		db:         db,
		collection: db.DB.Collection("files"),
	}
}

func (r *FileRepository) Create(f *File) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(context.Background(), f)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}
	return primitive.NilObjectID, nil
}
