package internal

import (
	"context"
	"log"
)

type JobNotificationProcessor struct {
	jobRepo   *JobRepository
	s3Service *S3Service
}

func NewPDFNotificationProcessor(
	jobRepo *JobRepository,
	s3Service *S3Service,
) *JobNotificationProcessor {
	return &JobNotificationProcessor{
		jobRepo:   jobRepo,
		s3Service: s3Service,
	}
}

func (p *JobNotificationProcessor) ProcessJob(ctx context.Context, job *JobNotificationMessage) error {
	log.Printf("Processing job: %s, status: %s", job.ID, job.Status)

	_, err := p.jobRepo.UpdateStatus(ctx, job.ID, job.Status, "")
	if err != nil {
		log.Printf("Error updating job status: %v", err)
		return err
	}

	log.Printf("Job processed successfully: %s", job.ID)
	return nil
}
