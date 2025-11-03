package internal

import (
	"io"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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
		log.Printf("Failed to upload file to S3: %v", err)
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
		log.Printf("Failed to publish NATS message: %v", err)
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

// Minimal implementation for other methods to compile
func (h *InstructionHandler) RunInstructionMessage(data []byte) {
	// TODO: Implement PDF processing logic
	log.Printf("Received message: %s", string(data))
}

func (h *InstructionHandler) CleanInstruction() error {
	log.Printf("Starting PDF cleanup job")
	return nil
}

func (h *InstructionHandler) GetInstructionDetails(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented yet"})
}

func (h *InstructionHandler) GetInstructionDetail(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented yet"})
}

func (h *InstructionHandler) GetInstructionDetailFile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented yet"})
}
