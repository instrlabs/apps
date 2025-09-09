package internal

import (
	"context"
	"log"
)

type Processor struct {
	mongo    *MongoDB
	s3Serv   *S3Service
	natsServ *NatsService
	imgServ  *ImageService
}

func NewProcessor(mongo *MongoDB, s3Serv *S3Service, natsServ *NatsService, imgServ *ImageService) *Processor {
	return &Processor{
		mongo:    mongo,
		s3Serv:   s3Serv,
		natsServ: natsServ,
		imgServ:  imgServ,
	}
}

func (p *Processor) Handle(ctx context.Context, job *JobMessage) error {
	log.Printf("jobID: %v, userID: %v", job.ID, job.UserID)

	// GET instruction by ID
	// EXEC instruction by key

	return nil
}
