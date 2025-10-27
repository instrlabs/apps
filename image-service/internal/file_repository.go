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

type InstructionDetailRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

func (r *InstructionDetailRepository) ListUncleaned() []InstructionDetail {
	ctx := context.Background()
	cur, err := r.collection.Find(ctx, bson.M{"is_cleaned": false})
	if err != nil {
		log.Infof("instruction_detail_repository.ListUncleaned: Find failed: %v", err)
		return []InstructionDetail{}
	}
	defer cur.Close(ctx)

	var out []InstructionDetail
	for cur.Next(ctx) {
		var detail InstructionDetail
		if err := cur.Decode(&detail); err != nil {
			log.Infof("instruction_detail_repository.ListUncleaned: cursor decode failed: %v", err)
			continue
		}
		out = append(out, detail)
	}
	return out
}

func (r *InstructionDetailRepository) GetByID(id primitive.ObjectID) *InstructionDetail {
	ctx := context.Background()
	var detail InstructionDetail
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&detail); err != nil {
		log.Infof("instruction_detail_repository.GetByID: FindOne failed for id=%s: %v", id.Hex(), err)
		return nil
	}
	return &detail
}

func NewInstructionDetailRepository(db *initx.Mongo) *InstructionDetailRepository {
	return &InstructionDetailRepository{
		db:         db,
		collection: db.DB.Collection("image_details"),
	}
}

func (r *InstructionDetailRepository) CreateMany(details []*InstructionDetail) error {
	if len(details) == 0 {
		return nil
	}

	docs := make([]interface{}, 0, len(details))
	for _, detail := range details {
		docs = append(docs, detail)
	}

	_, err := r.collection.InsertMany(context.Background(), docs)
	if err != nil {
		log.Infof("instruction_detail_repository.CreateMany: InsertMany failed: %v", err)
	}

	return err
}

func (r *InstructionDetailRepository) ListByInstruction(instrID primitive.ObjectID) []InstructionDetail {
	ctx := context.Background()
	cur, err := r.collection.Find(ctx, bson.M{"instruction_id": instrID})
	if err != nil {
		log.Infof("instruction_detail_repository.ListByInstruction: Find failed for instruction_id=%s: %v", instrID.Hex(), err)
		return []InstructionDetail{}
	}
	defer cur.Close(ctx)

	var out []InstructionDetail
	for cur.Next(ctx) {
		var detail InstructionDetail
		if err := cur.Decode(&detail); err != nil {
			log.Infof("instruction_detail_repository.ListByInstruction: cursor decode failed for instruction_id=%s: %v", instrID.Hex(), err)
			continue
		}
		out = append(out, detail)
	}
	return out
}

func (r *InstructionDetailRepository) UpdateStatus(id primitive.ObjectID, st FileStatus) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"status":     st,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		log.Infof("instruction_detail_repository.UpdateStatus: UpdateByID failed for id=%s status=%v: %v", id.Hex(), st, err)
	}
	return err
}

func (r *InstructionDetailRepository) UpdateStatusAndSize(id primitive.ObjectID, st FileStatus, size int64) error {
	_, err := r.collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"status":     st,
			"file_size":  size,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		log.Infof("instruction_detail_repository.UpdateStatusAndSize: UpdateByID failed for id=%s status=%v size=%d: %v", id.Hex(), st, size, err)
	}
	return err
}

func (r *InstructionDetailRepository) ListOlderThan(before time.Time) []InstructionDetail {
	ctx := context.Background()
	filter := bson.M{
		"created_at": bson.M{"$lt": before},
		"is_cleaned": bson.M{"$ne": true},
	}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Infof("instruction_detail_repository.ListOlderThan: Find failed for before=%s: %v", before.UTC().Format(time.RFC3339), err)
		return []InstructionDetail{}
	}
	defer cur.Close(ctx)

	var out []InstructionDetail
	for cur.Next(ctx) {
		var detail InstructionDetail
		if err := cur.Decode(&detail); err != nil {
			log.Infof("instruction_detail_repository.ListOlderThan: cursor decode failed: %v", err)
			continue
		}
		out = append(out, detail)
	}
	return out
}

func (r *InstructionDetailRepository) DeleteMany(ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := r.collection.DeleteMany(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Infof("instruction_detail_repository.DeleteMany: DeleteMany failed for count=%d: %v", len(ids), err)
	}
	return err
}

func (r *InstructionDetailRepository) ListPendingUpdatedBefore(before time.Time) []InstructionDetail {
	ctx := context.Background()
	filter := bson.M{
		"status":     FileStatusPending,
		"updated_at": bson.M{"$lt": before},
		"is_cleaned": bson.M{"$ne": true},
	}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Infof("instruction_detail_repository.ListPendingUpdatedBefore: Find failed for before=%s: %v", before.UTC().Format(time.RFC3339), err)
		return []InstructionDetail{}
	}
	defer cur.Close(ctx)

	var out []InstructionDetail
	for cur.Next(ctx) {
		var detail InstructionDetail
		if err := cur.Decode(&detail); err != nil {
			log.Infof("instruction_detail_repository.ListPendingUpdatedBefore: cursor decode failed: %v", err)
			continue
		}
		out = append(out, detail)
	}
	return out
}

func (r *InstructionDetailRepository) MarkCleaned(ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := r.collection.UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"is_cleaned": true}})
	if err != nil {
		log.Infof("instruction_detail_repository.MarkCleaned: UpdateMany failed for count=%d: %v", len(ids), err)
	}
	return err
}
