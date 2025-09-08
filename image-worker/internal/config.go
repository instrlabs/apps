package internal

import "os"

type Config struct {
	Environment              string
	Port                     string
	NatsURL                  string
	NatsSubjectRequests      string
	NatsSubjectNotifications string

	MongoURI string
	MongoDB  string

	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool
}

func NewConfig() *Config {
	return &Config{
		Environment:              getEnv("ENVIRONMENT", "development"),
		Port:                     getEnv("PORT", ":3000"),
		NatsURL:                  getEnv("NATS_URL", "nats://localhost:4222"),
		NatsSubjectRequests:      getEnv("NATS_SUBJECT_REQUESTS", "image.requests"),
		NatsSubjectNotifications: getEnv("NATS_SUBJECT_NOTIFICATIONS", "image.notifications"),

		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "labs"),

		S3Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
		S3Region:    getEnv("S3_REGION", "us-east-1"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    getEnv("S3_BUCKET", "labs"),
		S3UseSSL:    getEnvBool("S3_USE_SSL", false),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		switch v {
		case "true", "1", "yes", "TRUE", "Y", "y":
			return true
		case "false", "0", "no", "FALSE", "N", "n":
			return false
		}
	}
	return fallback
}
