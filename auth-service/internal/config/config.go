package config

import (
	"github.com/instrlabs/shared/functionx"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the auth service
type Config struct {
	Environment        string
	Port               string
	MongoURI           string
	MongoDB            string
	JWTSecret          string
	TokenExpiryHours   int
	RefreshExpiryHours int
	SMTPHost           string
	SMTPPort           string
	SMTPUsername       string
	SMTPPassword       string
	EmailFrom          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectUrl  string
	ApiUrl             string
	WebUrl             string
	PinEnabled         bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: functionx.GetEnvString("ENVIRONMENT", "development"),
		Port:        functionx.GetEnvString("PORT", "3001"),

		MongoURI: functionx.GetEnvString("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  functionx.GetEnvString("MONGO_DB", "auth_service"),

		JWTSecret:          functionx.GetEnvString("JWT_SECRET", "your-super-secret-jwt-key"),
		TokenExpiryHours:   functionx.GetEnvInt("TOKEN_EXPIRY_HOURS", 1),
		RefreshExpiryHours: functionx.GetEnvInt("REFRESH_EXPIRY_HOURS", 720),

		SMTPHost:     functionx.GetEnvString("SMTP_HOST", ""),
		SMTPPort:     functionx.GetEnvString("SMTP_PORT", ""),
		SMTPUsername: functionx.GetEnvString("SMTP_USERNAME", ""),
		SMTPPassword: functionx.GetEnvString("SMTP_PASSWORD", ""),
		EmailFrom:    functionx.GetEnvString("EMAIL_FROM", "noreply@example.com"),

		GoogleClientID:     functionx.GetEnvString("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: functionx.GetEnvString("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectUrl:  functionx.GetEnvString("GOOGLE_REDIRECT_URL", ""),

		ApiUrl: functionx.GetEnvString("API_URL", "http://localhost:3000"),
		WebUrl: functionx.GetEnvString("WEB_URL", "http://localhost:3000"),

		PinEnabled: functionx.GetEnvBool("PIN_ENABLED", false),
	}
}
