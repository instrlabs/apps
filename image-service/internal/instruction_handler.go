package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"mime"
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
	fileRepo    *FileRepository
	productServ *PaymentService
	imageSvc    *ImageService
}

func NewInstructionHandler(
	cfg *Config,
	s3 *initx.S3,
	nats *initx.Nats,
	instrRepo *InstructionRepository,
	fileRepo *FileRepository,
	productServ *PaymentService,
	imageSvc *ImageService) *InstructionHandler {
	return &InstructionHandler{cfg: cfg, s3: s3, nats: nats, instrRepo: instrRepo, fileRepo: fileRepo, productServ: productServ, imageSvc: imageSvc}
}

func (h *InstructionHandler) CreateInstruction(c *fiber.Ctx) error {
	productKey := c.Params("product") + c.Params("key")
	userID, _ := c.Locals("UserID").(string)
	product := h.productServ.GetProduct(userID, productKey)
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "product not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	instructionID := primitive.NewObjectID()
	objUserID, _ := primitive.ObjectIDFromHex(userID)
	objProductID, _ := primitive.ObjectIDFromHex(product.ID)

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
	// Not implemented against new model yet; return empty list to keep API stable
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok", "errors": nil, "data": []Instruction{}})
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

	// Validate file belongs to this instruction by scanning files
	files := h.fileRepo.ListByInstruction(instr.ID)
	found := false
	for _, f := range files {
		if f.FileName == fileName {
			found = true
			break
		}
	}
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "file not found in instruction", "errors": nil, "data": nil})
	}

	// Retrieve from S3
	b := h.s3.Get(fileName)
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

func (h *InstructionHandler) RunInstructionMessage(data []byte) {
	fileIDHex := string(bytes.TrimSpace(data))
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		log.Printf("RunInstructionMessage: invalid file id: %q err=%v", fileIDHex, err)
		return
	}

	// 1. Find file by ID
	input := h.fileRepo.GetByID(fileID)
	if input == nil || input.ID.IsZero() {
		log.Printf("RunInstructionMessage: file not found: %s", fileIDHex)
		return
	}

	// 2. Find instruction
	instr := h.instrRepo.GetByID(input.InstructionID)
	if instr == nil || instr.ID.IsZero() {
		log.Printf("RunInstructionMessage: instruction not found for file=%s", fileIDHex)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(input.InstructionID, input.ID, FileStatusFailed)
		return
	}

	// 3. Find product by instruction's product ID
	product := h.productServ.GetProduct(instr.UserID.Hex(), instr.ProductID.Hex())
	if product == nil {
		log.Printf("RunInstructionMessage: product not found for instruction=%s", instr.ID.Hex())
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
		return
	}

	// 4. Get binary from S3
	_ = h.fileRepo.UpdateStatus(input.ID, FileStatusProcessing)
	h.publishFileNotification(instr.ID, input.ID, FileStatusProcessing)
	inputBytes := h.s3.Get(input.FileName)
	if inputBytes == nil {
		log.Printf("RunInstructionMessage: input file missing on S3: %s", input.FileName)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
		return
	}

	// 5. Process based on product key
	var outputBytes []byte
	switch product.Key {
	case "image-compress":
		outputBytes, err = h.imageSvc.Compress(inputBytes)
		if err != nil {
			log.Printf("RunInstructionMessage: image-compress failed: %v", err)
			_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
			h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
			return
		}
	default:
		log.Printf("RunInstructionMessage: unsupported product key: %s", product.Key)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
		return
	}

	// 5b. Create output file doc and link it from input
	outID := primitive.NewObjectID()
	ext := filepath.Ext(input.FileName)
	outName := "images/" + outID.Hex() + ext
	outFile := &File{
		ID:            outID,
		InstructionID: instr.ID,
		OriginalName:  filepath.Base(outName),
		FileName:      outName,
		Size:          int64(len(outputBytes)),
		Status:        FileStatusDone,
		OutputID:      primitive.NilObjectID,
	}
	if err := h.fileRepo.CreateOne(outFile); err != nil {
		log.Printf("RunInstructionMessage: failed to create output file doc: %v", err)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
		return
	}
	if err := h.fileRepo.LinkOutput(input.ID, outID); err != nil {
		log.Printf("RunInstructionMessage: failed to link output: %v", err)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
		return
	}

	// 6. Upload output to S3
	if err := h.s3.Put(outFile.FileName, outputBytes); err != nil {
		log.Printf("RunInstructionMessage: failed to upload output to S3: %v", err)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
		return
	}

	// Publish DONE for output file
	h.publishFileNotification(instr.ID, outFile.ID, FileStatusDone)

	_ = h.fileRepo.UpdateStatus(input.ID, FileStatusDone)
	h.publishFileNotification(instr.ID, input.ID, FileStatusDone)
}

func (h *InstructionHandler) publishFileNotification(instrID, fileID primitive.ObjectID, status FileStatus) {
	n := FileNotification{InstructionID: instrID.Hex(), FileID: fileID.Hex(), FileStatus: string(status)}
	b, err := json.Marshal(n)
	if err != nil {
		log.Printf("publishFileNotification: marshal error: %v", err)
		return
	}
	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectNotificationsSSE, b); err != nil {
		log.Printf("publishFileNotification: publish error: %v", err)
	}
}

func (h *InstructionHandler) CleanInstruction(c *fiber.Ctx) error {
	log.Printf("CleanInstruction invoked")
	return nil
}
