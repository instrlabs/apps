package internal

import (
	"io"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	initx "github.com/histweety-labs/shared/init"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileHandler struct {
	cfg       *Config
	s3        *initx.S3
	nats      *initx.Nats
	instrRepo *InstructionRepository
	fileRepo  *FileRepository
}

func NewFileHandler(cfg *Config, s3 *initx.S3, nats *initx.Nats, instrRepo *InstructionRepository, fileRepo *FileRepository) *FileHandler {
	return &FileHandler{cfg: cfg, s3: s3, nats: nats, instrRepo: instrRepo, fileRepo: fileRepo}
}

func (h *FileHandler) CreateFile(c *fiber.Ctx) error {
	instrIDHex := c.Params("instruction_id", "")
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

	fileID := primitive.NewObjectID()
	ext := filepath.Ext(fh.Filename)
	fileName := "images/" + fileID.Hex() + ext

	fileDoc := &File{
		ID:            fileID,
		InstructionID: instr.ID,
		OriginalName:  fh.Filename,
		FileName:      fileName,
		Size:          int64(len(b)),
		Status:        FileStatusUploading,
		OutputID:      primitive.NilObjectID,
	}
	if err := h.fileRepo.CreateOne(fileDoc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create file record", "errors": nil, "data": nil})
	}

	if err := h.s3.Put(fileDoc.FileName, b); err != nil {
		_ = h.fileRepo.UpdateStatus(fileID, FileStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "upload failed", "errors": nil, "data": nil})
	}

	if err := h.nats.Conn.Publish(h.cfg.NatsSubjectImagesRequests, []byte(fileID.Hex())); err != nil {
		_ = h.fileRepo.UpdateStatus(fileID, FileStatusFailed)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to publish to nats", "errors": nil, "data": nil})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file created",
		"errors":  nil,
		"data":    fiber.Map{"file": fileDoc},
	})
}
