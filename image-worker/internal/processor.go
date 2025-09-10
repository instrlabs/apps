package internal

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	initx "github.com/histweety-labs/shared/init"
)

type Processor struct {
	mongo     *initx.Mongo
	s3        *initx.S3
	nats      *initx.Nats
	imgServ   *ImageService
	instrServ *InstructionService
}

func NewProcessor(mongo *initx.Mongo, s3 *initx.S3, nats *initx.Nats, imgServ *ImageService, instrServ *InstructionService) *Processor {
	return &Processor{
		mongo:     mongo,
		s3:        s3,
		nats:      nats,
		imgServ:   imgServ,
		instrServ: instrServ,
	}
}

func (p *Processor) Handle(ctx context.Context, job *JobMessage) error {
	log.Printf("jobID: %v, userID: %v", job.ID, job.UserID)

	if err := p.instrServ.UpdateStatus(ctx, job, InstructionStatusProcessing); err != nil {
		log.Printf("failed to set PROCESSING for %s: %v", job.ID, err)
	}

	instr := p.instrServ.GetInstruction(ctx, job)
	if instr == nil {
		_ = p.instrServ.UpdateStatus(ctx, job, InstructionStatusFailed)
		return fmt.Errorf("instruction not found")
	}

	var outputs []File
	for idx, in := range instr.Inputs {
		data := p.s3.Get(in.FileName)
		if data == nil {
			_ = p.instrServ.UpdateStatus(ctx, job, InstructionStatusFailed)
			return fmt.Errorf("missing input: %s", in.FileName)
		}
		compressed, err := p.imgServ.Compress(data)
		if err != nil {
			_ = p.instrServ.UpdateStatus(ctx, job, InstructionStatusFailed)
			return fmt.Errorf("compress: %w", err)
		}
		ext := filepath.Ext(in.FileName)
		outName := fmt.Sprintf("images/%s-out-%d%s", job.ID, idx, ext)
		if err := p.s3.Put(outName, compressed); err != nil {
			_ = p.instrServ.UpdateStatus(ctx, job, InstructionStatusFailed)
			return fmt.Errorf("s3 put: %w", err)
		}
		outputs = append(outputs, File{FileName: outName, Size: int64(len(compressed))})
	}

	if err := p.instrServ.UpdateOutputs(ctx, job, outputs); err != nil {
		_ = p.instrServ.UpdateStatus(ctx, job, InstructionStatusFailed)
		return fmt.Errorf("update outputs: %w", err)
	}

	if p.instrServ != nil {
		if err := p.instrServ.UpdateStatus(ctx, job, InstructionStatusCompleted); err != nil {
			log.Printf("failed to set COMPLETED for %s: %v", job.ID, err)
		}
	}
	return nil
}
