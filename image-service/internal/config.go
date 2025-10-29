package internal

import (
	initx "github.com/instrlabs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string

	MongoURI string
	MongoDB  string

	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool

	NatsURI                     string
	NatsSubjectImageRequests    string
	NatsSubjectNotificationsSSE string

	ApiUrl string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":3000"),

		MongoURI: initx.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  initx.GetEnv("MONGO_DB", "instrlabs-apps"),

		S3Endpoint:  initx.GetEnv("S3_ENDPOINT", "localhost:9000"),
		S3Region:    initx.GetEnv("S3_REGION", "us-east-1"),
		S3AccessKey: initx.GetEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: initx.GetEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    initx.GetEnv("S3_BUCKET", "instrlabs-apps"),
		S3UseSSL:    initx.GetEnvBool("S3_USE_SSL", false),

		NatsURI:                     initx.GetEnv("NATS_URI", "nats://nats:4222"),
		NatsSubjectImageRequests:    initx.GetEnv("NATS_SUBJECT_IMAGE_REQUESTS", "image.requests"),
		NatsSubjectNotificationsSSE: initx.GetEnv("NATS_SUBJECT_NOTIFICATIONS_SSE", "notifications.sse"),

		ApiUrl: initx.GetEnv("GATEWAY_URL", ""),
	}
}
