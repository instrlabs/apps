package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDFOperation string

const PDFJobSubject = "pdf.jobs"

const (
	PDFOperationConvertToJPG PDFOperation = "PDF/CONVERT_TO_JPG"
	PDFOperationCompress     PDFOperation = "PDF/COMPRESS"
	PDFOperationMerge        PDFOperation = "PDF/MERGE"
	PDFOperationSplit        PDFOperation = "PDF/SPLIT"
)

type PDFJob struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Filename       string             `json:"filename" bson:"filename"`
	FileSize       int64              `json:"file_size" bson:"file_size"`
	S3Path         string             `json:"s3_path" bson:"s3_path"`
	Operation      PDFOperation       `json:"operation" bson:"operation"`
	JobID          string             `json:"job_id" bson:"job_id"`
	OutputFilePath string             `json:"output_file_path,omitempty" bson:"output_file_path,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type PDFJobMessage struct {
	ID string `json:"id"`
}

type UpdatePDFJobRequest struct {
	OutputFilePath string `json:"output_file_path,omitempty"`
	Status         string `json:"status,omitempty"`
	Error          string `json:"error,omitempty"`
}
