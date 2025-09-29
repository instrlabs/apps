package internal

import (
	"context"
	"log"

	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

func (r *FileRepository) GetByID(id primitive.ObjectID) *File {
	ctx := context.Background()
	var f File
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&f); err != nil {
		log.Printf("file_repository.GetByID: FindOne failed for id=%s: %v", id.Hex(), err)
		return nil
	}
	return &f
}

func NewFileRepository(db *initx.Mongo) *FileRepository {
	return &FileRepository{
		db:         db,
		collection: db.DB.Collection("images_files"),
	}
}

func (r *FileRepository) CreateMany(files []*File) error {
	if len(files) == 0 {
		return nil
	}

	docs := make([]interface{}, 0, len(files))
	for _, f := range files {
		docs = append(docs, f)
	}

	_, err := r.collection.InsertMany(context.Background(), docs)
	if err != nil {
		log.Printf("file_repository.CreateMany: InsertMany failed: %v", err)
	}

	return err
}

func (r *FileRepository) CreateOne(f *File) error {
	_, err := r.collection.InsertOne(context.Background(), f)
	if err != nil {
		log.Printf("file_repository.CreateOne: InsertOne failed: %v", err)
		return err
	}

	return nil
}

func (r *FileRepository) ListByInstruction(instrID primitive.ObjectID) []File {
	ctx := context.Background()
	cur, err := r.collection.Find(ctx, bson.M{"instruction_id": instrID})
	if err != nil {
		log.Printf("file_repository.ListByInstruction: Find failed for instruction_id=%s: %v", instrID.Hex(), err)
		return []File{}
	}
	defer cur.Close(ctx)

	var out []File
	for cur.Next(ctx) {
		var f File
		if err := cur.Decode(&f); err != nil {
			log.Printf("file_repository.ListByInstruction: cursor decode failed for instruction_id=%s: %v", instrID.Hex(), err)
			continue
		}
		out = append(out, f)
	}
	return out
}

func (r *FileRepository) UpdateStatus(id primitive.ObjectID, st FileStatus) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"status": st,
		},
	})
	if err != nil {
		log.Printf("file_repository.UpdateStatus: UpdateByID failed for id=%s status=%v: %v", id.Hex(), st, err)
	}
	return err
}

// UpdateStatusAndSize updates both the status and size fields for a file in a single operation.
func (r *FileRepository) UpdateStatusAndSize(id primitive.ObjectID, st FileStatus, size int64) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"status": st,
			"size":   size,
		},
	})
	if err != nil {
		log.Printf("file_repository.UpdateStatusAndSize: UpdateByID failed for id=%s status=%v size=%d: %v", id.Hex(), st, size, err)
	}
	return err
}
