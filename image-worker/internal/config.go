package internal

import initx "github.com/histweety-labs/shared/init"

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

	ImageServiceURL string

	// Scheduler settings
	SchedulerIntervalMs  int
	PendingMaxAgeMinutes int
	RetryMax             int
	RetryBatchSize       int
	RetryLockTTLMinutes  int
}

func NewConfig() *Config {
	return &Config{
		Environment:              initx.GetEnv("ENVIRONMENT", "development"),
		Port:                     initx.GetEnv("PORT", ":3000"),
		NatsURL:                  initx.GetEnv("NATS_URL", "nats://localhost:4222"),
		NatsSubjectRequests:      initx.GetEnv("NATS_SUBJECT_REQUESTS", "image.requests"),
		NatsSubjectNotifications: initx.GetEnv("NATS_SUBJECT_NOTIFICATIONS", "image.notifications"),

		MongoURI: initx.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  initx.GetEnv("MONGO_DB", "labs"),

		S3Endpoint:  initx.GetEnv("S3_ENDPOINT", "localhost:9000"),
		S3Region:    initx.GetEnv("S3_REGION", "us-east-1"),
		S3AccessKey: initx.GetEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: initx.GetEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    initx.GetEnv("S3_BUCKET", "labs"),
		S3UseSSL:    initx.GetEnvBool("S3_USE_SSL", false),

		ImageServiceURL: initx.GetEnv("IMAGE_SERVICE_URL", "http://localhost:3001"),

		// Scheduler defaults (overrides can be wired later if needed)
		SchedulerIntervalMs:  initx.GetEnvInt("SCHEDULER_INTERVAL_MS", 60000),
		PendingMaxAgeMinutes: initx.GetEnvInt("PENDING_MAX_AGE_MINUTES", 30),
		RetryMax:             initx.GetEnvInt("RETRY_MAX", 3),
		RetryBatchSize:       initx.GetEnvInt("RETRY_BATCH_SIZE", 50),
		RetryLockTTLMinutes:  initx.GetEnvInt("RETRY_LOCK_TTL_MINUTES", 5),
	}
}
