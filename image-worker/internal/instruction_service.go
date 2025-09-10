package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type InstructionStatus string

const (
	InstructionStatusPending    InstructionStatus = "PENDING"
	InstructionStatusProcessing InstructionStatus = "PROCESSING"
	InstructionStatusCompleted  InstructionStatus = "COMPLETED"
	InstructionStatusFailed     InstructionStatus = "FAILED"
)

type File struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	FileName string `json:"file_name" bson:"file_name"`
	Type     string `json:"type" bson:"type"`
	Size     int64  `json:"size" bson:"size"`
}

type InstructionService struct {
	baseURL string
	client  *http.Client
}

type Instruction struct {
	ID      string            `json:"id"`
	Inputs  []File            `json:"inputs"`
	Outputs []File            `json:"outputs"`
	Status  InstructionStatus `json:"status"`
}

func NewInstructionService(cfg *Config) *InstructionService {
	tr := &http.Transport{
		DialContext:         (&net.Dialer{Timeout: 2 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 2 * time.Second,
	}
	return &InstructionService{
		baseURL: cfg.ImageServiceURL,
		client:  &http.Client{Timeout: 5 * time.Second, Transport: tr},
	}
}

func (s *InstructionService) UpdateStatus(ctx context.Context, id string, status InstructionStatus) error {
	url := fmt.Sprintf("%s/instructions/%s/status", s.baseURL, id)
	body := map[string]string{"status": string(status)}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("image-service returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *InstructionService) UpdateOutputs(ctx context.Context, id string, outputs []File) error {
	url := fmt.Sprintf("%s/instructions/%s/outputs", s.baseURL, id)
	body := map[string]any{"outputs": outputs}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("image-service returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *InstructionService) GetInstruction(ctx context.Context, id string) *Instruction {
	url := fmt.Sprintf("%s/instructions/%s", s.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("failed to create request for instruction %s: %v", id, err)
		return nil
	}
	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("failed to get instruction %s: %v", id, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("image-service returned status %d", resp.StatusCode)
		log.Printf("error getting instruction %s: %v", id, err)
		return nil
	}
	var envelope struct {
		Message string       `json:"message"`
		Errors  any          `json:"errors"`
		Data    *Instruction `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		log.Printf("failed to decode instruction %s: %v", id, err)
		return nil
	}
	return envelope.Data
}
