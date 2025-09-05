package internal

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var w *WorkerPool

type WorkerPool struct {
	jobs       chan string
	ctx        context.Context
	cancel     context.CancelFunc
	s3         *S3Service
	instrRepo  *InstructionRepository
	fileRepo   *FileRepository
	productSvc *ProductService
}

func NewWorkerPool(n int, s3 *S3Service, instrRepo *InstructionRepository, fileRepo *FileRepository, productSvc *ProductService) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPool{
		jobs:       make(chan string, 100),
		ctx:        ctx,
		cancel:     cancel,
		s3:         s3,
		instrRepo:  instrRepo,
		fileRepo:   fileRepo,
		productSvc: productSvc,
	}

	// expose to package-level for helpers that need access without changing signatures
	w = wp

	for i := 0; i < n; i++ {
		go func(id int) {
			for {
				select {
				case <-wp.ctx.Done():
					return
				case instrID := <-wp.jobs:
					log.Infof("worker-%d processing instructionID=%s", id, instrID)
					wp.processInstruction(instrID)
				}
			}
		}(i)
	}

	return wp
}

func (w *WorkerPool) processInstruction(instructionHex string) {
	if w.instrRepo == nil || w.s3 == nil {
		log.Errorf("worker missing dependencies; skipping instruction %s", instructionHex)
		return
	}

	id, err := primitive.ObjectIDFromHex(instructionHex)
	if err != nil {
		log.Errorf("invalid instruction id %s: %v", instructionHex, err)
		return
	}

	instr := w.instrRepo.GetByID(id)
	if instr == nil || instr.ID.IsZero() {
		log.Errorf("failed to get instruction %s", instructionHex)
		return
	}

	product := w.productSvc.GetProduct("image-compress")
	if product == nil {
		log.Errorf("failed to get product for instruction %s", instructionHex)
		return
	}

	_ = w.instrRepo.UpdateStatus(context.Background(), id, InstructionStatusProcessing)

	err = CompressJPEG(instructionHex)
	if err != nil {
		log.Errorf("failed to compress image for %s: %v", instructionHex, err)
		_ = w.instrRepo.UpdateStatus(context.Background(), id, InstructionStatusFailed)
		return
	}

	_ = w.instrRepo.UpdateStatus(context.Background(), id, InstructionStatusCompleted)
}

func (w *WorkerPool) Enqueue(instructionID string) {
	select {
	case <-w.ctx.Done():
		return
	default:
		w.jobs <- instructionID
	}
}

func (w *WorkerPool) Stop() { w.cancel() }
