package controllers

import (
	"github.com/gofiber/fiber/v2"
	"labs-service/constants"
	"labs-service/models"
	"labs-service/repositories"
	"labs-service/services"
	"log"
)

type PDFJobController struct {
	jobRepo     *repositories.JobRepository
	pdfJobRepo  *repositories.PDFJobRepository
	s3Service   *services.S3Service
	natsService *services.NatsService
	cfg         *constants.Config
}

func NewPDFJobController(
	jobRepo *repositories.JobRepository,
	pdfJobRepo *repositories.PDFJobRepository,
	s3Service *services.S3Service,
	natsService *services.NatsService,
	cfg *constants.Config,
) *PDFJobController {
	return &PDFJobController{
		jobRepo:     jobRepo,
		pdfJobRepo:  pdfJobRepo,
		s3Service:   s3Service,
		natsService: natsService,
		cfg:         cfg,
	}
}

func (c *PDFJobController) ConvertToJPG(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
		})
	}

	job := &models.Job{
		Filename: file.Filename,
		JobType:  models.JobTypePDFToJPG,
		Status:   models.JobStatusPending,
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

	pdfJob := &models.PDFJob{
		OriginalName: file.Filename,
		FileSize:     file.Size,
		S3Path:       s3Path,
		Operation:    models.PDFOperationToJPG,
		JobID:        job.ID.Hex(),
	}

	pdfJob, err = c.pdfJobRepo.Create(ctx.Context(), pdfJob)
	if err != nil {
		log.Printf("Failed to create PDF job: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create PDF job",
		})
	}

	err = c.natsService.PublishJobID("PDF_OPERATION", job.ID.Hex())
	if err != nil {
		log.Printf("Failed to publish job ID: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish job ID",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"job_id": job.ID.Hex(),
		"status": string(job.Status),
	})
}

func (c *PDFJobController) CompressPDF(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
		})
	}

	job := &models.Job{
		Filename: file.Filename,
		JobType:  models.JobTypePDFCompress,
		Status:   models.JobStatusPending,
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

	pdfJob := &models.PDFJob{
		OriginalName: file.Filename,
		FileSize:     file.Size,
		S3Path:       s3Path,
		Operation:    models.PDFOperationCompress,
		JobID:        job.ID.Hex(),
	}

	pdfJob, err = c.pdfJobRepo.Create(ctx.Context(), pdfJob)
	if err != nil {
		log.Printf("Failed to create PDF job: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create PDF job",
		})
	}

	err = c.natsService.PublishJobID("PDF_OPERATION", job.ID.Hex())
	if err != nil {
		log.Printf("Failed to publish job ID: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish job ID",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"job_id": job.ID.Hex(),
		"status": string(job.Status),
	})
}

func (c *PDFJobController) MergePDFs(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	files := form.File["files"]
	if len(files) < 2 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least two PDF files are required for merging",
		})
	}

	job := &models.Job{
		Filename: "merged.pdf",
		JobType:  models.JobTypePDFMerge,
		Status:   models.JobStatusPending,
	}

	job, err = c.jobRepo.Create(ctx.Context(), job)
	if err != nil {
		log.Printf("Failed to create job: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create job",
		})
	}

	s3Path, err := c.s3Service.UploadPDF(ctx.Context(), files[0], job.ID.Hex())
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	pdfJob := &models.PDFJob{
		OriginalName: files[0].Filename,
		FileSize:     files[0].Size,
		S3Path:       s3Path,
		Operation:    models.PDFOperationMerge,
		JobID:        job.ID.Hex(),
	}

	pdfJob, err = c.pdfJobRepo.Create(ctx.Context(), pdfJob)
	if err != nil {
		log.Printf("Failed to create PDF job: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create PDF job",
		})
	}

	for i := 1; i < len(files); i++ {
		s3Path, err := c.s3Service.UploadPDF(ctx.Context(), files[i], job.ID.Hex()+"-"+string(i))
		if err != nil {
			log.Printf("Failed to upload file %d: %v", i, err)
			_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to upload file",
			})
		}

		additionalPdfJob := &models.PDFJob{
			OriginalName: files[i].Filename,
			FileSize:     files[i].Size,
			S3Path:       s3Path,
			Operation:    models.PDFOperationMerge,
			JobID:        job.ID.Hex(),
		}

		_, err = c.pdfJobRepo.Create(ctx.Context(), additionalPdfJob)
		if err != nil {
			log.Printf("Failed to create PDF job for file %d: %v", i, err)
			_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create PDF job",
			})
		}
	}

	err = c.natsService.PublishJobID("PDF_OPERATION", job.ID.Hex())
	if err != nil {
		log.Printf("Failed to publish job ID: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish job ID",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"job_id": job.ID.Hex(),
		"status": string(job.Status),
	})
}

func (c *PDFJobController) SplitPDF(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
		})
	}

	pageRanges := ctx.FormValue("page_ranges", "")
	if pageRanges == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Page ranges are required",
		})
	}

	job := &models.Job{
		Filename: file.Filename,
		JobType:  models.JobTypePDFSplit,
		Status:   models.JobStatusPending,
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

	pdfJob := &models.PDFJob{
		OriginalName: file.Filename,
		FileSize:     file.Size,
		S3Path:       s3Path,
		Operation:    models.PDFOperationSplit,
		JobID:        job.ID.Hex(),
	}

	pdfJob, err = c.pdfJobRepo.Create(ctx.Context(), pdfJob)
	if err != nil {
		log.Printf("Failed to create PDF job: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create PDF job",
		})
	}

	err = c.natsService.PublishJobID("PDF_OPERATION", job.ID.Hex())
	if err != nil {
		log.Printf("Failed to publish job ID: %v", err)
		_, _ = c.jobRepo.UpdateStatus(ctx.Context(), job.ID.Hex(), models.JobStatusFailed, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish job ID",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"job_id": job.ID.Hex(),
		"status": string(job.Status),
	})
}

func (c *PDFJobController) GetPDFJobs(ctx *fiber.Ctx) error {
	pdfJobs, err := c.pdfJobRepo.FindAll(ctx.Context())
	if err != nil {
		log.Printf("Failed to get PDF jobs: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get PDF jobs",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"jobs": pdfJobs,
	})
}
