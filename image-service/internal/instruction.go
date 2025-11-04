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
	FilePath      string              `json:"file_path" bson:"file_path"`
	IsCleaned     bool                `json:"is_cleaned" bson:"is_cleaned"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" bson:"updated_at"`
}

type InstructionNotification struct {
	UserID              primitive.ObjectID `json:"user_id"`
	InstructionID       primitive.ObjectID `json:"instruction_id"`
	InstructionDetailID primitive.ObjectID `json:"instruction_detail_id"`
	Status              FileStatus         `json:"status"`
	CreatedAt           time.Time          `json:"created_at"`
}
