package internal

import (
	"github.com/instrlabs/shared/functionx"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	Origins     string
	JWTSecret   string
	CSRFEnabled bool

	AuthService                 string
	NatsURI                     string
	NatsSubjectNotificationsSSE string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: functionx.GetEnvString("ENVIRONMENT", "development"),
		Port:        functionx.GetEnvString("PORT", ":3000"),
		Origins:     functionx.GetEnvString("ORIGINS_ALLOWED", "http://localhost:8000"),
		JWTSecret:   functionx.GetEnvString("JWT_SECRET", ""),
		CSRFEnabled: functionx.GetEnvBool("CSRF_ENABLED", false),

		AuthService:                 functionx.GetEnvString("AUTH_SERVICE", "http://auth-service:3000"),
		NatsURI:                     functionx.GetEnvString("NATS_URI", "nats://localhost:4222"),
		NatsSubjectNotificationsSSE: functionx.GetEnvString("NATS_SUBJECT_NOTIFICATIONS_SSE", "notifications.sse"),
	}
}
