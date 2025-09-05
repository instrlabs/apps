package internal

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructionHandler struct {
	s3          *S3Service
	fileRepo    *FileRepository
	instrRepo   *InstructionRepository
	productServ *ProductService
}

func NewInstructionHandler(s3 *S3Service, fileRepo *FileRepository, instrRepo *InstructionRepository, productServ *ProductService) *InstructionHandler {
	return &InstructionHandler{s3: s3, fileRepo: fileRepo, instrRepo: instrRepo, productServ: productServ}
}

func (h *InstructionHandler) ImageCompress(c *fiber.Ctx) error {
	product := h.productServ.GetProductByKey("image-compress")
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "product not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	localUserID, _ := c.Locals("UserID").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	productID, _ := primitive.ObjectIDFromHex(product.ID)
	instructionID := primitive.NewObjectID()

	_ = h.instrRepo.Create(&Instruction{
		ID:        instructionID,
		UserID:    userID,
		ProductID: productID,
		Status:    InstructionStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	var headers []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil && form != nil && form.File != nil {
		if fs, ok := form.File["files"]; ok && len(fs) > 0 {
			headers = fs
		} else if fs, ok := form.File["file"]; ok && len(fs) > 0 {
			headers = fs
		}
	}

	for idx, fh := range headers {
		if err := h.s3.Upload(fh, fmt.Sprintf("%s-%d", instructionID, idx)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "upload failed",
				"errors":  nil,
				"data":    nil,
			})
		}

		if h.fileRepo != nil {
			_, _ = h.fileRepo.Create(&File{
				InstructionID: instructionID,
				Type:          FileTypeRequest,
				CreatedAt:     time.Now().UTC(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "upload accepted",
		"errors":  nil,
		"data":    fiber.Map{"instruction_id": instructionID},
	})
}
