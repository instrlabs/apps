package internal

import (
	initx "github.com/instrlabs/shared/init"
	"github.com/joho/godotenv"
)

type ServiceConfig struct {
	Name   string
	URL    string
	Prefix string
}

type Config struct {
	Environment string
	Port        string
	Origins     string
	JWTSecret   string
	Services    []ServiceConfig
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":3000"),
		Origins:     initx.GetEnv("ORIGINS_ALLOWED", "http://localhost:8000"),
		JWTSecret:   initx.GetEnv("JWT_SECRET", ""),
		Services: []ServiceConfig{
			{
				Name:   "auth-service",
				URL:    initx.GetEnv("AUTH_SERVICE_URL", "http://auth-service:3000"),
				Prefix: "/auth",
			},
			{
				Name:   "image-service",
				URL:    initx.GetEnv("IMAGE_SERVICE_URL", "http://image-service:3000"),
				Prefix: "/images",
			},
		},
	}
}
