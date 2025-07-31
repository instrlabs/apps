package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"labs-worker/constants"
	"labs-worker/models"
)

// DBService handles interactions with the database
type DBService struct {
	client     *mongo.Client
	db         *mongo.Database
	pdfJobColl *mongo.Collection
	cfg        *constants.Config
}

// NewDBService creates a new DBService
func NewDBService(cfg *constants.Config) (*DBService, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully")

	// Get database and collections
	db := client.Database(cfg.MongoDB)
	pdfJobColl := db.Collection("pdf_jobs")

	return &DBService{
		client:     client,
		db:         db,
		pdfJobColl: pdfJobColl,
		cfg:        cfg,
	}, nil
}

// Close closes the database connection
func (d *DBService) Close() {
	if d.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		d.client.Disconnect(ctx)
		log.Println("Disconnected from MongoDB")
	}
}

// GetPDFJobByID retrieves a PDF job by its job ID
func (d *DBService) GetPDFJobByID(ctx context.Context, jobID string) (*models.PDFJob, error) {
	// Find the job
	var job models.PDFJob
	err := d.pdfJobColl.FindOne(ctx, bson.M{"job_id": jobID}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("job not found: %s", jobID)
		}
		return nil, fmt.Errorf("failed to find job: %w", err)
	}

	return &job, nil
}

// UpdateJobStatus updates the status of a PDF job
func (d *DBService) UpdateJobStatus(ctx context.Context, jobID string, status models.JobStatus, outputPath string, errMsg string) error {
	// Create the update
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	// Add output path if provided
	if outputPath != "" {
		update["$set"].(bson.M)["output_file_path"] = outputPath
	}

	// Add error message if provided
	if errMsg != "" {
		update["$set"].(bson.M)["error"] = errMsg
	}

	// Update the job
	_, err := d.pdfJobColl.UpdateOne(ctx, bson.M{"job_id": jobID}, update)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}
