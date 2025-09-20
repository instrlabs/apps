package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStatus string

const (
	FileStatusPending    FileStatus = "PENDING"
	FileStatusProcessing FileStatus = "PROCESSING"
	FileStatusCompleted  FileStatus = "COMPLETED"
	FileStatusFailed     FileStatus = "FAILED"
)

type InstructionType string

type File struct {
	OriginalName string     `json:"original_name" bson:"original_name"`
	FileName     string     `json:"file_name" bson:"file_name"`
	Size         int64      `json:"size" bson:"size"`
	Status       FileStatus `json:"status" bson:"status"`
}

type Instruction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Files     []File             `json:"files" bson:"files"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type InstructionRequest struct {
	InstructionID string `json:"instruction_id"`
}

type InstructionNotification struct {
	InstructionID     string `json:"instruction_id"`
	InstructionStatus string `json:"instruction_status"`
}
