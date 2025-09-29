package internal

import (
	initx "github.com/instrlabs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment        string
	Port               string
	MongoURI           string
	MongoDB            string
	JWTSecret          string
	CookieDomain       string
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
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", ""),
		Port:        initx.GetEnv("PORT", ""),

		MongoURI: initx.GetEnv("MONGO_URI", ""),
		MongoDB:  initx.GetEnv("MONGO_DB", ""),

		JWTSecret:          initx.GetEnv("JWT_SECRET", ""),
		CookieDomain:       initx.GetEnv("COOKIE_DOMAIN", ""),
		TokenExpiryHours:   initx.GetEnvInt("TOKEN_EXPIRY_HOURS", 1),
		RefreshExpiryHours: initx.GetEnvInt("REFRESH_EXPIRY_HOURS", 720),

		SMTPHost:     initx.GetEnv("SMTP_HOST", ""),
		SMTPPort:     initx.GetEnv("SMTP_PORT", ""),
		SMTPUsername: initx.GetEnv("SMTP_USERNAME", ""),
		SMTPPassword: initx.GetEnv("SMTP_PASSWORD", ""),
		EmailFrom:    initx.GetEnv("EMAIL_FROM", ""),

		GoogleClientID:     initx.GetEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: initx.GetEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectUrl:  initx.GetEnv("GOOGLE_REDIRECT_URL", ""),

		ApiUrl: initx.GetEnv("API_URL", ""),
		WebUrl: initx.GetEnv("WEB_URL", ""),
	}
}
