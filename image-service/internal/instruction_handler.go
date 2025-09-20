package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os/exec"
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

func (h *InstructionHandler) CreateInstruction(c *fiber.Ctx) error {
	productId := c.Params("product_id")
	localUserID, _ := c.Locals("UserID").(string)
	product := h.productServ.GetProduct(localUserID, productId)
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "product not found",
			"errors":  nil,
			"data":    nil,
		})
	}

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
		_ = f.Close()

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

	h.publishNotification(instructionID.Hex(), InstructionStatusPending)

	if data, err := json.Marshal(&InstructionRequest{
		UserID:        userID.Hex(),
		InstructionID: instructionID.Hex(),
	}); err == nil {
		_ = h.nats.Conn.Publish(h.cfg.NatsSubjectImagesRequests, data)
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

func (h *InstructionHandler) RunInstruction(c *fiber.Ctx) error {
	var req InstructionRequest
	if c != nil && len(c.Body()) > 0 {
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			log.Printf("RunInstruction: invalid body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid body", "errors": nil, "data": nil})
		}
		// process
		h.RunInstructionMessage(c.Body())
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "instruction running", "errors": nil, "data": fiber.Map{"instruction_id": req.InstructionID}})
	}
	log.Println("RunInstruction invoked without context body")
	return nil
}

func (h *InstructionHandler) publishNotification(instructionID string, status InstructionStatus) {
	if h == nil || h.nats == nil || h.nats.Conn == nil {
		return
	}
	data, err := json.Marshal(&InstructionNotification{
		InstructionID:     instructionID,
		InstructionStatus: string(status),
	})
	if err == nil {
		_ = h.nats.Conn.Publish(h.cfg.NatsSubjectNotificationsSSE, data)
	}
}

func (h *InstructionHandler) processImageCompression(instr *Instruction) {
	var outputs []File
	for idx, input := range instr.Inputs {
		inputBytes := h.s3.Get("images/" + input.FileName)
		if inputBytes == nil {
			log.Printf("processImageCompression: input file not found: %s", input.FileName)
			continue
		}

		ext := filepath.Ext(input.FileName)
		outputName := fmt.Sprintf("%s-out-%d%s", instr.ID.Hex(), idx, ext)

		cmd := exec.Command("convert", "-", "-quality", "60", "-")
		cmd.Stdin = bytes.NewReader(inputBytes)
		outputBytes, err := cmd.Output()
		if err != nil {
			log.Printf("compression failed: %v", err)
			continue
		}

		if err := h.s3.Put("images/"+outputName, outputBytes); err != nil {
			log.Printf("failed to upload compressed file: %v", err)
			continue
		}

		outputs = append(outputs, File{
			FileName: outputName,
			Size:     int64(len(outputBytes)),
		})
	}

	if err := h.instrRepo.UpdateOutputs(instr.ID, outputs); err != nil {
		log.Printf("failed to update outputs: %v", err)
		_ = h.instrRepo.UpdateStatus(instr.ID, InstructionStatusFailed)
		h.publishNotification(instr.ID.Hex(), InstructionStatusFailed)
		return
	}

	_ = h.instrRepo.UpdateStatus(instr.ID, InstructionStatusCompleted)
	h.publishNotification(instr.ID.Hex(), InstructionStatusCompleted)
}

func (h *InstructionHandler) RunInstructionMessage(data []byte) {
	var msg InstructionRequest
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("RunInstructionMessage: unmarshal error: %v", err)
		return
	}

	id, err := primitive.ObjectIDFromHex(msg.InstructionID)
	if err != nil {
		log.Printf("RunInstructionMessage: invalid instruction id: %v", err)
		return
	}

	instr := h.instrRepo.GetByID(id)
	if instr == nil || instr.ID.IsZero() {
		log.Printf("RunInstructionMessage: instruction not found")
		return
	}

	product := h.productServ.GetProduct(instr.UserID.Hex(), instr.ProductID.Hex())
	if product == nil {
		log.Printf("RunInstructionMessage: product not found")
	}

	_ = h.instrRepo.UpdateStatus(instr.ID, InstructionStatusProcessing)
	h.publishNotification(instr.ID.Hex(), InstructionStatusProcessing)

	switch product.Key {
	case "image-compress":
		h.processImageCompression(instr)
	}
}

func (h *InstructionHandler) CleanInstruction(c *fiber.Ctx) error {
	log.Printf("CleanInstruction invoked")
	return nil
}
