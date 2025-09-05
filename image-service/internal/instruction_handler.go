package internal

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructionHandler struct {
	s3       *S3Service
	fileRepo *FileRepository
}

func NewInstructionHandler(s3 *S3Service, fileRepo *FileRepository) *InstructionHandler {
	return &InstructionHandler{s3: s3, fileRepo: fileRepo}
}

func (h *InstructionHandler) ImageCompress(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "file is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	jobID := primitive.NewObjectID()

	err = h.s3.Upload(fileHeader, jobID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "upload failed",
			"errors":  nil,
			"data":    nil,
		})
	}

	f := &File{
		JobID:     jobID,
		Type:      FileTypeRequest,
		CreatedAt: time.Now().UTC(),
	}
	if h.fileRepo != nil {
		_, _ = h.fileRepo.Create(c.Context(), f)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "upload accepted",
		"errors":  nil,
		"data": fiber.Map{
			"job_id": jobID.Hex(),
		},
	})
}
