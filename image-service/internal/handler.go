package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	initx "github.com/histweety-labs/shared/init"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructionHandler struct {
	cfg         *Config
	s3          *initx.S3
	nats        *initx.Nats
	instrRepo   *InstructionRepository
	productServ *ProductService
}

func NewInstructionHandler(
	cfg *Config,
	s3 *initx.S3,
	nats *initx.Nats,
	instrRepo *InstructionRepository,
	productServ *ProductService) *InstructionHandler {
	return &InstructionHandler{cfg: cfg, s3: s3, nats: nats, instrRepo: instrRepo, productServ: productServ}
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

	// We'll upload files to S3 first, then create the instruction once uploads succeed
	instr := &Instruction{
		ID:        instructionID,
		UserID:    userID,
		ProductID: productID,
		Status:    InstructionStatusPending,
		Inputs:    []File{},
		Outputs:   []File{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

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
		f, _ := fh.Open()
		b, _ := io.ReadAll(f)
		f.Close()

		ext := filepath.Ext(fh.Filename)
		objectName := fmt.Sprintf("images/%s-%d%s", instructionID.Hex(), idx, ext)
		if err := h.s3.Put(objectName, b); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "upload failed",
				"errors":  nil,
				"data":    nil,
			})
		}

		instr.Inputs = append(instr.Inputs, File{
			FileName: objectName,
			Type:     "request",
			Size:     int64(len(b)),
		})
	}
	// After successful uploads, persist the instruction with collected inputs
	if err := h.instrRepo.Create(instr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create instruction",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Publish job request
	payload := map[string]string{
		"id":     instructionID.Hex(),
		"userId": userID.Hex(),
	}
	if data, err := json.Marshal(payload); err == nil && h.nats != nil && h.nats.Conn != nil {
		_ = h.nats.Conn.Publish(h.cfg.NatsSubjectRequests, data)
	}

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
	if err := h.instrRepo.UpdateStatus(id, body.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update status", "errors": err.Error(), "data": nil})
	}
	instr := h.instrRepo.GetByID(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "updated", "errors": nil, "data": instr})
}
