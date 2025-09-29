package internal

import (
	initx "github.com/instrlabs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string

	Origins   string
	JWTSecret string

	NatsURL                     string
	NatsSubjectNotificationsSSE string
}

func NewConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":3001"),

		Origins:   initx.GetEnv("CORS_ALLOWED_ORIGINS", "http://localhost:8000"),
		JWTSecret: initx.GetEnv("JWT_SECRET", ""),

		NatsURL:                     initx.GetEnv("NATS_URL", "nats://localhost:4222"),
		NatsSubjectNotificationsSSE: initx.GetEnv("NATS_SUBJECT_NOTIFICATIONS_SSE", "notifications.sse"),
	}
}
