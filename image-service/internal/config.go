package internal

import (
	"github.com/instrlabs/shared/functionx"
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

	ApiUrl            string
	ProductServiceURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: functionx.GetEnvString("ENVIRONMENT", "development"),
		Port:        functionx.GetEnvString("PORT", ":3000"),

		MongoURI: functionx.GetEnvString("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  functionx.GetEnvString("MONGO_DB", "instrlabs-apps"),

		S3Endpoint:  functionx.GetEnvString("S3_ENDPOINT", "localhost:9000"),
		S3Region:    functionx.GetEnvString("S3_REGION", "us-east-1"),
		S3AccessKey: functionx.GetEnvString("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: functionx.GetEnvString("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    functionx.GetEnvString("S3_BUCKET", "instrlabs-apps"),
		S3UseSSL:    functionx.GetEnvBool("S3_USE_SSL", false),

		NatsURI:                     functionx.GetEnvString("NATS_URI", "nats://nats:4222"),
		NatsSubjectImageRequests:    functionx.GetEnvString("NATS_SUBJECT_IMAGE_REQUESTS", "image.requests"),
		NatsSubjectNotificationsSSE: functionx.GetEnvString("NATS_SUBJECT_NOTIFICATIONS_SSE", "notifications.sse"),

		ApiUrl:            functionx.GetEnvString("API_URL", ""),
		ProductServiceURL: functionx.GetEnvString("PRODUCT_SERVICE_URL", "http://product-service:3005"),
	}
}
