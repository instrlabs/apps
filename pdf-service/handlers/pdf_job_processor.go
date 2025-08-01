package handlers

import (
	"context"
	"log"
	"pdf-service/models"
	"pdf-service/repositories"
	"pdf-service/services"
)

type PDFJobProcessor struct {
	jobRepo    *repositories.JobRepository
	pdfJobRepo *repositories.PDFJobRepository
	s3Service  *services.S3Service
}

func NewPDFJobProcessor(
	jobRepo *repositories.JobRepository,
	pdfJobRepo *repositories.PDFJobRepository,
	s3Service *services.S3Service,
) *PDFJobProcessor {
	return &PDFJobProcessor{
		jobRepo:    jobRepo,
		pdfJobRepo: pdfJobRepo,
		s3Service:  s3Service,
	}
}

func (p *PDFJobProcessor) ProcessJob(ctx context.Context, job *models.PDFJob) error {
	log.Printf("Processing PDF job: %s, operation: %s", job.JobID, job.Operation)

	_, err := p.jobRepo.UpdateStatus(ctx, job.JobID, models.JobStatusProcessing, "")
	if err != nil {
		log.Printf("Error updating job status: %v", err)
		return err
	}

	switch job.Operation {
	case models.PDFOperationConvertToJPG:
		log.Printf("Converting PDF to JPG: %s", job.S3Path)
		// TODO: Implement conversion logic
	case models.PDFOperationCompress:
		log.Printf("Compressing PDF: %s", job.S3Path)
		// TODO: Implement compression logic
	case models.PDFOperationMerge:
		log.Printf("Merging PDFs: %s", job.S3Path)
		// TODO: Implement merge logic
	case models.PDFOperationSplit:
		log.Printf("Splitting PDF: %s", job.S3Path)
		// TODO: Implement split logic
	default:
		log.Printf("Unknown operation: %s", job.Operation)
	}

	updatePDFJob := &models.UpdatePDFJobRequest{
		OutputFilePath: job.S3Path + ".processed",
	}

	err = p.pdfJobRepo.Update(ctx, job.ID.Hex(), updatePDFJob)
	if err != nil {
		log.Printf("Error updating job: %v", err)
		return err
	}

	log.Printf("Job processed successfully: %s", job.JobID)
	return nil
}
