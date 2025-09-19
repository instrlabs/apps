package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstructionStatus string

const (
	InstructionStatusPending    InstructionStatus = "PENDING"
	InstructionStatusProcessing InstructionStatus = "PROCESSING"
	InstructionStatusCompleted  InstructionStatus = "COMPLETED"
	InstructionStatusFailed     InstructionStatus = "FAILED"
)

type InstructionType string

type File struct {
	FileName string `json:"file_name" bson:"file_name"`
	Size     int64  `json:"size" bson:"size"`
}

type Instruction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Inputs    []File             `json:"inputs" bson:"inputs"`
	Outputs   []File             `json:"outputs" bson:"outputs"`
	Status    InstructionStatus  `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type InstructionRequest struct {
	UserID        string `json:"user_id"`
	InstructionID string `json:"instruction_id"`
}

type InstructionNotification struct {
	UserID            string `json:"user_id"`
	InstructionID     string `json:"instruction_id"`
	InstructionStatus string `json:"instruction_status"`
}
