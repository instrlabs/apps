package internal

import (
	"github.com/instrlabs/shared/functionx"
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
		Environment: functionx.GetEnvString("ENVIRONMENT", "development"),
		Port:        functionx.GetEnvString("PORT", ":3000"),
		Origins:     functionx.GetEnvString("ORIGINS_ALLOWED", "http://localhost:8000"),
		JWTSecret:   functionx.GetEnvString("JWT_SECRET", ""),
		CSRFEnabled: functionx.GetEnvBool("CSRF_ENABLED", true),
		Services: []ServiceConfig{
			{
				Name:   "auth-service",
				URL:    functionx.GetEnvString("AUTH_SERVICE", "http://auth-service:3000"),
				Prefix: "/auth",
			},
			{
				Name:   "image-service",
				URL:    functionx.GetEnvString("IMAGE_SERVICE", "http://image-service:3000"),
				Prefix: "/images",
			},
			{
				Name:   "pdf-service",
				URL:    functionx.GetEnvString("PDF_SERVICE", "http://pdf-service:3000"),
				Prefix: "/pdfs",
			},
			{
				Name:   "product-service",
				URL:    functionx.GetEnvString("PRODUCT_SERVICE", "http://product-service:3005"),
				Prefix: "/products",
			},
		},
	}
}
