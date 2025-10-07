package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructionHandler struct {
	cfg         *Config
	s3          *initx.S3
	nats        *initx.Nats
	instrRepo   *InstructionRepository
	fileRepo    *FileRepository
	productRepo *ProductRepository
	imageSvc    *ImageService
}

func NewInstructionHandler(
	cfg *Config,
	s3 *initx.S3,
	nats *initx.Nats,
	instrRepo *InstructionRepository,
	fileRepo *FileRepository,
	productRepo *ProductRepository,
	imageSvc *ImageService) *InstructionHandler {
	return &InstructionHandler{cfg: cfg, s3: s3, nats: nats, instrRepo: instrRepo, fileRepo: fileRepo, productRepo: productRepo, imageSvc: imageSvc}
}

func (h *InstructionHandler) CreateInstruction(c *fiber.Ctx) error {
	type payload struct {
		ProductKey string `json:"productKey"`
	}
	var body payload
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"errors":  nil,
			"data":    nil,
		})
	}
	if body.ProductKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "productKey is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	userID, _ := c.Locals("UserID").(string)
	product, _ := h.productRepo.FindByKey(body.ProductKey)
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "product not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	instructionID := primitive.NewObjectID()
	objUserID, _ := primitive.ObjectIDFromHex(userID)
	objProductID := product.ID

	instr := &Instruction{
		ID:        instructionID,
		UserID:    objUserID,
		ProductID: objProductID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := h.instrRepo.Create(instr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create instruction",
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "instruction created",
		"errors":  nil,
		"data":    fiber.Map{"instruction": instr},
	})
}

func (h *InstructionHandler) ListInstructions(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data": map[string]interface{}{
			"instructions": []Instruction{},
		},
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok", "errors": nil, "data": map[string]interface{}{"instruction": instr}})
}

func (h *InstructionHandler) GetInstructionDetails(c *fiber.Ctx) error {
	idHex := c.Params("id")
	if idHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id required", "errors": nil, "data": nil})
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

	files := h.fileRepo.ListByInstruction(instr.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data":    fiber.Map{"files": files},
	})
}

func (h *InstructionHandler) CreateInstructionDetails(c *fiber.Ctx) error {
	instrIDHex := c.Params("id", "")
	if instrIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "instruction id required", "errors": nil, "data": nil})
	}
	instrID, err := primitive.ObjectIDFromHex(instrIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid instruction id", "errors": nil, "data": nil})
	}

	instr := h.instrRepo.GetByID(instrID)
	if instr == nil || instr.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction not found", "errors": nil, "data": nil})
	}

	localUserID, _ := c.Locals("UserID").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden", "errors": nil, "data": nil})
	}

	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "no file uploaded", "errors": nil, "data": nil})
	}

	f, _ := fh.Open()
	b, _ := io.ReadAll(f)
	_ = f.Close()

	inputID := primitive.NewObjectID()
	outputID := primitive.NewObjectID()
	ext := filepath.Ext(fh.Filename)

	inName := "images/" + inputID.Hex() + ext
	outName := "images/" + outputID.Hex() + ext

	input := &File{
		ID:            inputID,
		InstructionID: instr.ID,
		OriginalName:  fh.Filename,
		FileName:      inName,
		Size:          int64(len(b)),
		Status:        FileStatusPending,
		OutputID:      &outputID,
	}

	output := &File{
		ID:            outputID,
		InstructionID: instr.ID,
		OriginalName:  filepath.Base(outName),
		FileName:      outName,
		Size:          0,
		Status:        FileStatusPending,
	}

	if err := h.fileRepo.CreateMany([]*File{input, output}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create file records",
			"errors":  nil,
			"data":    nil,
		})
	}

	if err := h.s3.Put(input.FileName, b); err != nil {
		_ = h.fileRepo.UpdateStatus(inputID, FileStatusFailed)
		_ = h.fileRepo.UpdateStatus(outputID, FileStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "upload failed",
			"errors":  nil,
			"data":    nil,
		})
	}

	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectImageRequests, []byte(inputID.Hex())); err != nil {
		_ = h.fileRepo.UpdateStatus(inputID, FileStatusFailed)
		_ = h.fileRepo.UpdateStatus(outputID, FileStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to publish to nats", "errors": nil, "data": nil})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file created",
		"errors":  nil,
		"data": fiber.Map{
			"input":  input,
			"output": output,
		},
	})
}

func (h *InstructionHandler) RunInstructionMessage(data []byte) {
	fileIDHex := string(bytes.TrimSpace(data))
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		log.Infof("RunInstructionMessage: invalid file id: %q err=%v", fileIDHex, err)
		return
	}

	// 1. Find file by ID
	input := h.fileRepo.GetByID(fileID)
	if input == nil || input.ID.IsZero() {
		log.Infof("RunInstructionMessage: input file not found: %s", fileIDHex)
		return
	}

	output := h.fileRepo.GetByID(*input.OutputID)
	if output == nil || output.ID.IsZero() {
		log.Infof("RunInstructionMessage: output file not found: %s", fileIDHex)
		return
	}

	// 2. Find instruction
	instr := h.instrRepo.GetByID(input.InstructionID)
	if instr == nil || instr.ID.IsZero() {
		log.Infof("RunInstructionMessage: instruction not found: %s", fileIDHex)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.fileRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(primitive.NilObjectID, input.InstructionID)
		return
	}

	// 3. Find product by instruction's product ID
	product, _ := h.productRepo.FindByID(instr.ProductID)
	if product == nil {
		log.Infof("RunInstructionMessage: product not found: %s", instr.ProductID.Hex())
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.fileRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID)
		return
	}

	// 4. Get binary from S3
	_ = h.fileRepo.UpdateStatus(input.ID, FileStatusProcessing)
	h.publishFileNotification(instr.UserID, instr.ID)
	inputBytes := h.s3.Get(input.FileName)
	if inputBytes == nil {
		log.Infof("RunInstructionMessage: input file missing on S3: %s", input.FileName)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.fileRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID)
		return
	}

	// 5. Process based on product key
	var outputBytes []byte
	switch product.Key {
	case "image-compress":
		outputBytes, err = h.imageSvc.Compress(inputBytes)
		if err != nil {
			log.Infof("RunInstructionMessage: image-compress failed: %v", err)
			_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
			_ = h.fileRepo.UpdateStatus(output.ID, FileStatusFailed)
			h.publishFileNotification(instr.UserID, instr.ID)
			return
		}
	default:
		log.Infof("RunInstructionMessage: unsupported product key: %s", product.Key)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.fileRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID)
		return
	}

	_ = h.fileRepo.UpdateStatus(input.ID, FileStatusDone)
	_ = h.fileRepo.UpdateStatus(output.ID, FileStatusProcessing)
	h.publishFileNotification(instr.UserID, instr.ID)

	// 6. Upload output to S3
	if err := h.s3.Put(output.FileName, outputBytes); err != nil {
		log.Infof("RunInstructionMessage: failed to upload output to S3: %v", err)
		_ = h.fileRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID)
		return
	}

	_ = h.fileRepo.UpdateStatusAndSize(output.ID, FileStatusDone, int64(len(outputBytes)))
	h.publishFileNotification(instr.UserID, instr.ID)
}

func (h *InstructionHandler) publishFileNotification(userID, instrID primitive.ObjectID) {
	n := InstructionNotification{UserID: userID.Hex(), InstructionID: instrID.Hex()}
	b, err := json.Marshal(n)
	if err != nil {
		log.Infof("publishFileNotification: marshal error: %v", err)
		return
	}
	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectNotificationsSSE, b); err != nil {
		log.Infof("publishFileNotification: publish error: %v", err)
	}
}

func (h *InstructionHandler) CleanInstruction(c *fiber.Ctx) error {
	log.Infof("CleanInstruction invoked")
	return nil
}

// GetInstructionFileBytes finds a file by instructionId and fileId and returns its bytes
func (h *InstructionHandler) GetInstructionFileBytes(c *fiber.Ctx) error {
	instrIDHex := c.Params("id", "")
	fileIDHex := c.Params("fileId", "")
	if instrIDHex == "" || fileIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "instruction id and file id required", "errors": nil, "data": nil})
	}

	instrID, err := primitive.ObjectIDFromHex(instrIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid instruction id", "errors": nil, "data": nil})
	}
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid file id", "errors": nil, "data": nil})
	}

	instr := h.instrRepo.GetByID(instrID)
	if instr == nil || instr.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction not found", "errors": nil, "data": nil})
	}

	// authorization: ensure the requester owns the instruction
	localUserID, _ := c.Locals("UserID").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden", "errors": nil, "data": nil})
	}

	f := h.fileRepo.GetByID(fileID)
	if f == nil || f.ID.IsZero() || f.InstructionID != instr.ID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "file not found", "errors": nil, "data": nil})
	}

	b := h.s3.Get(f.FileName)
	if b == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "file blob not found", "errors": nil, "data": nil})
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", "attachment; filename="+filepath.Base(f.FileName))
	return c.Status(fiber.StatusOK).Send(b)
}
