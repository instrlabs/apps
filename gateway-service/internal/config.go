package internal

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type ServiceConfig struct {
	Name   string
	URL    string
	Prefix string
}

type Config struct {
	Port     string
	Services []ServiceConfig
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file, using default environment variables")
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Warn("PORT not set, using default: ", port)
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth.localhost"
		log.Warn("AUTH_SERVICE_URL not set, using default: ", authServiceURL)
	}

	paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentServiceURL == "" {
		paymentServiceURL = "http://payment.localhost"
		log.Warn("PAYMENT_SERVICE_URL not set, using default: ", paymentServiceURL)
	}

	return &Config{
		Port: port,
		Services: []ServiceConfig{
			{
				Name:   "auth-service",
				URL:    authServiceURL,
				Prefix: "/auth",
			},
			{
				Name:   "payment-service",
				URL:    paymentServiceURL,
				Prefix: "/payment",
			},
		},
	}
}
