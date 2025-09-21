package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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
	productKey := c.Params("product_key")
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
	fileName := "images/" + inputID.Hex() + ext

	fileDoc := &File{
		ID:            inputID,
		InstructionID: instr.ID,
		OriginalName:  fh.Filename,
		FileName:      fileName,
		Size:          int64(len(b)),
		Status:        FileStatusUploading,
		OutputID:      outputID,
	}
	if err := h.fileRepo.CreateOne(fileDoc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create file record", "errors": nil, "data": nil})
	}

	if err := h.s3.Put(fileDoc.FileName, b); err != nil {
		_ = h.fileRepo.UpdateStatus(inputID, FileStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "upload failed", "errors": nil, "data": nil})
	}

	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectImagesRequests, []byte(inputID.Hex())); err != nil {
		_ = h.fileRepo.UpdateStatus(inputID, FileStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to publish to nats", "errors": nil, "data": nil})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file created",
		"errors":  nil,
		"data":    fiber.Map{"file": fileDoc},
	})
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
	case "images-compress":
		outputBytes, err = h.imageSvc.Compress(inputBytes)
		if err != nil {
			log.Printf("RunInstructionMessage: image-compress failed: %v", err)
			_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
			h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
			return
		} else {
			_ = h.fileRepo.UpdateStatus(input.ID, FileStatusDone)
			h.publishFileNotification(instr.ID, input.ID, FileStatusDone)
		}
	default:
		log.Printf("RunInstructionMessage: unsupported product key: %s", product.Key)
		_ = h.fileRepo.UpdateStatus(input.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, input.ID, FileStatusFailed)
		return
	}

	// 5b. Create output file doc and link it from input
	ext := filepath.Ext(input.FileName)
	outName := "images/" + input.OutputID.Hex() + ext
	outFile := &File{
		ID:            input.OutputID,
		InstructionID: instr.ID,
		OriginalName:  filepath.Base(outName),
		FileName:      outName,
		Size:          int64(len(outputBytes)),
		Status:        FileStatusProcessing,
	}
	if err := h.fileRepo.CreateOne(outFile); err != nil {
		log.Printf("RunInstructionMessage: failed to create output file doc: %v", err)
		_ = h.fileRepo.UpdateStatus(outFile.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, outFile.ID, FileStatusFailed)
		return
	} else {
		_ = h.fileRepo.UpdateStatus(outFile.ID, FileStatusUploading)
		h.publishFileNotification(instr.ID, outFile.ID, FileStatusUploading)
	}

	// 6. Upload output to S3
	if err := h.s3.Put(outFile.FileName, outputBytes); err != nil {
		log.Printf("RunInstructionMessage: failed to upload output to S3: %v", err)
		_ = h.fileRepo.UpdateStatus(outFile.ID, FileStatusFailed)
		h.publishFileNotification(instr.ID, outFile.ID, FileStatusFailed)
		return
	} else {
		_ = h.fileRepo.UpdateStatus(outFile.ID, FileStatusDone)
		h.publishFileNotification(instr.ID, outFile.ID, FileStatusDone)
	}
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
