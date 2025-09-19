package internal

import (
	initx "github.com/histweety-labs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	MongoURI    string
	MongoDB     string

	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool

	NatsURL                  string
	NatsSubjectRequests      string
	NatsSubjectNotifications string

	PaymentServiceURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":3000"),
		MongoURI:    initx.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:     initx.GetEnv("MONGO_DB", "labs"),

		S3Endpoint:  initx.GetEnv("S3_ENDPOINT", "localhost:9000"),
		S3Region:    initx.GetEnv("S3_REGION", "us-east-1"),
		S3AccessKey: initx.GetEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: initx.GetEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    initx.GetEnv("S3_BUCKET", "labs"),
		S3UseSSL:    initx.GetEnvBool("S3_USE_SSL", false),

		NatsURL:                  initx.GetEnv("NATS_URL", "nats://localhost:4222"),
		NatsSubjectRequests:      initx.GetEnv("NATS_SUBJECT_REQUESTS", "image.requests"),
		NatsSubjectNotifications: initx.GetEnv("NATS_SUBJECT_NOTIFICATIONS", "image.notifications"),

		PaymentServiceURL: initx.GetEnv("PAYMENT_SERVICE_URL", "http://localhost:3000/payment"),
	}
}
