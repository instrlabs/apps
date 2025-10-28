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
	detailRepo  *InstructionDetailRepository
	productRepo *ProductRepository
	imageSvc    *ImageService
}

func NewInstructionHandler(
	cfg *Config,
	s3 *initx.S3,
	nats *initx.Nats,
	instrRepo *InstructionRepository,
	detailRepo *InstructionDetailRepository,
	productRepo *ProductRepository,
	imageSvc *ImageService) *InstructionHandler {
	return &InstructionHandler{cfg: cfg, s3: s3, nats: nats, instrRepo: instrRepo, detailRepo: detailRepo, productRepo: productRepo, imageSvc: imageSvc}
}

func (h *InstructionHandler) CreateInstruction(c *fiber.Ctx) error {
	type payload struct {
		ProductID string `json:"product_id"`
	}
	var body payload
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"errors":  nil,
			"data":    nil,
		})
	}
	if body.ProductID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ProductID is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	userID, _ := c.Locals("userId").(string)
	productID, _ := primitive.ObjectIDFromHex(body.ProductID)
	product, _ := h.productRepo.FindByID(productID)
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
	userId, _ := c.Locals("userId").(string)
	instructions, err := h.instrRepo.ListLatest(userId, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to list instructions",
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data": map[string]interface{}{
			"instructions": instructions,
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

	localUserID, _ := c.Locals("userId").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden", "errors": nil, "data": nil})
	}

	files := h.detailRepo.ListByInstruction(instr.ID)
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

	localUserID, _ := c.Locals("userId").(string)
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

	now := time.Now().UTC()
	input := &InstructionDetail{
		ID:            inputID,
		InstructionID: instr.ID,
		FileName:      inName,
		FileSize:      int64(len(b)),
		MimeType:      fh.Header.Get("Content-Type"),
		Status:        FileStatusPending,
		OutputID:      &outputID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	output := &InstructionDetail{
		ID:            outputID,
		InstructionID: instr.ID,
		FileName:      outName,
		FileSize:      0,
		MimeType:      fh.Header.Get("Content-Type"),
		Status:        FileStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.detailRepo.CreateMany([]*InstructionDetail{input, output}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create file records",
			"errors":  nil,
			"data":    nil,
		})
	}

	if err := h.s3.Put(input.FileName, b); err != nil {
		_ = h.detailRepo.UpdateStatus(inputID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(outputID, FileStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "upload failed",
			"errors":  nil,
			"data":    nil,
		})
	}

	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectImageRequests, []byte(inputID.Hex())); err != nil {
		_ = h.detailRepo.UpdateStatus(inputID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(outputID, FileStatusFailed)
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
	input := h.detailRepo.GetByID(fileID)
	if input == nil || input.ID.IsZero() {
		log.Infof("RunInstructionMessage: input file not found: %s", fileIDHex)
		return
	}

	output := h.detailRepo.GetByID(*input.OutputID)
	if output == nil || output.ID.IsZero() {
		log.Infof("RunInstructionMessage: output file not found: %s", fileIDHex)
		return
	}

	// 2. Find instruction
	instr := h.instrRepo.GetByID(input.InstructionID)
	if instr == nil || instr.ID.IsZero() {
		log.Infof("RunInstructionMessage: instruction not found: %s", fileIDHex)
		_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(primitive.NilObjectID, input.InstructionID, input.ID)
		return
	}

	// 3. Find product by instruction's product ID
	product, _ := h.productRepo.FindByID(instr.ProductID)
	if product == nil {
		log.Infof("RunInstructionMessage: product not found: %s", instr.ProductID.Hex())
		_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID, input.ID)
		return
	}

	// 4. Get binary from S3
	_ = h.detailRepo.UpdateStatus(input.ID, FileStatusProcessing)
	h.publishFileNotification(instr.UserID, instr.ID, input.ID)
	inputBytes := h.s3.Get(input.FileName)
	if inputBytes == nil {
		log.Infof("RunInstructionMessage: input file missing on S3: %s", input.FileName)
		_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID, input.ID)
		return
	}

	// 5. Process based on product key
	var outputBytes []byte
	switch product.Key {
	case "images/compress":
		outputBytes, err = h.imageSvc.Compress(inputBytes)
		if err != nil {
			log.Infof("RunInstructionMessage: image-compress failed: %v", err)
			_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
			_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
			h.publishFileNotification(instr.UserID, instr.ID, input.ID)
			return
		}
	default:
		log.Infof("RunInstructionMessage: unsupported product key: %s", product.Key)
		_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID, input.ID)
		return
	}

	// Input file is DONE - notify for input completion
	_ = h.detailRepo.UpdateStatus(input.ID, FileStatusDone)
	h.publishFileNotification(instr.UserID, instr.ID, input.ID)

	// Output file is now PROCESSING - notify for output processing start
	_ = h.detailRepo.UpdateStatus(output.ID, FileStatusProcessing)
	h.publishFileNotification(instr.UserID, instr.ID, output.ID)

	// 6. Upload output to S3
	if err := h.s3.Put(output.FileName, outputBytes); err != nil {
		log.Infof("RunInstructionMessage: failed to upload output to S3: %v", err)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instr.UserID, instr.ID, output.ID)
		return
	}

	// Output file is DONE - notify for output completion
	_ = h.detailRepo.UpdateStatusAndSize(output.ID, FileStatusDone, int64(len(outputBytes)))
	h.publishFileNotification(instr.UserID, instr.ID, output.ID)
}

func (h *InstructionHandler) publishFileNotification(userID, instrID, fileID primitive.ObjectID) {
	n := InstructionNotification{UserID: userID.Hex(), InstructionID: instrID.Hex(), InstructionDetailID: fileID.Hex()}
	b, err := json.Marshal(n)
	if err != nil {
		log.Infof("publishFileNotification: marshal error: %v", err)
		return
	}
	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectNotificationsSSE, b); err != nil {
		log.Infof("publishFileNotification: publish error: %v", err)
	}
}

func (h *InstructionHandler) CleanInstruction() error {
	cutoff := time.Now().Add(-1 * time.Hour)

	files := h.detailRepo.ListOlderThan(cutoff)
	if len(files) > 0 {
		for _, f := range files {
			if f.FileName == "" {
				continue
			}
			if err := h.s3.Delete(f.FileName); err != nil {
				log.Infof("CleanInstruction: failed to delete S3 object %s: %v", f.FileName, err)
			}
		}

		ids := make([]primitive.ObjectID, 0, len(files))
		for _, f := range files {
			if !f.ID.IsZero() {
				ids = append(ids, f.ID)
			}
		}
		if err := h.detailRepo.MarkCleaned(ids); err != nil {
			log.Infof("CleanInstruction: MarkCleaned failed for %d files: %v", len(ids), err)
		}
		log.Infof("CleanInstruction: marked cleaned %d files older than %s", len(ids), cutoff.UTC().Format(time.RFC3339))
	} else {
		log.Infof("CleanInstruction: no files older than %s", cutoff.UTC().Format(time.RFC3339))
	}

	stale := h.detailRepo.ListPendingUpdatedBefore(cutoff)
	if len(stale) > 0 {
		for _, f := range stale {
			if f.FileName != "" {
				if err := h.s3.Delete(f.FileName); err != nil {
					log.Infof("CleanInstruction: failed to delete stale PENDING S3 object %s: %v", f.FileName, err)
				}
			}
			_ = h.detailRepo.UpdateStatus(f.ID, FileStatusFailed)
		}
		ids := make([]primitive.ObjectID, 0, len(stale))
		for _, f := range stale {
			if !f.ID.IsZero() {
				ids = append(ids, f.ID)
			}
		}
		if err := h.detailRepo.MarkCleaned(ids); err != nil {
			log.Infof("CleanInstruction: MarkCleaned failed for stale PENDING %d files: %v", len(ids), err)
		}
		log.Infof("CleanInstruction: marked FAILED and cleaned %d stale PENDING files (updated_at < %s)", len(ids), cutoff.UTC().Format(time.RFC3339))
	}

	return nil
}

func (h *InstructionHandler) GetInstructionDetilFile(c *fiber.Ctx) error {
	instrIDHex := c.Params("id", "")
	detailIDHex := c.Params("detailId", "")
	if instrIDHex == "" || detailIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "instruction id and detail id required", "errors": nil, "data": nil})
	}

	instrID, err := primitive.ObjectIDFromHex(instrIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid instruction id", "errors": nil, "data": nil})
	}
	fileID, err := primitive.ObjectIDFromHex(detailIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid detail id", "errors": nil, "data": nil})
	}

	instr := h.instrRepo.GetByID(instrID)
	if instr == nil || instr.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction not found", "errors": nil, "data": nil})
	}

	localUserID, _ := c.Locals("userId").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden", "errors": nil, "data": nil})
	}

	f := h.detailRepo.GetByID(fileID)
	if f == nil || f.ID.IsZero() || f.InstructionID != instr.ID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "file not found", "errors": nil, "data": nil})
	}

	b := h.s3.Get(f.FileName)
	if b == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "file blob not found", "errors": nil, "data": nil})
	}

	c.Set("content-type", "application/octet-stream")
	c.Set("content-disposition", "attachment; filename="+filepath.Base(f.FileName))
	return c.Status(fiber.StatusOK).Send(b)
}

func (h *InstructionHandler) ListUncleanedFiles(c *fiber.Ctx) error {
	files := h.detailRepo.ListUncleaned()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data": fiber.Map{
			"files": files,
		},
	})
}

func (h *InstructionHandler) GetInstructionDetail(c *fiber.Ctx) error {
	instrIDHex := c.Params("id")
	detailIDHex := c.Params("detailId")

	if instrIDHex == "" || detailIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "instruction id and detail id required", "errors": nil, "data": nil})
	}

	instrID, err := primitive.ObjectIDFromHex(instrIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid instruction id", "errors": nil, "data": nil})
	}

	detailID, err := primitive.ObjectIDFromHex(detailIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid detail id", "errors": nil, "data": nil})
	}

	// Verify instruction exists and user has permission
	instr := h.instrRepo.GetByID(instrID)
	if instr == nil || instr.ID.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction not found", "errors": nil, "data": nil})
	}

	localUserID, _ := c.Locals("userId").(string)
	userID, _ := primitive.ObjectIDFromHex(localUserID)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "forbidden", "errors": nil, "data": nil})
	}

	detail := h.detailRepo.GetByID(detailID)
	if detail == nil || detail.ID.IsZero() || detail.InstructionID != instr.ID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "instruction detail not found", "errors": nil, "data": nil})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data":    fiber.Map{"detail": detail},
	})
}
