package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
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

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Warn("PORT not set, using default: ", port)
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service.localhost"
		log.Warn("AUTH_SERVICE_URL not set, using default: ", authServiceURL)
	}

	labsServiceURL := os.Getenv("LABS_SERVICE_URL")
	if labsServiceURL == "" {
		labsServiceURL = "http://labs-service.localhost"
		log.Warn("AUTH_SERVICE_URL not set, using default: ", labsServiceURL)
	}

	return Config{
		Port: port,
		Services: []ServiceConfig{
			{
				Name:   "auth-service",
				URL:    authServiceURL,
				Prefix: "/auth",
			},
			{
				Name:   "labs-service",
				URL:    labsServiceURL,
				Prefix: "/labs",
			},
		},
	}
}
