package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"

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

func (p *PDFJobProcessor) compressPDF(ctx context.Context, pdfPath string) (string, error) {
	outputPath := pdfPath + ".compressed.pdf"

	err := api.OptimizeFile(pdfPath, outputPath, nil)

	if err != nil {
		return "", fmt.Errorf("failed to compress PDF: %w", err)
	}

	return outputPath, nil
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
		return err
	}

	localPDFPath, err := p.s3Service.DownloadPDF(ctx, pdfJob.S3Path)
	if err != nil {
		log.Printf("Error downloading PDF: %v", err)
		p.natsService.PublishJobNotification(pdfJob.JobID, models.JobStatusFailed)
		return err
	}
	defer os.Remove(localPDFPath)

	var outputPath string
	var contentType string

	switch pdfJob.Operation {
	case models.PDFOperationCompress:
		log.Printf("Compressing PDF: %s", job.ID)
		outputPath, err = p.compressPDF(ctx, localPDFPath)
		contentType = "application/pdf"
	default:
		log.Printf("Unknown operation: %s", pdfJob.Operation)
		p.natsService.PublishJobNotification(pdfJob.JobID, models.JobStatusFailed)
		return fmt.Errorf("unknown operation: %s", pdfJob.Operation)
	}

	if err != nil {
		log.Printf("Error processing PDF: %v", err)
		p.natsService.PublishJobNotification(pdfJob.JobID, models.JobStatusFailed)
		return err
	}

	// Ensure the output file is cleaned up
	defer os.Remove(outputPath)

	s3OutputPath := fmt.Sprintf("processed/%s/%s", job.ID, filepath.Base(outputPath))
	s3OutputPath, err = p.s3Service.UploadProcessedFile(ctx, outputPath, s3OutputPath, contentType)
	if err != nil {
		log.Printf("Error uploading processed file: %v", err)
		p.natsService.PublishJobNotification(pdfJob.JobID, models.JobStatusFailed)
		return err
	}

	updateRequest := &models.UpdatePDFJobRequest{
		OutputFilePath: s3OutputPath,
	}

	err = p.pdfJobRepo.Update(ctx, pdfJob.ID.Hex(), updateRequest)
	if err != nil {
		log.Printf("Error updating job: %v", err)
		p.natsService.PublishJobNotification(pdfJob.JobID, models.JobStatusFailed)
		return err
	}

	err = p.natsService.PublishJobNotification(pdfJob.JobID, models.JobStatusCompleted)
	if err != nil {
		log.Printf("Error publishing job completion notification: %v", err)
		return err
	}

	log.Printf("Job processed successfully: %s", job.ID)
	return nil
}
