package internal

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Port        string
	MongoURI    string
	MongoDB     string

	// S3 configuration
	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool

	// NATS configuration
	NatsURL                     string
	NatsSubjectJobNotifications string

	// PDF Service configuration
	PDFServiceURL string
}

func NewConfig() *Config {
	env := getEnv("ENVIRONMENT", "development")
	port := getEnv("PORT", ":8080")
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "labs_service")

	// S3 configuration
	s3Endpoint := getEnv("S3_ENDPOINT", "localhost:9000")
	s3Region := getEnv("S3_REGION", "us-east-1")
	s3AccessKey := getEnv("S3_ACCESS_KEY", "minioadmin")
	s3SecretKey := getEnv("S3_SECRET_KEY", "minioadmin")
	s3Bucket := getEnv("S3_BUCKET", "pdfs")
	s3UseSSL := getEnvBool("S3_USE_SSL", false)

	// NATS configuration
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	natsSubjectJobNotifications := getEnv("NATS_SUBJECT_JOB_NOTIFICATIONS", "job.notifications")

	// PDF Service configuration
	pdfServiceURL := getEnv("PDF_SERVICE_URL", "http://pdf-service:3000")

	return &Config{
		Environment: env,
		Port:        port,
		MongoURI:    mongoURI,
		MongoDB:     mongoDB,

		S3Endpoint:  s3Endpoint,
		S3Region:    s3Region,
		S3AccessKey: s3AccessKey,
		S3SecretKey: s3SecretKey,
		S3Bucket:    s3Bucket,
		S3UseSSL:    s3UseSSL,

		NatsURL:                     natsURL,
		NatsSubjectJobNotifications: natsSubjectJobNotifications,

		PDFServiceURL: pdfServiceURL,
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
