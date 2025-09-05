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

	return &Config{
		Port: port,
		Services: []ServiceConfig{
			{
				Name:   "auth-service",
				URL:    os.Getenv("AUTH_SERVICE_URL"),
				Prefix: "/auth",
			},
			{
				Name:   "payment-service",
				URL:    os.Getenv("PAYMENT_SERVICE_URL"),
				Prefix: "/payment",
			},
			{
				Name:   "image-service",
				URL:    os.Getenv("IMAGE_SERVICE_URL"),
				Prefix: "/image",
			},
		},
	}
}
