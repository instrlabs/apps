package internal

import (
	initx "github.com/histweety-labs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string

	NatsURL                  string
	NatsSubjectNotifications string
}

func NewConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":3001"),

		NatsURL:                  initx.GetEnv("NATS_URL", "nats://localhost:4222"),
		NatsSubjectNotifications: initx.GetEnv("NATS_SUBJECT_NOTIFICATIONS", "notifications.sse"),
	}
}
