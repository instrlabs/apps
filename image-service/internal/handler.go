package internal

import (
	"context"
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
	nats        *NatsService
}

func NewInstructionHandler(
	s3 *S3Service,
	fileRepo *FileRepository,
	instrRepo *InstructionRepository,
	productServ *ProductService,
	nats *NatsService) *InstructionHandler {
	return &InstructionHandler{s3: s3, fileRepo: fileRepo, instrRepo: instrRepo, productServ: productServ, nats: nats}
}

func (h *InstructionHandler) ImageCompress(c *fiber.Ctx) error {
	product := h.productServ.GetProduct("image-compress")
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
		if err := h.s3.Upload(fh, fmt.Sprintf("%s-%d", instructionID.Hex(), idx)); err != nil {
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

	_ = h.nats.PublishJobRequest(instructionID.Hex(), userID.Hex())

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "upload accepted",
		"errors":  nil,
		"data":    fiber.Map{"instruction_id": instructionID},
	})
}

func (h *InstructionHandler) GetInstructionByID(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id", "errors": nil, "data": nil})
	}
	instr := h.instrRepo.GetByID(id)
	if instr == nil || instr.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction not found", "errors": nil, "data": nil})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok", "errors": nil, "data": instr})
}

func (h *InstructionHandler) UpdateInstructionStatus(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id", "errors": nil, "data": nil})
	}
	var body struct {
		Status InstructionStatus `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid body", "errors": nil, "data": nil})
	}
	switch body.Status {
	case InstructionStatusPending, InstructionStatusProcessing, InstructionStatusCompleted, InstructionStatusFailed:
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid status", "errors": nil, "data": nil})
	}
	if err := h.instrRepo.UpdateStatus(context.Background(), id, body.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update status", "errors": err.Error(), "data": nil})
	}
	instr := h.instrRepo.GetByID(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "updated", "errors": nil, "data": instr})
}
