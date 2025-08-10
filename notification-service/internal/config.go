package internal

import (
	"os"
)

type Config struct {
	Environment string
	Port        string

	NatsURL                     string
	NatsSubjectJobNotifications string

	WebSocketPath string
}

func NewConfig() *Config {
	env := getEnv("ENVIRONMENT", "development")
	port := getEnv("PORT", ":3030")

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	natsSubjectJobNotifications := getEnv("NATS_SUBJECT_JOB_NOTIFICATIONS", "job.notifications")

	webSocketPath := getEnv("WEBSOCKET_PATH", "/ws")

	return &Config{
		Environment: env,
		Port:        port,

		NatsURL:                     natsURL,
		NatsSubjectJobNotifications: natsSubjectJobNotifications,

		WebSocketPath: webSocketPath,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if value == "true" || value == "1" || value == "yes" {
			return true
		} else if value == "false" || value == "0" || value == "no" {
			return false
		}
	}
	return fallback
}
