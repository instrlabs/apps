package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strconv"
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
	productServ *PaymentService
}

func NewInstructionHandler(
	cfg *Config,
	s3 *initx.S3,
	nats *initx.Nats,
	instrRepo *InstructionRepository,
	productServ *PaymentService) *InstructionHandler {
	return &InstructionHandler{cfg: cfg, s3: s3, nats: nats, instrRepo: instrRepo, productServ: productServ}
}

func (h *InstructionHandler) ImageCompress(c *fiber.Ctx) error {
	product := h.productServ.GetProduct(c, "image-compress")
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

	if len(headers) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "no file uploaded",
			"errors":  nil,
			"data":    nil,
		})
	}

	for idx, fh := range headers {
		f, _ := fh.Open()
		b, _ := io.ReadAll(f)
		f.Close()

		ext := filepath.Ext(fh.Filename)
		objectName := fmt.Sprintf("%s-%d%s", instructionID.Hex(), idx, ext)
		if err := h.s3.Put("images/"+objectName, b); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "upload failed",
				"errors":  nil,
				"data":    nil,
			})
		}

		instr.Inputs = append(instr.Inputs, File{
			FileName: objectName,
			Size:     int64(len(b)),
		})
	}
	if err := h.instrRepo.Create(instr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create instruction",
			"errors":  nil,
			"data":    nil,
		})
	}

	if data, err := json.Marshal(&InstructionRequest{
		UserID:        userID.Hex(),
		InstructionID: instructionID.Hex(),
	}); err == nil {
		_ = h.nats.Conn.Publish(h.cfg.NatsSubjectRequests, data)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "upload accepted",
		"errors":  nil,
		"data":    fiber.Map{"instruction_id": instructionID},
	})
}

func (h *InstructionHandler) ListInstructions(c *fiber.Ctx) error {
	localUserID, _ := c.Locals("UserID").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	list := h.instrRepo.ListByUser(userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok", "errors": nil, "data": list})
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

func (h *InstructionHandler) UpdateInstructionOutputs(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id", "errors": nil, "data": nil})
	}
	var body struct {
		Outputs []File `json:"outputs"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid body", "errors": nil, "data": nil})
	}
	for _, f := range body.Outputs {
		if f.FileName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "file_name required", "errors": nil, "data": nil})
		}
	}
	if err := h.instrRepo.UpdateOutputs(id, body.Outputs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to update outputs", "errors": err.Error(), "data": nil})
	}
	instr := h.instrRepo.GetByID(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "outputs updated", "errors": nil, "data": instr})
}

// GetInstructionFile streams a specific file for an instruction if owned by the user
func (h *InstructionHandler) GetInstructionFile(c *fiber.Ctx) error {
	idHex := c.Params("id")
	fileName := c.Params("file_name")
	if idHex == "" || fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id and file_name required", "errors": nil, "data": nil})
	}

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id", "errors": nil, "data": nil})
	}

	instr := h.instrRepo.GetByID(id)
	if instr == nil || instr.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction not found", "errors": nil, "data": nil})
	}

	localUserID, _ := c.Locals("UserID").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden", "errors": nil, "data": nil})
	}

	found := false
	for _, f := range instr.Inputs {
		if f.FileName == fileName {
			found = true
			break
		}
	}
	if !found {
		for _, f := range instr.Outputs {
			if f.FileName == fileName {
				found = true
				break
			}
		}
	}
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "file not found in instruction", "errors": nil, "data": nil})
	}

	// Retrieve from S3
	b := h.s3.Get("images/" + fileName)
	if b == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "file not found", "errors": nil, "data": nil})
	}

	ext := filepath.Ext(fileName)
	ct := mime.TypeByExtension(ext)
	c.Set("Content-Type", ct)
	c.Set("Content-Length", strconv.FormatInt(int64(len(b)), 10))
	c.Attachment(fileName)
	return c.Status(fiber.StatusOK).Send(b)
}

// GetInstructionFiles returns all input and output files for a given instruction if owned by the user
func (h *InstructionHandler) GetInstructionFiles(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id", "errors": nil, "data": nil})
	}

	instr := h.instrRepo.GetByID(id)
	if instr == nil || instr.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction not found", "errors": nil, "data": nil})
	}

	localUserID, _ := c.Locals("UserID").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden", "errors": nil, "data": nil})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data": fiber.Map{
			"inputs":  instr.Inputs,
			"outputs": instr.Outputs,
		},
	})
}
