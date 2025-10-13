package internal

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

func (r *FileRepository) ListUncleaned() []File {
	ctx := context.Background()
	cur, err := r.collection.Find(ctx, bson.M{"is_cleaned": false})
	if err != nil {
		log.Infof("file_repository.ListUncleaned: Find failed: %v", err)
		return []File{}
	}
	defer cur.Close(ctx)

	var out []File
	for cur.Next(ctx) {
		var f File
		if err := cur.Decode(&f); err != nil {
			log.Infof("file_repository.ListUncleaned: cursor decode failed: %v", err)
			continue
		}
		out = append(out, f)
	}
	return out
}

func (r *FileRepository) GetByID(id primitive.ObjectID) *File {
	ctx := context.Background()
	var f File
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&f); err != nil {
		log.Infof("file_repository.GetByID: FindOne failed for id=%s: %v", id.Hex(), err)
		return nil
	}
	return &f
}

func NewFileRepository(db *initx.Mongo) *FileRepository {
	return &FileRepository{
		db:         db,
		collection: db.DB.Collection("image_files"),
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
		log.Infof("file_repository.CreateMany: InsertMany failed: %v", err)
	}

	return err
}

func (r *FileRepository) ListByInstruction(instrID primitive.ObjectID) []File {
	ctx := context.Background()
	cur, err := r.collection.Find(ctx, bson.M{"instruction_id": instrID})
	if err != nil {
		log.Infof("file_repository.ListByInstruction: Find failed for instruction_id=%s: %v", instrID.Hex(), err)
		return []File{}
	}
	defer cur.Close(ctx)

	var out []File
	for cur.Next(ctx) {
		var f File
		if err := cur.Decode(&f); err != nil {
			log.Infof("file_repository.ListByInstruction: cursor decode failed for instruction_id=%s: %v", instrID.Hex(), err)
			continue
		}
		out = append(out, f)
	}
	return out
}

func (r *FileRepository) UpdateStatus(id primitive.ObjectID, st FileStatus) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"status":     st,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		log.Infof("file_repository.UpdateStatus: UpdateByID failed for id=%s status=%v: %v", id.Hex(), st, err)
	}
	return err
}

func (r *FileRepository) UpdateStatusAndSize(id primitive.ObjectID, st FileStatus, size int64) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"status":     st,
			"size":       size,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		log.Infof("file_repository.UpdateStatusAndSize: UpdateByID failed for id=%s status=%v size=%d: %v", id.Hex(), st, size, err)
	}
	return err
}

func (r *FileRepository) ListOlderThan(before time.Time) []File {
	ctx := context.Background()
	filter := bson.M{
		"created_at": bson.M{"$lt": before},
		"is_cleaned": bson.M{"$ne": true},
	}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Infof("file_repository.ListOlderThan: Find failed for before=%s: %v", before.UTC().Format(time.RFC3339), err)
		return []File{}
	}
	defer cur.Close(ctx)

	var out []File
	for cur.Next(ctx) {
		var f File
		if err := cur.Decode(&f); err != nil {
			log.Infof("file_repository.ListOlderThan: cursor decode failed: %v", err)
			continue
		}
		out = append(out, f)
	}
	return out
}

func (r *FileRepository) DeleteMany(ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := r.collection.DeleteMany(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Infof("file_repository.DeleteMany: DeleteMany failed for count=%d: %v", len(ids), err)
	}
	return err
}

func (r *FileRepository) ListPendingUpdatedBefore(before time.Time) []File {
	ctx := context.Background()
	filter := bson.M{
		"status":     FileStatusPending,
		"updated_at": bson.M{"$lt": before},
		"is_cleaned": bson.M{"$ne": true},
	}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Infof("file_repository.ListPendingUpdatedBefore: Find failed for before=%s: %v", before.UTC().Format(time.RFC3339), err)
		return []File{}
	}
	defer cur.Close(ctx)

	var out []File
	for cur.Next(ctx) {
		var f File
		if err := cur.Decode(&f); err != nil {
			log.Infof("file_repository.ListPendingUpdatedBefore: cursor decode failed: %v", err)
			continue
		}
		out = append(out, f)
	}
	return out
}

func (r *FileRepository) MarkCleaned(ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := r.collection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"is_cleaned": true}})
	if err != nil {
		log.Infof("file_repository.MarkCleaned: UpdateMany failed for count=%d: %v", len(ids), err)
	}
	return err
}
