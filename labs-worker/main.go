package main

import (
	"context"
	"labs-worker/constants"
	"labs-worker/models"
	"labs-worker/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting PDF to JPG worker...")

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

	dbService, err := services.NewDBService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DB service: %v", err)
	}
	defer dbService.Close()

	muPDFService := services.NewMuPDFService()
	defer muPDFService.Close()

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

	err = natsService.SubscribeToJobs(func(jobID string) {
		processJob(ctx, jobID, dbService, s3Service, muPDFService, natsService)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to jobs: %v", err)
	}

	log.Println("Worker is running. Press Ctrl+C to exit.")

	<-ctx.Done()
	log.Println("Shutting down...")
}

func processJob(ctx context.Context, jobID string, dbService *services.DBService, s3Service *services.S3Service, muPDFService *services.MuPDFService, natsService *services.NatsService) {
	log.Printf("Processing job: %s", jobID)

	job, err := dbService.GetPDFJobByID(ctx, jobID)
	if err != nil {
		log.Printf("Failed to get job: %v", err)
		return
	}

	if job.Operation != models.PDFOperationToJPG {
		log.Printf("Job is not a PDF to JPG conversion: %s", job.Operation)
		return
	}

	err = dbService.UpdateJobStatus(ctx, jobID, models.JobStatusProcessing, "", "")
	if err != nil {
		log.Printf("Failed to update job status: %v", err)
		return
	}

	pdfPath, err := s3Service.DownloadPDF(ctx, job.S3Path)
	if err != nil {
		log.Printf("Failed to download PDF: %v", err)
		errMsg := err.Error()
		dbService.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, "", errMsg)
		natsService.PublishCompletion(jobID, string(models.JobStatusFailed), err)
		return
	}
	defer os.Remove(pdfPath) // Clean up the temporary file

	jpgPath, err := muPDFService.ConvertPDFToJPG(pdfPath)
	if err != nil {
		log.Printf("Failed to convert PDF to JPG: %v", err)
		errMsg := err.Error()
		dbService.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, "", errMsg)
		natsService.PublishCompletion(jobID, string(models.JobStatusFailed), err)
		return
	}
	defer os.Remove(jpgPath) // Clean up the temporary file

	outputPath, err := s3Service.UploadJPG(ctx, jpgPath, jobID)
	if err != nil {
		log.Printf("Failed to upload JPG: %v", err)
		errMsg := err.Error()
		dbService.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, "", errMsg)
		natsService.PublishCompletion(jobID, string(models.JobStatusFailed), err)
		return
	}

	err = dbService.UpdateJobStatus(ctx, jobID, models.JobStatusCompleted, outputPath, "")
	if err != nil {
		log.Printf("Failed to update job status: %v", err)
		natsService.PublishCompletion(jobID, string(models.JobStatusFailed), err)
		return
	}

	err = natsService.PublishCompletion(jobID, string(models.JobStatusCompleted), nil)
	if err != nil {
		log.Printf("Failed to publish completion event: %v", err)
		return
	}

	log.Printf("Job completed successfully: %s", jobID)
}
