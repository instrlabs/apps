package internal

import (
	initx "github.com/histweety-labs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment           string
	Port                  string
	MongoURI              string
	MongoDB               string
	JWTSecret             string
	TokenExpiryHours      int
	SMTPHost              string
	SMTPPort              string
	SMTPUsername          string
	SMTPPassword          string
	EmailFrom             string
	ResetTokenExpiryHours int
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleRedirectUrl     string
	FEResetPassword       string
	FEOAuthRedirect       string
	CORSAllowedOrigins    string
	CookieDomain          string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment:           initx.GetEnv("ENVIRONMENT", ""),
		Port:                  initx.GetEnv("PORT", ""),
		MongoURI:              initx.GetEnv("MONGO_URI", ""),
		MongoDB:               initx.GetEnv("MONGO_DB", ""),
		JWTSecret:             initx.GetEnv("JWT_SECRET", ""),
		TokenExpiryHours:      initx.GetEnvInt("TOKEN_EXPIRY_HOURS", 1),
		SMTPHost:              initx.GetEnv("SMTP_HOST", ""),
		SMTPPort:              initx.GetEnv("SMTP_PORT", ""),
		SMTPUsername:          initx.GetEnv("SMTP_USERNAME", ""),
		SMTPPassword:          initx.GetEnv("SMTP_PASSWORD", ""),
		EmailFrom:             initx.GetEnv("EMAIL_FROM", ""),
		ResetTokenExpiryHours: initx.GetEnvInt("RESET_TOKEN_EXPIRY_HOURS", 24),
		GoogleClientID:        initx.GetEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:    initx.GetEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectUrl:     initx.GetEnv("GOOGLE_REDIRECT_URL", ""),
		FEResetPassword:       initx.GetEnv("FE_RESET_PASSWORD", ""),
		FEOAuthRedirect:       initx.GetEnv("FE_OAUTH_REDIRECT", ""),
		CORSAllowedOrigins:    initx.GetEnv("CORS_ALLOWED_ORIGINS", "http://web.localhost"),
		CookieDomain:          initx.GetEnv("COOKIE_DOMAIN", ".localhost"),
	}
}
