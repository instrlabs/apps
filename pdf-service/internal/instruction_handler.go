package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/pdf-service/pkg/utils"
	"github.com/instrlabs/shared/modelx"
	"github.com/minio/minio-go/v7"
	natsgo "github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructionHandler struct {
	cfg           *Config
	s3            *minio.Client
	nats          *natsgo.Conn
	instrRepo     *InstructionRepository
	detailRepo    *InstructionDetailRepository
	productClient *ProductClient
	pdfSvc        *PDFService
}

func NewInstructionHandler(
	cfg *Config,
	s3 *minio.Client,
	nats *natsgo.Conn,
	instrRepo *InstructionRepository,
	detailRepo *InstructionDetailRepository,
	productClient *ProductClient,
	pdfSvc *PDFService) *InstructionHandler {
	return &InstructionHandler{cfg: cfg, s3: s3, nats: nats, instrRepo: instrRepo, detailRepo: detailRepo, productClient: productClient, pdfSvc: pdfSvc}
}

func (h *InstructionHandler) putObject(objectName string, data []byte) error {
	ctx := context.Background()
	_, err := h.s3.PutObject(ctx, h.cfg.S3Bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (h *InstructionHandler) getObject(objectName string) ([]byte, error) {
	ctx := context.Background()
	obj, err := h.s3.GetObject(ctx, h.cfg.S3Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

func (h *InstructionHandler) deleteObject(objectName string) error {
	ctx := context.Background()
	return h.s3.RemoveObject(ctx, h.cfg.S3Bucket, objectName, minio.RemoveObjectOptions{})
}

func (h *InstructionHandler) publishFileNotification(userID, instrID, fileID primitive.ObjectID) error {
	pdfNotification := modelx.InstructionNotification{
		UserID:              userID,
		InstructionID:       instrID,
		InstructionDetailID: fileID,
		CreatedAt:           time.Now().UTC(),
	}

	b, _ := json.Marshal(pdfNotification)
	err := h.nats.Publish(h.cfg.NatsSubjectNotificationsSSE, b)

	return err
}

func (h *InstructionHandler) CreateInstruction(c *fiber.Ctx) error {
	type RequestBody struct {
		ProductIDHex string `json:"product_id"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"errors":  nil,
			"data":    nil,
		})
	}

	if body.ProductIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ProductID is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	productID, _ := primitive.ObjectIDFromHex(body.ProductIDHex)
	product, _ := h.productClient.FindByID(productID, "pdf")
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "product not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	userID, _ := c.Locals("userId").(string)
	userObjID, _ := primitive.ObjectIDFromHex(userID)

	instruction := &modelx.Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    userObjID,
		ProductID: productID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := h.instrRepo.Create(instruction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create instruction",
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "instruction creation successful",
		"errors":  nil,
		"data":    fiber.Map{"instruction": instruction},
	})
}

func (h *InstructionHandler) CreateInstructionDetails(c *fiber.Ctx) error {
	instrIDHex := c.Params("id")
	instrID, err := primitive.ObjectIDFromHex(instrIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid instruction id",
			"errors":  nil,
			"data":    nil,
		})
	}

	var instr modelx.Instruction
	if err := h.instrRepo.GetByID(instrID, instr); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "instruction not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	userIDHex, _ := c.Locals("userId").(string)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "user not authorized to create instruction details",
			"errors":  nil,
			"data":    nil,
		})
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to read uploaded file",
			"errors":  nil,
			"data":    nil,
		})
	}

	const maxFileSize = 50 * 1024 * 1024
	if fh.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "fileSize too large",
			"errors":  nil,
			"data":    nil,
		})
	}

	allowedTypes := []string{"application/pdf"}
	contentType := fh.Header.Get("Content-Type")
	if !utils.Contains(allowedTypes, contentType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid file type",
			"errors":  nil,
			"data":    nil,
		})
	}

	f, err := fh.Open()
	if err != nil {
		log.Errorf("CreateInstructionDetails: Failed to open uploaded file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to read uploaded file",
			"errors":  nil,
			"data":    nil,
		})
	}

	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Errorf("CreateInstructionDetails: Failed to read uploaded file content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to read file content",
			"errors":  nil,
			"data":    nil,
		})
	}

	inputID := primitive.NewObjectID()
	outputID := primitive.NewObjectID()
	ext := filepath.Ext(fh.Filename)
	inputFileName := inputID.Hex() + ext
	outputFileName := outputID.Hex() + ext
	inputFilePath := "pdfs/" + inputFileName
	outputFilePath := "pdfs/" + outputFileName

	now := time.Now().UTC()

	input := &modelx.InstructionFile{
		ID:            inputID,
		InstructionID: instr.ID,
		FileName:      inputFilePath,
		FileSize:      int64(len(b)),
		MimeType:      utils.GetMimeTypeFromName(fh.Filename),
		Status:        modelx.InstructionDetailStatusPending,
		OutputID:      &outputID,
		FilePath:      inputFilePath,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	output := &modelx.InstructionFile{
		ID:            outputID,
		InstructionID: instr.ID,
		FileName:      outputFilePath,
		FileSize:      0,
		MimeType:      utils.GetMimeTypeFromName(fh.Filename),
		Status:        modelx.InstructionDetailStatusPending,
		FilePath:      outputFilePath,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.detailRepo.CreateMany([]*modelx.InstructionFile{input, output}); err != nil {
		log.Errorf("CreateInstructionDetails: Failed to create instruction details: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create file records",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	if err := h.putObject(input.FileName, b); err != nil {
		log.Errorf("CreateInstructionDetails: Failed to upload file to S3: %v", err)
		_ = h.detailRepo.UpdateStatus(inputID, modelx.InstructionDetailStatusFailed)
		_ = h.detailRepo.UpdateStatus(outputID, modelx.InstructionDetailStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "upload failed",
			"errors":  []string{fmt.Sprintf("failed to upload file to storage: %v", err)},
			"data":    nil,
		})
	}

	if err := h.nats.Publish(h.cfg.NatsSubjectPdfRequests, []byte(inputID.Hex())); err != nil {
		log.Errorf("CreateInstructionDetails: Failed to queue processing request: %v", err)
		_ = h.detailRepo.UpdateStatus(inputID, modelx.InstructionDetailStatusFailed)
		_ = h.detailRepo.UpdateStatus(outputID, modelx.InstructionDetailStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to queue processing request",
			"errors":  []string{fmt.Sprintf("failed to publish to message queue: %v", err)},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file created",
		"errors":  nil,
		"data":    fiber.Map{"input": input, "output": output},
	})
}

func (h *InstructionHandler) GetInstructionByID(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id",
			"errors":  nil,
			"data":    nil,
		})
	}

	var instr modelx.Instruction
	if err := h.instrRepo.GetByID(id, instr); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "instruction not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data":    fiber.Map{"instruction": instr},
	})
}

func (h *InstructionHandler) GetInstructionDetail(c *fiber.Ctx) error {
	instrIDHex := c.Params("id")
	detailIDHex := c.Params("detail_id")

	if instrIDHex == "" || detailIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "instruction id and detail id required",
			"errors":  nil,
			"data":    nil,
		})
	}

	instrID, err := primitive.ObjectIDFromHex(instrIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid instruction id",
			"errors":  nil,
			"data":    nil,
		})
	}

	detailID, err := primitive.ObjectIDFromHex(detailIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid detail id",
			"errors":  nil,
			"data":    nil,
		})
	}

	var instr modelx.Instruction
	if err := h.instrRepo.GetByID(instrID, instr); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "instruction not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	userIDHex, _ := c.Locals("userId").(string)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
			"errors":  nil,
			"data":    nil,
		})
	}

	var detail modelx.InstructionFile
	if err := h.detailRepo.GetByID(detailID, detail); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "instruction detail not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"errors":  nil,
		"data":    fiber.Map{"detail": detail},
	})
}

func (h *InstructionHandler) GetInstructionDetailFile(c *fiber.Ctx) error {
	instrIDHex := c.Params("id")
	detailIDHex := c.Params("detail_id")
	if instrIDHex == "" || detailIDHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "instruction id and detail id required",
			"errors":  nil,
			"data":    nil,
		})
	}

	instrID, err := primitive.ObjectIDFromHex(instrIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid instruction id",
			"errors":  nil,
			"data":    nil,
		})
	}
	detailID, err := primitive.ObjectIDFromHex(detailIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid detail id",
			"errors":  nil,
			"data":    nil,
		})
	}

	var instr modelx.Instruction
	if err := h.instrRepo.GetByID(instrID, instr); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "instruction not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	userIDHex, _ := c.Locals("userId").(string)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)
	if instr.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
			"errors":  nil,
			"data":    nil,
		})
	}

	var instrDetail modelx.InstructionFile
	if err := h.detailRepo.GetByID(detailID, instrDetail); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "detail not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	b, err := h.getObject(instrDetail.FileName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "file blob not found",
			"errors":  nil,
			"data":    nil,
		})
	}

	c.Set("content-type", "application/octet-stream")
	c.Set("content-disposition", "attachment; filename="+filepath.Base(instrDetail.FileName))

	return c.Status(fiber.StatusOK).Send(b)
}

func (h *InstructionHandler) RunInstructionMessage(data []byte) {
	fileIDHex := string(bytes.TrimSpace(data))
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		log.Errorf("RunInstructionMessage: Invalid file ID: %v", err)
		return
	}

	var input modelx.InstructionFile
	if err := h.detailRepo.GetByID(fileID, input); err != nil {
		log.Info("RunInstructionMessage: Input file not found")
		return
	}

	var output modelx.InstructionFile
	if err := h.detailRepo.GetByID(*input.OutputID, output); err != nil {
		log.Info("RunInstructionMessage: Output file not found")
		return
	}

	var instr modelx.Instruction
	if err := h.instrRepo.GetByID(input.InstructionID, instr); err != nil {
		log.Info("RunInstructionMessage: Instruction not found")
		_ = h.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		return
	}

	product, _ := h.productClient.FindByID(instr.ProductID, "pdf")
	if product == nil {
		log.Info("RunInstructionMessage: Product not found")
		_ = h.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = h.publishFileNotification(instr.UserID, instr.ID, input.ID)
		return
	}

	_ = h.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusProcessing)
	_ = h.publishFileNotification(instr.UserID, instr.ID, input.ID)

	inputBytes, err := h.getObject(input.FilePath)
	if err != nil {
		log.Errorf("RunInstructionMessage: Input file missing on S3: %v", err)
		_ = h.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = h.publishFileNotification(instr.UserID, instr.ID, input.ID)
		return
	}

	outputBytes, err := h.pdfSvc.Compress(inputBytes)
	if err != nil {
		log.Errorf("RunInstructionMessage: PDF processing failed: %v", err)
		_ = h.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = h.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = h.publishFileNotification(instr.UserID, instr.ID, input.ID)
		return
	}

	_ = h.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusSuccess)
	_ = h.publishFileNotification(instr.UserID, instr.ID, input.ID)
	_ = h.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusProcessing)
	_ = h.publishFileNotification(instr.UserID, instr.ID, output.ID)

	if err := h.putObject(output.FileName, outputBytes); err != nil {
		log.Errorf("RunInstructionMessage: Failed to upload output to S3: %v", err)
		_ = h.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = h.publishFileNotification(instr.UserID, instr.ID, output.ID)
		return
	}

	_ = h.detailRepo.UpdateFileSize(output.ID, int64(len(outputBytes)))
	_ = h.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusSuccess)
	_ = h.publishFileNotification(instr.UserID, instr.ID, output.ID)
}

func (h *InstructionHandler) CleanInstruction() {
	log.Info("CleanInstruction: Cleaning instruction")
}
