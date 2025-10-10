package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	InstructionID primitive.ObjectID  `json:"instruction_id" bson:"instruction_id"`
	OriginalName  string              `json:"original_name" bson:"original_name"`
	FileName      string              `json:"file_name" bson:"file_name"`
	Size          int64               `json:"size" bson:"size"`
	Status        FileStatus          `json:"status" bson:"status"`
	OutputID      *primitive.ObjectID `json:"output_id,omitempty" bson:"output_id,omitempty"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" bson:"updated_at"`
	IsCleaned     bool                `json:"is_cleaned" bson:"is_cleaned"`
}
