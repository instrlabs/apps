package internal

import (
	"context"
	"log"
)

type Processor struct {
	mongo *MongoDB
	s3    *S3Service
	nats  *NatsService
}

func NewProcessor(mongo *MongoDB, s3 *S3Service, nats *NatsService) *Processor {
	return &Processor{
		mongo: mongo,
		s3:    s3,
		nats:  nats,
	}
}

func (p *Processor) Handle(ctx context.Context, job *JobMessage) error {
	log.Printf("jobID: %v, userID: %v", job.ID, job.UserID)

	return nil
}
