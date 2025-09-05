package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileType string

const (
	FileTypeRequest  FileType = "REQUEST"
	FileTypeResponse FileType = "RESPONSE"
)

type File struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InstructionID primitive.ObjectID `json:"instruction_id" bson:"instruction_id"`
	Type          FileType           `json:"type" bson:"type"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}
