package internal

import (
	initx "github.com/histweety-labs/shared/init"
	"github.com/joho/godotenv"
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
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", ""),
		Port:        initx.GetEnv("PORT", ":3040"),

		MongoURI: initx.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  initx.GetEnv("MONGO_DB", "payment_service"),

		NatsURL:                    initx.GetEnv("NATS_URL", "nats://localhost:4222"),
		NatsSubjectPaymentEvents:   initx.GetEnv("NATS_SUBJECT_PAYMENT_EVENTS", "payment.events"),
		NatsSubjectPaymentRequests: initx.GetEnv("NATS_SUBJECT_PAYMENT_REQUESTS", "payment.requests"),

		MidtransServerKey:       initx.GetEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey:       initx.GetEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransEnvironment:     initx.GetEnv("MIDTRANS_ENVIRONMENT", "sandbox"),
		MidtransNotificationURL: initx.GetEnv("MIDTRANS_NOTIFICATION_URL", ""),

		CORSAllowedOrigins: initx.GetEnv("CORS_ALLOWED_ORIGINS", "http://web.localhost"),
	}
}

// LoadConfig mirrors auth-service naming and delegates to NewConfig for compatibility
func LoadConfig() *Config {
	return NewConfig()
}
