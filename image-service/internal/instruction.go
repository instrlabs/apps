package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStatus string

const (
	FileStatusFailed     FileStatus = "FAILED"
	FileStatusPending    FileStatus = "PENDING"
	FileStatusProcessing FileStatus = "PROCESSING"
	FileStatusDone       FileStatus = "DONE"
)

type InstructionType string

type Instruction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type InstructionDetail struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	InstructionID primitive.ObjectID  `json:"instruction_id" bson:"instruction_id"`
	FileName      string              `json:"file_name" bson:"file_name"`
	FileSize      int64               `json:"file_size" bson:"file_size"`
	MimeType      string              `json:"mime_type" bson:"mime_type"`
	Status        FileStatus          `json:"status" bson:"status"`
	OutputID      *primitive.ObjectID `json:"output_id,omitempty" bson:"output_id,omitempty"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" bson:"updated_at"`
	IsCleaned     bool                `json:"is_cleaned" bson:"is_cleaned"`
}

type InstructionNotification struct {
	UserID              string `json:"user_id"`
	InstructionID       string `json:"instruction_id"`
	InstructionDetailID string `json:"instruction_detail_id"`
}
