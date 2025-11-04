package internal

import (
	"context"
	"log"
	"time"

	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstructionDetailRepository struct {
	collection *mongo.Collection
}

func NewInstructionDetailRepository(db *initx.Mongo) *InstructionDetailRepository {
	return &InstructionDetailRepository{
		collection: db.DB.Collection("pdf_instruction_details"),
	}
}

func (r *InstructionDetailRepository) CreateMany(details []InstructionDetail) ([]InstructionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs := make([]interface{}, len(details))
	for i := range details {
		details[i].ID = primitive.NewObjectID()
		details[i].CreatedAt = time.Now()
		details[i].UpdatedAt = time.Now()
		docs[i] = details[i]
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("Failed to create instruction details: %v", err)
		return nil, err
	}

	return details, nil
}

func (r *InstructionDetailRepository) GetByID(id primitive.ObjectID) (*InstructionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var detail InstructionDetail
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&detail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Failed to get instruction detail by ID %s: %v", id.Hex(), err)
		return nil, err
	}

	return &detail, nil
}

func (r *InstructionDetailRepository) ListByInstruction(instructionID primitive.ObjectID) ([]InstructionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"instruction_id": instructionID})
	if err != nil {
		log.Printf("Failed to list instruction details for instruction %s: %v", instructionID.Hex(), err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var details []InstructionDetail
	if err := cursor.All(ctx, &details); err != nil {
		log.Printf("Failed to decode instruction details: %v", err)
		return nil, err
	}

	return details, nil
}

func (r *InstructionDetailRepository) UpdateStatus(id primitive.ObjectID, status FileStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Printf("Failed to update status for instruction detail %s: %v", id.Hex(), err)
		return err
	}

	return nil
}

func (r *InstructionDetailRepository) UpdateStatusAndSize(id primitive.ObjectID, status FileStatus, fileSize int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"file_size":  fileSize,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Printf("Failed to update status and file size for instruction detail %s: %v", id.Hex(), err)
		return err
	}

	return nil
}

func (r *InstructionDetailRepository) ListOlderThan(olderThan time.Time) ([]InstructionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"status": FileStatusDone,
		"updated_at": bson.M{
			"$lt": olderThan,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list files older than %v: %v", olderThan, err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var details []InstructionDetail
	if err := cursor.All(ctx, &details); err != nil {
		log.Printf("Failed to decode older files: %v", err)
		return nil, err
	}

	return details, nil
}

func (r *InstructionDetailRepository) ListPendingUpdatedBefore(before time.Time) ([]InstructionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"status": FileStatusProcessing,
		"updated_at": bson.M{
			"$lt": before,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list stale pending files: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var details []InstructionDetail
	if err := cursor.All(ctx, &details); err != nil {
		log.Printf("Failed to decode stale pending files: %v", err)
		return nil, err
	}

	return details, nil
}

func (r *InstructionDetailRepository) MarkCleaned(ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectIDs := make([]interface{}, len(ids))
	for i, id := range ids {
		objectIDs[i] = id
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{
		"$set": bson.M{
			"status":     "CLEANED",
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to mark files as cleaned: %v", err)
		return err
	}

	return nil
}

func (r *InstructionDetailRepository) ListUncleaned() ([]InstructionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"status": bson.M{"$ne": "CLEANED"}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list uncleaned files: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var details []InstructionDetail
	if err := cursor.All(ctx, &details); err != nil {
		log.Printf("Failed to decode uncleaned files: %v", err)
		return nil, err
	}

	return details, nil
}
