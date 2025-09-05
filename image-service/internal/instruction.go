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

type Instruction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Status    InstructionStatus  `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
