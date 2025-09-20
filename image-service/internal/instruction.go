package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStatus string

const (
	FileStatusFailed     FileStatus = "FAILED"
	FileStatusPending    FileStatus = "PENDING"
	FileStatusUploading  FileStatus = "UPLOADING"
	FileStatusProcessing FileStatus = "PROCESSING"
	FileStatusDone       FileStatus = "DONE"
)

type InstructionType string

type File struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InstructionID primitive.ObjectID `json:"instruction_id" bson:"instruction_id"`
	OriginalName  string             `json:"original_name" bson:"original_name"`
	Size          int64              `json:"size" bson:"size"`
	Status        FileStatus         `json:"status" bson:"status"`
	OutputID      primitive.ObjectID `json:"output_id" bson:"output_id"`
}

type Instruction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
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
