package handlers

import (
	"context"
	"log"
	"pdf-service/models"
	"pdf-service/repositories"
	"pdf-service/services"
)

type PDFJobProcessor struct {
	jobRepo     *repositories.JobRepository
	pdfJobRepo  *repositories.PDFJobRepository
	s3Service   *services.S3Service
	natsService *services.NatsService
}

func NewPDFJobProcessor(
	jobRepo *repositories.JobRepository,
	pdfJobRepo *repositories.PDFJobRepository,
	s3Service *services.S3Service,
	natsService *services.NatsService,
) *PDFJobProcessor {
	return &PDFJobProcessor{
		jobRepo:     jobRepo,
		pdfJobRepo:  pdfJobRepo,
		s3Service:   s3Service,
		natsService: natsService,
	}
}

func (p *PDFJobProcessor) ProcessJob(ctx context.Context, job *models.PDFJobMessage) error {
	log.Printf("Processing PDF job: %s", job.ID)

	pdfJob, err := p.pdfJobRepo.FindByID(ctx, job.ID)
	if err != nil {
		log.Printf("Error finding job: %v", err)
		return err
	}

	err = p.natsService.PublishJobNotification(pdfJob.JobID, models.JobStatusProcessing)
	if err != nil {
		log.Printf("Error publishing job notification: %v", err)
		return err
	}

	switch pdfJob.Operation {
	case models.PDFOperationConvertToJPG:
		log.Printf("Converting PDF to JPG: %s", job.ID)
		// TODO: Implement conversion logic
	case models.PDFOperationCompress:
		log.Printf("Compressing PDF: %s", job.ID)
		// TODO: Implement compression logic
	case models.PDFOperationMerge:
		log.Printf("Merging PDFs: %s", job.ID)
		// TODO: Implement merge logic
	case models.PDFOperationSplit:
		log.Printf("Splitting PDF: %s", job.ID)
		// TODO: Implement split logic
	default:
		log.Printf("Unknown operation: %s", pdfJob.Operation)
	}

	log.Printf("Job processed successfully: %s", job.ID)
	return nil
}
