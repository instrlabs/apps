package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"labs-service/models"
	"labs-service/services"
)

type PDFJobRepository struct {
	db         *services.MongoDB
	collection *mongo.Collection
}

func NewPDFJobRepository(db *services.MongoDB) *PDFJobRepository {
	return &PDFJobRepository{
		db:         db,
		collection: db.Collection("pdf_jobs"),
	}
}

func (r *PDFJobRepository) Create(ctx context.Context, pdf *models.PDFJob) (*models.PDFJob, error) {
	pdf.CreatedAt = time.Now()
	pdf.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, pdf)
	if err != nil {
		return nil, err
	}

	pdf.ID = result.InsertedID.(primitive.ObjectID)
	return pdf, nil
}

func (r *PDFJobRepository) FindByID(ctx context.Context, id string) (*models.PDFJob, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pdf models.PDFJob
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&pdf)
	if err != nil {
		return nil, err
	}

	return &pdf, nil
}

func (r *PDFJobRepository) FindByJobID(ctx context.Context, jobID string) (*models.PDFJob, error) {
	var pdf models.PDFJob
	err := r.collection.FindOne(ctx, bson.M{"job_id": jobID}).Decode(&pdf)
	if err != nil {
		return nil, err
	}

	return &pdf, nil
}

func (r *PDFJobRepository) FindAll(ctx context.Context) ([]models.PDFJob, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pdfs []models.PDFJob
	if err := cursor.All(ctx, &pdfs); err != nil {
		return nil, err
	}

	return pdfs, nil
}

func (r *PDFJobRepository) Update(ctx context.Context, id string, update *models.UpdatePDFJobRequest) (*models.PDFJob, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if update.OutputFilePath != "" {
		updateDoc["$set"].(bson.M)["output_file_path"] = update.OutputFilePath
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedPDF models.PDFJob
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, updateDoc, opts).Decode(&updatedPDF)
	if err != nil {
		return nil, err
	}

	return &updatedPDF, nil
}

func (r *PDFJobRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
