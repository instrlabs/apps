package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"strconv"
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
	pdfSvc      *PDFService
}

func NewInstructionHandler(
	cfg *Config,
	s3 *initx.S3,
	nats *initx.Nats,
	instrRepo *InstructionRepository,
	detailRepo *InstructionDetailRepository,
	productRepo *ProductRepository,
	pdfSvc *PDFService,
) *InstructionHandler {
	return &InstructionHandler{
		cfg:         cfg,
		s3:          s3,
		nats:        nats,
		instrRepo:   instrRepo,
		detailRepo:  detailRepo,
		productRepo: productRepo,
		pdfSvc:      pdfSvc,
	}
}

func (h *InstructionHandler) CreateInstruction(c *fiber.Ctx) error {
	var req struct {
		ProductID string `json:"productId"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Verify product exists
	product, err := h.productRepo.FindByID(productID, "pdf")
	if err != nil || product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	userID := c.Locals("userId").(string)
	instruction := &Instruction{
		UserID:    userID,
		ProductID: productID,
	}

	createdInstruction, err := h.instrRepo.Create(instruction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create instruction",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Instruction created successfully",
		"errors":  nil,
		"data":    createdInstruction,
	})
}

func (h *InstructionHandler) CreateInstructionDetails(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid instruction ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Get instruction to verify ownership
	instruction, err := h.instrRepo.GetByID(id)
	if err != nil || instruction == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Instruction not found",
			"errors":  err,
			"data":    nil,
		})
	}

	localUserID := c.Locals("userId").(string)
	if instruction.UserID != localUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
			"errors":  nil,
			"data":    nil,
		})
	}

	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "File upload required",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Read file
	f, _ := fh.Open()
	b, _ := io.ReadAll(f)
	_ = f.Close()

	// Validate PDF file
	if err := h.pdfSvc.Validate(b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid PDF file",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	inputID := primitive.NewObjectID()
	outputID := primitive.NewObjectID()
	ext := filepath.Ext(fh.Filename)

	inName := "pdfs/" + inputID.Hex() + "_input" + ext
	outName := "pdfs/" + outputID.Hex() + "_output" + ext

	now := time.Now().UTC()
	input := &InstructionDetail{
		ID:            inputID,
		InstructionID: instruction.ID,
		FileName:      fh.Filename,
		FileSize:      int64(len(b)),
		Status:        FileStatusDone,
		Type:          "input",
		FilePath:      inName,
		InputID:       nil,
		OutputID:      &outputID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	output := &InstructionDetail{
		ID:            outputID,
		InstructionID: instruction.ID,
		FileName:      "compressed_" + fh.Filename,
		FileSize:      0,
		Status:        FileStatusPending,
		Type:          "output",
		FilePath:      outName,
		InputID:       &inputID,
		OutputID:      nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Upload input file to S3
	if err := h.s3.Put(inName, b); err != nil {
		log.Infof("Failed to upload file to S3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload file",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Create details
	details := []InstructionDetail{*input, *output}
	createdDetails, err := h.detailRepo.CreateMany(details)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create instruction details",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Publish NATS message for processing
	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectPdfRequests, []byte(inputID.Hex())); err != nil {
		log.Infof("Failed to publish NATS message: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"errors":  nil,
		"data":    createdDetails,
	})
}

func (h *InstructionHandler) ListInstructions(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	limitStr := c.Query("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	instructions, err := h.instrRepo.ListLatest(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch instructions",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"errors":  nil,
		"data":    instructions,
	})
}

func (h *InstructionHandler) GetInstructionByID(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid instruction ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	instruction, err := h.instrRepo.GetByID(id)
	if err != nil || instruction == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Instruction not found",
			"errors":  err,
			"data":    nil,
		})
	}

	localUserID := c.Locals("userId").(string)
	if instruction.UserID != localUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"errors":  nil,
		"data":    instruction,
	})
}

func (h *InstructionHandler) ListUncleanedFiles(c *fiber.Ctx) error {
	details, err := h.detailRepo.ListUncleaned()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch uncleaned files",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"errors":  nil,
		"data":    details,
	})
}

// RunInstructionMessage processes PDF compression request from NATS message
func (h *InstructionHandler) RunInstructionMessage(data []byte) {
	inputIDHex := string(bytes.TrimSpace(data))
	inputID, err := primitive.ObjectIDFromHex(inputIDHex)
	if err != nil {
		log.Infof("RunInstructionMessage: invalid input ID: %q err=%v", inputIDHex, err)
		return
	}

	// 1. Get input file record
	input, err := h.detailRepo.GetByID(inputID)
	if err != nil || input == nil || input.ID.IsZero() {
		log.Infof("RunInstructionMessage: input file not found: %s", inputIDHex)
		return
	}

	// 2. Get output file record (linked via OutputID)
	if input.OutputID == nil {
		log.Infof("RunInstructionMessage: input file has no output: %s", inputIDHex)
		return
	}

	output, err := h.detailRepo.GetByID(*input.OutputID)
	if err != nil || output == nil || output.ID.IsZero() {
		log.Infof("RunInstructionMessage: output file not found: %s", input.OutputID.Hex())
		return
	}

	// 3. Get instruction for user context
	instruction, err := h.instrRepo.GetByID(input.InstructionID)
	if err != nil || instruction == nil || instruction.ID.IsZero() {
		log.Infof("RunInstructionMessage: instruction not found: %s", input.InstructionID.Hex())
		_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instruction, output)
		return
	}

	// 4. Get input file from S3
	_ = h.detailRepo.UpdateStatus(input.ID, FileStatusProcessing)
	h.publishFileNotification(instruction, input)

	inputBytes := h.s3.Get(input.FilePath)
	if inputBytes == nil {
		log.Infof("RunInstructionMessage: input file missing on S3: %s", input.FilePath)
		_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instruction, input)
		return
	}

	// 5. Compress PDF
	compressedBytes, err := h.pdfSvc.Compress(inputBytes)
	if err != nil {
		log.Infof("RunInstructionMessage: PDF compression failed: %v", err)
		_ = h.detailRepo.UpdateStatus(input.ID, FileStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instruction, input)
		return
	}

	// Input processing is complete
	_ = h.detailRepo.UpdateStatus(input.ID, FileStatusDone)
	h.publishFileNotification(instruction, input)

	// 6. Update output status and upload compressed file
	_ = h.detailRepo.UpdateStatus(output.ID, FileStatusProcessing)
	h.publishFileNotification(instruction, output)

	if err := h.s3.Put(output.FilePath, compressedBytes); err != nil {
		log.Infof("RunInstructionMessage: failed to upload compressed PDF to S3: %v", err)
		_ = h.detailRepo.UpdateStatus(output.ID, FileStatusFailed)
		h.publishFileNotification(instruction, output)
		return
	}

	// Mark output as complete with file size
	_ = h.detailRepo.UpdateStatusAndSize(output.ID, FileStatusDone, int64(len(compressedBytes)))
	h.publishFileNotification(instruction, output)

	log.Infof("RunInstructionMessage: PDF compressed successfully: %s -> %.1f%% reduction",
		inputID.Hex(), float64(len(inputBytes)-len(compressedBytes))/float64(len(inputBytes))*100)
}

func (h *InstructionHandler) CleanInstruction() error {
	cutoff := time.Now().Add(-1 * time.Hour)

	// Find old completed files (DONE status, older than 1 hour)
	files, err := h.detailRepo.ListOlderThan(cutoff)
	if err != nil {
		log.Infof("CleanInstruction: failed to list old files: %v", err)
		return err
	}

	if len(files) > 0 {
		// Delete from S3
		for _, f := range files {
			if f.FilePath == "" {
				continue
			}
			if err := h.s3.Delete(f.FilePath); err != nil {
				log.Infof("CleanInstruction: failed to delete S3 object %s: %v", f.FilePath, err)
			}
		}

		// Mark as cleaned in database
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

	// Find stale pending/processing files (PROCESSING status, updated before cutoff)
	stale, err := h.detailRepo.ListPendingUpdatedBefore(cutoff)
	if err != nil {
		log.Infof("CleanInstruction: failed to list stale files: %v", err)
		return err
	}

	if len(stale) > 0 {
		// Delete stale files from S3 and mark as failed
		for _, f := range stale {
			if f.FilePath != "" {
				if err := h.s3.Delete(f.FilePath); err != nil {
					log.Infof("CleanInstruction: failed to delete stale PROCESSING S3 object %s: %v", f.FilePath, err)
				}
			}
			_ = h.detailRepo.UpdateStatus(f.ID, FileStatusFailed)
		}

		// Mark as cleaned in database
		ids := make([]primitive.ObjectID, 0, len(stale))
		for _, f := range stale {
			if !f.ID.IsZero() {
				ids = append(ids, f.ID)
			}
		}
		if err := h.detailRepo.MarkCleaned(ids); err != nil {
			log.Infof("CleanInstruction: MarkCleaned failed for stale PROCESSING %d files: %v", len(ids), err)
		}
		log.Infof("CleanInstruction: marked FAILED and cleaned %d stale PROCESSING files (updated_at < %s)", len(ids), cutoff.UTC().Format(time.RFC3339))
	}

	return nil
}

func (h *InstructionHandler) GetInstructionDetails(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid instruction ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Verify instruction exists and user owns it
	instruction, err := h.instrRepo.GetByID(id)
	if err != nil || instruction == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Instruction not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	localUserID := c.Locals("userId").(string)
	if instruction.UserID != localUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
			"errors":  nil,
			"data":    nil,
		})
	}

	// List all details for this instruction
	details, err := h.detailRepo.ListByInstruction(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch instruction details",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"errors":  nil,
		"data":    details,
	})
}

func (h *InstructionHandler) GetInstructionDetail(c *fiber.Ctx) error {
	idHex := c.Params("id")
	detailIDHex := c.Params("detailId")

	if idHex == "" || detailIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Instruction ID and detail ID required",
			"errors":  nil,
			"data":    nil,
		})
	}

	instrID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid instruction ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	detailID, err := primitive.ObjectIDFromHex(detailIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid detail ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Verify instruction exists and user owns it
	instruction, err := h.instrRepo.GetByID(instrID)
	if err != nil || instruction == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Instruction not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	localUserID := c.Locals("userId").(string)
	if instruction.UserID != localUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Get the specific detail
	detail, err := h.detailRepo.GetByID(detailID)
	if err != nil || detail == nil || detail.ID.IsZero() || detail.InstructionID != instrID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Instruction detail not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"errors":  nil,
		"data":    detail,
	})
}

func (h *InstructionHandler) GetInstructionDetailFile(c *fiber.Ctx) error {
	idHex := c.Params("id")
	detailIDHex := c.Params("detailId")

	if idHex == "" || detailIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Instruction ID and detail ID required",
			"errors":  nil,
			"data":    nil,
		})
	}

	instrID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid instruction ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	detailID, err := primitive.ObjectIDFromHex(detailIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid detail ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Verify instruction exists and user owns it
	instruction, err := h.instrRepo.GetByID(instrID)
	if err != nil || instruction == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Instruction not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	localUserID := c.Locals("userId").(string)
	if instruction.UserID != localUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Get the detail record
	detail, err := h.detailRepo.GetByID(detailID)
	if err != nil || detail == nil || detail.ID.IsZero() || detail.InstructionID != instrID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "File not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Fetch file from S3
	fileBytes := h.s3.Get(detail.FilePath)
	if fileBytes == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "File blob not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Set response headers for file download
	c.Set("content-type", "application/octet-stream")
	c.Set("content-disposition", "attachment; filename="+filepath.Base(detail.FileName))
	return c.Status(fiber.StatusOK).Send(fileBytes)
}

func (h *InstructionHandler) publishFileNotification(instruction *Instruction, detail *InstructionDetail) {
	if instruction == nil || detail == nil {
		return
	}

	notification := InstructionNotification{
		UserID:              instruction.UserID,
		InstructionID:       instruction.ID,
		InstructionDetailID: detail.ID,
		Status:              detail.Status,
		Type:                detail.Type,
		CreatedAt:           time.Now().UTC(),
	}

	notifBytes, err := json.Marshal(notification)
	if err != nil {
		log.Infof("publishFileNotification: failed to marshal notification: %v", err)
		return
	}

	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectNotificationsSSE, notifBytes); err != nil {
		log.Infof("publishFileNotification: failed to publish notification: %v", err)
	}
}
