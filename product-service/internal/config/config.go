package config

import (
	"fmt"
	"time"

	"github.com/instrlabs/shared/functionx"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the product service
type Config struct {
	// Service
	ServiceName string `env:"SERVICE_NAME,required"`
	Port        string `env:"PORT,default=3005"`
	Environment string `env:"ENVIRONMENT,default=development"`

	// Database
	MongoURI     string `env:"MONGO_URI,required"`
	MongoDB      string `env:"MONGO_DB,required"`
	MongoTimeout int    `env:"MONGO_TIMEOUT,default=10"`

	// Security
	JWTSecret string `env:"JWT_SECRET,required"`

	// CORS
	Origins     string `env:"CORS_ORIGINS,default=http://localhost:3000"`
	CSRFEnabled bool   `env:"CSRF_ENABLED,default=true"`

	// Rate limiting
	RateLimit  int           `env:"RATE_LIMIT,default=100"`
	RateWindow time.Duration `env:"RATE_WINDOW,default=60s"`

	// Timeouts
	ReadTimeout  int `env:"READ_TIMEOUT,default=30"`
	WriteTimeout int `env:"WRITE_TIMEOUT,default=30"`
	IdleTimeout  int `env:"IDLE_TIMEOUT,default=60"`

	// External URLs
	ApiUrl string `env:"API_URL,required"`
	WebUrl string `env:"WEB_URL,required"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		ServiceName: functionx.GetEnvString("SERVICE_NAME", "product-service"),
		Port:        functionx.GetEnvString("PORT", "3005"),
		Environment: functionx.GetEnvString("ENVIRONMENT", "development"),

		MongoURI:     functionx.GetEnvString("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      functionx.GetEnvString("MONGO_DB", "instrlabs-apps"),
		MongoTimeout: functionx.GetEnvInt("MONGO_TIMEOUT", 10),

		JWTSecret: functionx.GetEnvString("JWT_SECRET", "your-super-secret-jwt-key"),

		Origins:     functionx.GetEnvString("CORS_ORIGINS", "http://localhost:3000"),
		CSRFEnabled: functionx.GetEnvBool("CSRF_ENABLED", true),

		RateLimit:  functionx.GetEnvInt("RATE_LIMIT", 100),
		RateWindow: time.Duration(functionx.GetEnvInt("RATE_WINDOW", 60)) * time.Second,

		ReadTimeout:  functionx.GetEnvInt("READ_TIMEOUT", 30),
		WriteTimeout: functionx.GetEnvInt("WRITE_TIMEOUT", 30),
		IdleTimeout:  functionx.GetEnvInt("IDLE_TIMEOUT", 60),

		ApiUrl: functionx.GetEnvString("API_URL", "http://localhost:3000"),
		WebUrl: functionx.GetEnvString("WEB_URL", "http://localhost:3000"),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("SERVICE_NAME is required")
	}

	if c.MongoURI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}

	if c.MongoDB == "" {
		return fmt.Errorf("MONGO_DB is required")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if c.ApiUrl == "" {
		return fmt.Errorf("API_URL is required")
	}

	if c.WebUrl == "" {
		return fmt.Errorf("WEB_URL is required")
	}

	// Validate timeout values
	if c.MongoTimeout <= 0 {
		return fmt.Errorf("MONGO_TIMEOUT must be greater than 0")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("READ_TIMEOUT must be greater than 0")
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("WRITE_TIMEOUT must be greater than 0")
	}

	if c.IdleTimeout <= 0 {
		return fmt.Errorf("IDLE_TIMEOUT must be greater than 0")
	}

	// Validate rate limiting
	if c.RateLimit <= 0 {
		return fmt.Errorf("RATE_LIMIT must be greater than 0")
	}

	return nil
}
