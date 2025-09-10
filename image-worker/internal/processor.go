package internal

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	initx "github.com/histweety-labs/shared/init"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Processor struct {
	mongo     *initx.Mongo
	s3        *initx.S3
	nats      *initx.Nats
	imgServ   *ImageService
	instrRepo *InstructionRepository
}

func NewProcessor(mongo *initx.Mongo, s3 *initx.S3, nats *initx.Nats, imgServ *ImageService, instrRepo *InstructionRepository) *Processor {
	return &Processor{
		mongo:     mongo,
		s3:        s3,
		nats:      nats,
		imgServ:   imgServ,
		instrRepo: instrRepo,
	}
}

func (p *Processor) Handle(ctx context.Context, job *JobMessage) error {
	log.Printf("jobID: %v, userID: %v", job.ID, job.UserID)

	id, err := primitive.ObjectIDFromHex(job.ID)
	if err != nil {
		return fmt.Errorf("invalid job id: %w", err)
	}

	if err := p.instrRepo.UpdateStatus(id, InstructionStatusProcessing); err != nil {
		log.Printf("failed to set PROCESSING for %s: %v", job.ID, err)
	}

	instr := p.instrRepo.GetByID(id)
	if instr == nil || instr.ID.IsZero() {
		_ = p.instrRepo.UpdateStatus(id, InstructionStatusFailed)
		return fmt.Errorf("instruction not found")
	}

	var outputs []File
	for idx, in := range instr.Inputs {
		data := p.s3.Get(in.FileName)
		if data == nil {
			_ = p.instrRepo.UpdateStatus(id, InstructionStatusFailed)
			return fmt.Errorf("missing input: %s", in.FileName)
		}
		compressed, err := p.imgServ.Compress(data)
		if err != nil {
			_ = p.instrRepo.UpdateStatus(id, InstructionStatusFailed)
			return fmt.Errorf("compress: %w", err)
		}
		ext := filepath.Ext(in.FileName)
		outName := fmt.Sprintf("images/%s-out-%d%s", job.ID, idx, ext)
		if err := p.s3.Put(outName, compressed); err != nil {
			_ = p.instrRepo.UpdateStatus(id, InstructionStatusFailed)
			return fmt.Errorf("s3 put: %w", err)
		}
		outputs = append(outputs, File{FileName: outName, Size: int64(len(compressed))})
	}

	if err := p.instrRepo.UpdateOutputs(id, outputs); err != nil {
		_ = p.instrRepo.UpdateStatus(id, InstructionStatusFailed)
		return fmt.Errorf("update outputs: %w", err)
	}

	if p.instrRepo != nil {
		if err := p.instrRepo.UpdateStatus(id, InstructionStatusCompleted); err != nil {
			log.Printf("failed to set COMPLETED for %s: %v", job.ID, err)
		}
	}
	return nil
}
