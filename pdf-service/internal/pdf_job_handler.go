package internal

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PDFJobHandler struct {
	pdfJobRepo  *PDFJobRepository
	s3Service   *S3Service
	natsService *NatsService
}

func NewPDFJobHandler(
	pdfJobRepo *PDFJobRepository,
	s3Service *S3Service,
	natsService *NatsService,
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
	var request struct {
		JobID     string `json:"job_id"`
		Operation string `json:"operation"`
		Filename  string `json:"filename"`
		FileSize  int64  `json:"file_size"`
		S3Path    string `json:"s3_path"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.JobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JobID is required",
		})
	}

	if request.Operation == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Operation is required",
		})
	}

	if request.S3Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "S3Path is required",
		})
	}

	pdfOperation := PDFOperation(request.Operation)
	if pdfOperation != PDFOperationConvertToJPG &&
		pdfOperation != PDFOperationCompress &&
		pdfOperation != PDFOperationMerge &&
		pdfOperation != PDFOperationSplit {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid operation",
		})
	}

	job := &PDFJob{
		Filename:  request.Filename,
		FileSize:  request.FileSize,
		S3Path:    request.S3Path,
		Operation: pdfOperation,
		JobID:     request.JobID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := h.pdfJobRepo.Create(context.Background(), job)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create job",
		})
	}

	err = h.natsService.PublishPDFJob(job.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish job to NATS",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": job,
	})
}

func (h *PDFJobHandler) GetFile(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	job, err := h.pdfJobRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "PDF job not found",
		})
	}

	// Determine which file to serve (original or processed)
	s3Path := job.S3Path
	if job.OutputFilePath != "" {
		s3Path = job.OutputFilePath
	}

	// Download the file from S3
	tempFilePath, err := h.s3Service.DownloadPDF(c.Context(), s3Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve file from storage",
		})
	}
	defer os.Remove(tempFilePath) // Clean up the temp file after serving

	// No need to determine content type as c.Download handles it automatically

	// Set the filename for the download
	filename := job.Filename
	if filepath.Ext(s3Path) != filepath.Ext(filename) {
		filename = filepath.Base(s3Path)
	}

	return c.Download(tempFilePath, filename)
}
