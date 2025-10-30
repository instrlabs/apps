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

	AuthService                 string
	NatsURI                     string
	NatsSubjectNotificationsSSE string
}

func NewConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":3001"),

		Origins:   initx.GetEnv("ORIGINS_ALLOWED", "http://localhost:8000"),
		JWTSecret: initx.GetEnv("JWT_SECRET", ""),

		AuthService:                 initx.GetEnv("AUTH_SERVICE", "http://auth-service:3000"),
		NatsURI:                     initx.GetEnv("NATS_URI", "nats://localhost:4222"),
		NatsSubjectNotificationsSSE: initx.GetEnv("NATS_SUBJECT_NOTIFICATIONS_SSE", "notifications.sse"),
	}
}
