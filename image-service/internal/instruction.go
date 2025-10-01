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

type InstructionNotification struct {
	UserID        string `json:"user_id"`
	InstructionID string `json:"instruction_id"`
}
