package handlers

import (
	"context"
	"labs-service/models"
	"labs-service/repositories"
	"labs-service/services"
	"log"
)

type JobNotificationProcessor struct {
	jobRepo   *repositories.JobRepository
	s3Service *services.S3Service
}

func NewPDFNotificationProcessor(
	jobRepo *repositories.JobRepository,
	s3Service *services.S3Service,
) *JobNotificationProcessor {
	return &JobNotificationProcessor{
		jobRepo:   jobRepo,
		s3Service: s3Service,
	}
}

func (p *JobNotificationProcessor) ProcessJob(ctx context.Context, job *models.JobNotificationMessage) error {
	log.Printf("Processing job: %s, status: %s", job.ID, job.Status)

	_, err := p.jobRepo.UpdateStatus(ctx, job.ID, job.Status, "")
	if err != nil {
		log.Printf("Error updating job status: %v", err)
		return err
	}

	log.Printf("Job processed successfully: %s", job.ID)
	return nil
}
