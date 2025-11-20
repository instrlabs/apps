package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/shared/modelx"
	"github.com/minio/minio-go/v7"
	natsgo "github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructionProcessor struct {
	cfg           *Config
	s3            *minio.Client
	nats          *natsgo.Conn
	instrRepo     *InstructionRepository
	detailRepo    *InstructionDetailRepository
	productClient *ProductClient
	imageSvc      *ImageService
}

func NewInstructionProcessor(
	cfg *Config,
	s3 *minio.Client,
	nats *natsgo.Conn,
	instrRepo *InstructionRepository,
	detailRepo *InstructionDetailRepository,
	productClient *ProductClient,
	imageSvc *ImageService) *InstructionProcessor {
	return &InstructionProcessor{
		cfg:           cfg,
		s3:            s3,
		nats:          nats,
		instrRepo:     instrRepo,
		detailRepo:    detailRepo,
		productClient: productClient,
		imageSvc:      imageSvc,
	}
}

func (p *InstructionProcessor) putObject(objectName string, data []byte) error {
	ctx := context.Background()
	_, err := p.s3.PutObject(ctx, p.cfg.S3Bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (p *InstructionProcessor) getObject(objectName string) ([]byte, error) {
	ctx := context.Background()
	obj, err := p.s3.GetObject(ctx, p.cfg.S3Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

func (p *InstructionProcessor) publishFileNotification(userID *primitive.ObjectID, guestID *string, instrID, fileID primitive.ObjectID) error {
	imageNotification := modelx.InstructionNotification{
		UserID:              userID,
		GuestID:             guestID,
		InstructionID:       instrID,
		InstructionDetailID: fileID,
		CreatedAt:           time.Now().UTC(),
	}

	b, _ := json.Marshal(imageNotification)
	err := p.nats.Publish(p.cfg.NatsSubjectNotificationsSSE, b)

	return err
}

func (p *InstructionProcessor) RunInstructionMessage(data []byte) {
	fileIDHex := string(bytes.TrimSpace(data))
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		log.Errorf("RunInstructionMessage: Invalid file ID: %v", err)
		return
	}

	var input modelx.InstructionFile
	if err := p.detailRepo.GetByID(fileID, input); err != nil {
		log.Info("RunInstructionMessage: Input file not found")
		return
	}

	var output modelx.InstructionFile
	if err := p.detailRepo.GetByID(*input.OutputID, output); err != nil {
		log.Info("RunInstructionMessage: Output file not found")
		return
	}

	var instr modelx.Instruction
	if err := p.instrRepo.GetByID(input.InstructionID, instr); err != nil {
		log.Info("RunInstructionMessage: Instruction not found")
		_ = p.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = p.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		return
	}

	product, _ := p.productClient.FindByID(instr.ProductID, "image")
	if product == nil {
		log.Info("RunInstructionMessage: Product not found")
		_ = p.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = p.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, input.ID)
		return
	}

	_ = p.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusProcessing)
	_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, input.ID)

	inputBytes, err := p.getObject(input.FilePath)
	if err != nil {
		log.Errorf("RunInstructionMessage: Input file missing on S3: %v", err)
		_ = p.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = p.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, input.ID)
		return
	}

	outputBytes, err := p.imageSvc.Run(product.Key, inputBytes)
	if err != nil {
		log.Errorf("RunInstructionMessage: Image processing failed: %v", err)
		_ = p.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusFailed)
		_ = p.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, input.ID)
		return
	}

	_ = p.detailRepo.UpdateStatus(input.ID, modelx.InstructionDetailStatusSuccess)
	_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, input.ID)
	_ = p.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusProcessing)
	_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, output.ID)

	if err := p.putObject(output.FileName, outputBytes); err != nil {
		log.Errorf("RunInstructionMessage: Failed to upload output to S3: %v", err)
		_ = p.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusFailed)
		_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, output.ID)
		return
	}

	_ = p.detailRepo.UpdateFileSize(output.ID, int64(len(outputBytes)))
	_ = p.detailRepo.UpdateStatus(output.ID, modelx.InstructionDetailStatusSuccess)
	_ = p.publishFileNotification(instr.UserID, instr.GuestID, instr.ID, output.ID)
}

func (p *InstructionProcessor) CleanInstruction() {
	log.Info("CleanInstruction: Cleaning instruction")
}
