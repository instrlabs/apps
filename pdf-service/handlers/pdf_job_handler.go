package handlers

import (
	"context"
	"pdf-service/models"
	"pdf-service/repositories"
	"pdf-service/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PDFJobHandler struct {
	pdfJobRepo  repositories.PDFJobRepositoryInterface
	s3Service   *services.S3Service
	natsService *services.NatsService
}

func NewPDFJobHandler(
	pdfJobRepo repositories.PDFJobRepositoryInterface,
	s3Service *services.S3Service,
	natsService *services.NatsService,
) *PDFJobHandler {
	return &PDFJobHandler{
		pdfJobRepo:  pdfJobRepo,
		s3Service:   s3Service,
		natsService: natsService,
	}
}

func (h *PDFJobHandler) GetJobs(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit parameter",
		})
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid offset parameter",
		})
	}

	jobs, err := h.pdfJobRepo.FindAll(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve jobs",
		})
	}

	return c.JSON(fiber.Map{
		"data": jobs,
	})
}

func (h *PDFJobHandler) CreateJob(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Expected multipart form",
		})
	}

	operation := c.FormValue("operation")
	if operation == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Operation is required",
		})
	}

	pdfOperation := models.PDFOperation(operation)
	switch pdfOperation {
	case models.PDFOperationConvertToJPG, models.PDFOperationCompress, models.PDFOperationMerge, models.PDFOperationSplit:
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid operation",
		})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	file := files[0]

	jobID := uuid.New().String()

	s3Path, err := h.s3Service.UploadPDF(c.Context(), file, jobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	job := &models.PDFJob{
		OriginalName: file.Filename,
		FileSize:     file.Size,
		S3Path:       s3Path,
		Operation:    pdfOperation,
		JobID:        jobID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = h.pdfJobRepo.Create(context.Background(), job)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create job",
		})
	}

	err = h.natsService.PublishPDFJob(job)
	if err != nil {
		// Log error but don't fail the request
		// In a production environment, you might want to implement a retry mechanism
		// or mark the job as failed in the database
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": job,
	})
}
