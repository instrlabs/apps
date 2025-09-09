package internal

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/histweety-labs/shared/init"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

func NewFileRepository(db *initx.Mongo) *FileRepository {
	return &FileRepository{
		db:         db,
		collection: db.DB.Collection("files"),
	}
}

func (r *FileRepository) Create(f *File) error {
	_, err := r.collection.InsertOne(context.Background(), f)
	if err != nil {
		log.Errorf("Failed to insert file: %v", err)
		return err
	}
	return nil
}
