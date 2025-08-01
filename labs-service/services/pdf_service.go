package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"labs-service/constants"
	"labs-service/models"
	"net/http"
	"time"
)

type PDFService struct {
	baseURL string
}

type PDFCreateJobRequest struct {
	JobID     string `json:"job_id"`
	Operation string `json:"operation"`
	S3Path    string `json:"s3_path,omitempty"`
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

func (s *PDFService) CreateJob(ctx *fiber.Ctx, jobID string, jobType models.JobType, s3Path string) (*PDFCreateJobResponse, error) {
	jsonData, err := json.Marshal(interface{}(PDFCreateJobRequest{
		JobID:     jobID,
		Operation: string(jobType),
		S3Path:    s3Path,
	}))
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(s.baseURL+"/jobs", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response PDFCreateJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func NewPDFService(cfg *constants.Config) *PDFService {
	return &PDFService{
		baseURL: cfg.PDFServiceURL,
	}
}
