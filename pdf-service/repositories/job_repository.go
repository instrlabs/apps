package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"pdf-service/models"
)

type JobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(db *mongo.Database) *JobRepository {
	return &JobRepository{
		collection: db.Collection("jobs"),
	}
}

func (r *JobRepository) Create(ctx context.Context, job *models.Job) (*models.Job, error) {
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}

	job.ID = result.InsertedID.(primitive.ObjectID)
	return job, nil
}

func (r *JobRepository) FindByID(ctx context.Context, id string) (*models.Job, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var job models.Job
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id string, status models.JobStatus, errorMsg string, s3path ...string) (*models.Job, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	if errorMsg != "" {
		update["$set"].(bson.M)["error"] = errorMsg
	}

	if len(s3path) > 0 && s3path[0] != "" {
		update["$set"].(bson.M)["s3path"] = s3path[0]
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedJob models.Job
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updatedJob)
	if err != nil {
		return nil, err
	}

	return &updatedJob, nil
}
