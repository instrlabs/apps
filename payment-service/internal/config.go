package internal

import (
	"os"
)

type Config struct {
	Environment string
	Port        string

	// Database
	MongoURI string
	MongoDB  string

	// NATS
	NatsURL                    string
	NatsSubjectPaymentEvents   string
	NatsSubjectPaymentRequests string

	// Midtrans
	MidtransServerKey       string
	MidtransClientKey       string
	MidtransEnvironment     string // sandbox or production
	MidtransNotificationURL string

	// CORS
	CORSAllowedOrigins string
}

func NewConfig() *Config {
	env := getEnv("ENVIRONMENT", "development")
	port := getEnv("PORT", ":3040")

	// Database
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "payment_service")

	// NATS
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	natsSubjectPaymentEvents := getEnv("NATS_SUBJECT_PAYMENT_EVENTS", "payment.events")
	natsSubjectPaymentRequests := getEnv("NATS_SUBJECT_PAYMENT_REQUESTS", "payment.requests")

	// Midtrans
	midtransServerKey := getEnv("MIDTRANS_SERVER_KEY", "")
	midtransClientKey := getEnv("MIDTRANS_CLIENT_KEY", "")
	midtransEnvironment := getEnv("MIDTRANS_ENVIRONMENT", "sandbox")
	midtransNotificationURL := getEnv("MIDTRANS_NOTIFICATION_URL", "")

	// CORS
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://web.localhost")

	return &Config{
		Environment: env,
		Port:        port,

		MongoURI: mongoURI,
		MongoDB:  mongoDB,

		NatsURL:                    natsURL,
		NatsSubjectPaymentEvents:   natsSubjectPaymentEvents,
		NatsSubjectPaymentRequests: natsSubjectPaymentRequests,

		MidtransServerKey:       midtransServerKey,
		MidtransClientKey:       midtransClientKey,
		MidtransEnvironment:     midtransEnvironment,
		MidtransNotificationURL: midtransNotificationURL,

		CORSAllowedOrigins: corsOrigins,
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
