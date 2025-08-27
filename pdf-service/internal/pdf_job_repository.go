package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const pdfJobsCollection = "pdf_jobs"

type PDFJobRepository struct {
	collection *mongo.Collection
}

func NewPDFJobRepository(db *mongo.Database) *PDFJobRepository {
	return &PDFJobRepository{
		collection: db.Collection(pdfJobsCollection),
	}
}

func (r *PDFJobRepository) Create(ctx context.Context, job *PDFJob) error {
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		return err
	}

	job.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PDFJobRepository) FindByID(ctx context.Context, id string) (*PDFJob, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var job PDFJob
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *PDFJobRepository) FindAll(ctx context.Context, limit, offset int64) ([]*PDFJob, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*PDFJob
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *PDFJobRepository) Update(ctx context.Context, id string, update *UpdatePDFJobRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateDoc := bson.M{"updated_at": time.Now()}

	if update.OutputFilePath != "" {
		updateDoc["output_file_path"] = update.OutputFilePath
	}

	if update.Status != "" {
		updateDoc["status"] = update.Status
	}

	if update.Error != "" {
		updateDoc["error"] = update.Error
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateDoc},
	)

	return err
}
