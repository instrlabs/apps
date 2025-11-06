package internal

import (
	"github.com/instrlabs/shared/functionx"
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
		Environment: functionx.GetEnvString("ENVIRONMENT", "development"),
		Port:        functionx.GetEnvString("PORT", ":3001"),

		Origins:   functionx.GetEnvString("ORIGINS_ALLOWED", "http://localhost:8000"),
		JWTSecret: functionx.GetEnvString("JWT_SECRET", ""),

		AuthService:                 functionx.GetEnvString("AUTH_SERVICE", "http://auth-service:3000"),
		NatsURI:                     functionx.GetEnvString("NATS_URI", "nats://localhost:4222"),
		NatsSubjectNotificationsSSE: functionx.GetEnvString("NATS_SUBJECT_NOTIFICATIONS_SSE", "notifications.sse"),
	}
}
