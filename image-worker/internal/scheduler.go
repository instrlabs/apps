package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	initx "github.com/histweety-labs/shared/init"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Scheduler struct {
	instrRepo *InstructionRepository
	nats      *initx.Nats
	cfg       *Config
	stopCh    chan struct{}
}

func NewScheduler(repo *InstructionRepository, nats *initx.Nats, cfg *Config) *Scheduler {
	return &Scheduler{
		instrRepo: repo,
		nats:      nats,
		cfg:       cfg,
		stopCh:    make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	interval := time.Duration(s.cfg.SchedulerIntervalMs) * time.Millisecond
	log.Printf("scheduler: starting with interval %s", interval)
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := s.tick(); err != nil {
					log.Printf("scheduler tick error: %v", err)
				}
			case <-s.stopCh:
				log.Printf("scheduler: stopping")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() { close(s.stopCh) }

func (s *Scheduler) tick() error {
	batchSize := int64(s.cfg.RetryBatchSize)
	maxAgeMin := s.cfg.PendingMaxAgeMinutes
	retryMax := s.cfg.RetryMax
	lockTTLMin := s.cfg.RetryLockTTLMinutes

	now := time.Now().UTC()
	claimed := int64(0)
	for claimed < batchSize {
		doc, err := s.instrRepo.ClaimNextPendingForRetry(now, maxAgeMin, retryMax, lockTTLMin)
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				break
			}
			return err
		}
		if doc == nil || len(doc) == 0 {
			break
		}

		claimed++
		go s.processClaim(doc) // process concurrently but lightweight
	}
	return nil
}

func (s *Scheduler) processClaim(doc bson.M) {
	id, _ := doc["_id"].(primitive.ObjectID)
	userID, _ := doc["user_id"].(primitive.ObjectID)
	retryCount := int64(0)
	if v, ok := doc["retry_count"]; ok {
		switch x := v.(type) {
		case int32:
			retryCount = int64(x)
		case int64:
			retryCount = x
		}
	}

	job := &JobMessage{ID: id.Hex(), UserID: userID.Hex()}
	msg := &nats.Msg{Subject: s.cfg.NatsSubjectRequests}
	msg.Data = mustJSON(job)
	if msg.Header == nil {
		msg.Header = nats.Header{}
	}
	msg.Header.Set("Nats-Msg-Id", "instr-"+id.Hex()+"-retry-"+itoa(int(retryCount+1)))

	if s.nats == nil || s.nats.Conn == nil {
		log.Printf("scheduler: NATS connection unavailable; releasing lock for %s", id.Hex())
		s.releaseLock(id)
		return
	}

	if err := s.nats.Conn.PublishMsg(msg); err != nil {
		log.Printf("scheduler: publish failed for %s: %v", id.Hex(), err)
		s.releaseLock(id)
		return
	}

	// Mark as retried if still pending
	s.markRetried(id)
}

func (s *Scheduler) releaseLock(id primitive.ObjectID) {
	_ = s.instrRepo.ReleaseRetryLock(id)
}

func (s *Scheduler) markRetried(id primitive.ObjectID) {
	now := time.Now().UTC()
	if err := s.instrRepo.MarkRetried(id, now); err != nil {
		log.Printf("scheduler: markRetried update error for %s: %v", id.Hex(), err)
	}
}

// Helper functions
func itoa(i int) string { return fmt.Sprintf("%d", i) }

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
