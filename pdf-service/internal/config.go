package internal

import (
	initx "github.com/instrlabs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	// Environment
	Environment string
	Port        string

	// Database
	MongoURI string
	MongoDB  string

	// S3
	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool

	// NATS
	NatsURI                     string
	NatsSubjectPdfRequests      string
	NatsSubjectNotificationsSSE string

	// API
	ApiUrl            string
	ProductServiceURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", "3004"),

		MongoURI: initx.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  initx.GetEnv("MONGO_DB", "instrlabs"),

		S3Endpoint:  initx.GetEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Region:    initx.GetEnv("S3_REGION", "us-east-1"),
		S3AccessKey: initx.GetEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: initx.GetEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    initx.GetEnv("S3_BUCKET", "instrlabs"),
		S3UseSSL:    initx.GetEnvBool("S3_USE_SSL", false),

		NatsURI:                     initx.GetEnv("NATS_URI", "nats://localhost:4222"),
		NatsSubjectPdfRequests:      initx.GetEnv("NATS_SUBJECT_PDF_REQUESTS", "pdf.requests"),
		NatsSubjectNotificationsSSE: initx.GetEnv("NATS_SUBJECT_NOTIFICATIONS_SSE", "notifications.sse"),

		ApiUrl:            initx.GetEnv("API_URL", "http://localhost:3000"),
		ProductServiceURL: initx.GetEnv("PRODUCT_SERVICE_URL", "http://product-service:3005"),
	}
}
