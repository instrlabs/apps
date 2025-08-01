package repositories

import (
	"context"
	"pdf-service/models"
)

// PDFJobRepositoryInterface defines the contract for PDF job repository operations
type PDFJobRepositoryInterface interface {
	// Create inserts a new PDF job into the database
	Create(ctx context.Context, job *models.PDFJob) error

	// FindByID retrieves a PDF job by its ID
	FindByID(ctx context.Context, id string) (*models.PDFJob, error)

	// FindAll retrieves all PDF jobs with pagination
	FindAll(ctx context.Context, limit, offset int64) ([]*models.PDFJob, error)

	// Update updates a PDF job with new information
	Update(ctx context.Context, id string, update *models.UpdatePDFJobRequest) error
}
