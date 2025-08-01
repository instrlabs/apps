package handlers

import (
	"github.com/gofiber/fiber/v2"
	"labs-service/constants"
	"labs-service/models"
	"labs-service/repositories"
	"labs-service/services"
	"log"
)

type PDFJobHandler struct {
	jobRepo     *repositories.JobRepository
	pdfService  *services.PDFService
	s3Service   *services.S3Service
	natsService *services.NatsService
	cfg         *constants.Config
}

func NewPDFJobHandler(
	jobRepo *repositories.JobRepository,
	pdfService *services.PDFService,
	s3Service *services.S3Service,
	natsService *services.NatsService,
	cfg *constants.Config,
) *PDFJobHandler {
	return &PDFJobHandler{
		jobRepo:     jobRepo,
		pdfService:  pdfService,
		s3Service:   s3Service,
		natsService: natsService,
		cfg:         cfg,
	}
}

func (c *PDFJobHandler) CompressPDF(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
		})
	}

	job := &models.Job{
		OriginalFilename: file.Filename,
		JobType:          models.JobTypePDFCompress,
		Status:           models.JobStatusPending,
	}

	job, err = c.jobRepo.Create(ctx.Context(), job)
	if err != nil {
		log.Printf("Failed to create job: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create job",
		})
	}

	s3Path, err := c.s3Service.UploadPDF(ctx.Context(), file, job.ID.Hex())
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	pdfJob, err := c.pdfService.CreateJob(ctx, services.PDFCreateJobRequest{
		JobID:     job.ID.Hex(),
		Operation: models.JobTypePDFCompress,
		Filename:  job.ID.Hex() + ".pdf",
		FileSize:  file.Size,
		S3Path:    s3Path,
	})
	if err != nil {
		log.Printf("Failed to create PDF job: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create PDF job",
		})
	}

	err = c.natsService.PublishJobID(c.cfg.PDFJobsSubject, pdfJob.ID)
	if err != nil {
		log.Printf("Failed to publish job ID: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish job ID",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"job_id":     job.ID.Hex(),
		"pdf_job_id": pdfJob.ID,
		"status":     string(job.Status),
	})
}
