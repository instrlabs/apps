package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type PDFService struct {
	baseURL string
}

type PDFCreateJobRequest struct {
	JobID     string  `json:"job_id"`
	Operation JobType `json:"operation"`
	Filename  string  `json:"filename"`
	FileSize  int64   `json:"file_size"`
	S3Path    string  `json:"s3_path"`
}

type PDFCreateJobResponse struct {
	ID             string    `json:"id"`
	Filename       string    `json:"filename"`
	FileSize       int64     `json:"file_size"`
	Operation      string    `json:"operation"`
	JobID          string    `json:"job_id"`
	OutputFilePath string    `json:"output_file_path,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s *PDFService) CreateJob(ctx *fiber.Ctx, req PDFCreateJobRequest) (*PDFCreateJobResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(s.baseURL+"/jobs", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted &&
		resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Data PDFCreateJobResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func NewPDFService(cfg *Config) *PDFService {
	return &PDFService{
		baseURL: cfg.PDFServiceURL,
	}
}
