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
	CSRFEnabled bool
	Services    []ServiceConfig
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":3000"),
		Origins:     initx.GetEnv("ORIGINS_ALLOWED", "http://localhost:8000"),
		JWTSecret:   initx.GetEnv("JWT_SECRET", ""),
		CSRFEnabled: initx.GetEnvBool("CSRF_ENABLED", true),
		Services: []ServiceConfig{
			{
				Name:   "auth-service",
				URL:    initx.GetEnv("AUTH_SERVICE", "http://auth-service:3000"),
				Prefix: "/auth",
			},
			{
				Name:   "image-service",
				URL:    initx.GetEnv("IMAGE_SERVICE", "http://image-service:3000"),
				Prefix: "/images",
			},
			{
				Name:   "pdf-service",
				URL:    initx.GetEnv("PDF_SERVICE", "http://pdf-service:3000"),
				Prefix: "/pdfs",
			},
			{
				Name:   "product-service",
				URL:    initx.GetEnv("PRODUCT_SERVICE", "http://product-service:3005"),
				Prefix: "/products",
			},
		},
	}
}
