package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
)

type JobType string

type Job struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OriginalFilename string             `json:"original_filename" bson:"original_filename"`
	JobType          JobType            `json:"job_type" bson:"job_type"`
	Status           JobStatus          `json:"status" bson:"status"`
	Error            string             `json:"error,omitempty" bson:"error,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type JobNotificationMessage struct {
	ID     string    `json:"id"`
	Status JobStatus `json:"status"`
}

type JobResponse struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	JobType          string    `json:"job_type"`
	Status           string    `json:"status"`
	Error            string    `json:"error,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (j *Job) ToResponse() JobResponse {
	return JobResponse{
		ID:               j.ID.Hex(),
		OriginalFilename: j.OriginalFilename,
		JobType:          string(j.JobType),
		Status:           string(j.Status),
		Error:            j.Error,
		CreatedAt:        j.CreatedAt,
		UpdatedAt:        j.UpdatedAt,
	}
}
