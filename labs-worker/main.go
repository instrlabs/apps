package main

import (
	"context"
	"labs-worker/constants"
	π "labs-worker/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting labs worker...")

	cfg := constants.NewConfig()

	s3Service, err := services.NewS3Service(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}

	natsService, err := services.NewNatsService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize NATS service: %v", err)
	}
	defer natsService.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	err = natsService.SubscribeToPDFJobs(func(jobID string) {
		processPDFJob(ctx, jobID, s3Service, natsService)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to pdf.jobs: %v", err)
	}

	log.Println("Worker is running. Press Ctrl+C to exit.")

	<-ctx.Done()
	log.Println("Shutting down...")
}

func processPDFJob(ctx context.Context, jobID string, s3Service *services.S3Service, natsService *services.NatsService) {
	log.Printf("Processing job: %s", jobID)

	log.Printf("Job completed successfully: %s", jobID)
}
