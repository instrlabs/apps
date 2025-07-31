package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDFOperation string

const (
	PDFOperationToJPG    PDFOperation = "TO_JPG"
	PDFOperationCompress PDFOperation = "COMPRESS"
	PDFOperationMerge    PDFOperation = "MERGE"
	PDFOperationSplit    PDFOperation = "SPLIT"
)

type PDFJob struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OriginalName   string             `json:"original_name" bson:"original_name"`
	FileSize       int64              `json:"file_size" bson:"file_size"`
	PageCount      int                `json:"page_count" bson:"page_count"`
	S3Path         string             `json:"s3_path" bson:"s3_path"`
	Operation      PDFOperation       `json:"operation" bson:"operation"`
	JobID          string             `json:"job_id" bson:"job_id"`
	OutputFilePath string             `json:"output_file_path,omitempty" bson:"output_file_path,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type UpdatePDFJobRequest struct {
	OutputFilePath string `json:"output_file_path,omitempty"`
}
