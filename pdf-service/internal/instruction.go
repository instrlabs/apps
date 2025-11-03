package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instruction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"userId" bson:"userId"`
	ProductID primitive.ObjectID `json:"productId" bson:"productId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type InstructionDetail struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	InstructionID primitive.ObjectID  `json:"-" bson:"instructionId"`
	FileName      string              `json:"fileName" bson:"fileName"`
	FileSize      int64               `json:"fileSize" bson:"fileSize"`
	Status        FileStatus          `json:"status" bson:"status"`
	Type          string              `json:"type" bson:"type"`            // "input" or "output"
	InputID       *primitive.ObjectID `json:"-" bson:"inputId,omitempty"`  // Links output to input
	OutputID      *primitive.ObjectID `json:"-" bson:"outputId,omitempty"` // Links input to output
	FilePath      string              `json:"filePath" bson:"filePath"`    // S3 file path
	CreatedAt     time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt" bson:"updatedAt"`
}

type FileStatus string

const (
	FileStatusPending    FileStatus = "PENDING"
	FileStatusProcessing FileStatus = "PROCESSING"
	FileStatusDone       FileStatus = "DONE"
	FileStatusFailed     FileStatus = "FAILED"
)

type InstructionNotification struct {
	UserID              string             `json:"userId"`
	InstructionID       primitive.ObjectID `json:"instructionId"`
	InstructionDetailID primitive.ObjectID `json:"instructionDetailId"`
	Status              FileStatus         `json:"status"`
	Type                string             `json:"type"`
	CreatedAt           time.Time          `json:"createdAt"`
}
