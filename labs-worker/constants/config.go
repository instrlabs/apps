package constants

import (
	"os"
)

type Config struct {
	Environment string

	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool

	NatsURL               string
	NatsSubjectPDFJobs    string
	NatsSubjectPDFResults string
}

func NewConfig() *Config {
	env := getEnv("ENVIRONMENT", "development")

	s3Endpoint := getEnv("S3_ENDPOINT", "localhost:9000")
	s3Region := getEnv("S3_REGION", "us-east-1")
	s3AccessKey := getEnv("S3_ACCESS_KEY", "minioadmin")
	s3SecretKey := getEnv("S3_SECRET_KEY", "minioadmin")
	s3Bucket := getEnv("S3_BUCKET", "labs")
	s3UseSSL := getEnvBool("S3_USE_SSL", false)

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	natsSubjectPDFJobs := getEnv("NATS_SUBJECT_PDF_JOBS", "pdf.jobs")
	natsSubjectPDFResults := getEnv("NATS_SUBJECT_PDF_RESULTS", "pdf.results")

	return &Config{
		Environment: env,

		S3Endpoint:  s3Endpoint,
		S3Region:    s3Region,
		S3AccessKey: s3AccessKey,
		S3SecretKey: s3SecretKey,
		S3Bucket:    s3Bucket,
		S3UseSSL:    s3UseSSL,

		NatsURL:               natsURL,
		NatsSubjectPDFJobs:    natsSubjectPDFJobs,
		NatsSubjectPDFResults: natsSubjectPDFResults,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if value == "true" || value == "1" || value == "yes" {
			return true
		} else if value == "false" || value == "0" || value == "no" {
			return false
		}
	}
	return fallback
}
