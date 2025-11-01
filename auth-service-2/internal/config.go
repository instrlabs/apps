package internal

import (
	initx "github.com/instrlabs/shared/init"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string

	// Database
	MongoURI string
	MongoDB  string

	// JWT Configuration
	JWTSecret          string
	AccessTokenExpiry  int // hours
	RefreshTokenExpiry int // hours

	// OAuth - Google
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// PIN Authentication
	PinLength      int
	PinExpiryMins  int
	PinEnabled     bool

	// Email/SMTP
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	EmailFrom    string

	// URLs
	APIBaseURL string
	WebURL     string

	// Cookie
	CookieDomain string
	CookieSecure bool
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", ":8002"),

		MongoURI: initx.GetEnv("MONGO_URI", ""),
		MongoDB:  initx.GetEnv("MONGO_DB", ""),

		JWTSecret:          initx.GetEnv("JWT_SECRET", ""),
		AccessTokenExpiry:  initx.GetEnvInt("ACCESS_TOKEN_EXPIRY_HOURS", 1),
		RefreshTokenExpiry: initx.GetEnvInt("REFRESH_TOKEN_EXPIRY_HOURS", 720),

		GoogleClientID:     initx.GetEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: initx.GetEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  initx.GetEnv("GOOGLE_REDIRECT_URL", ""),

		PinLength:     initx.GetEnvInt("PIN_LENGTH", 6),
		PinExpiryMins: initx.GetEnvInt("PIN_EXPIRY_MINS", 15),
		PinEnabled:    initx.GetEnvBool("PIN_ENABLED", true),

		SMTPHost:     initx.GetEnv("SMTP_HOST", ""),
		SMTPPort:     initx.GetEnv("SMTP_PORT", "587"),
		SMTPUsername: initx.GetEnv("SMTP_USERNAME", ""),
		SMTPPassword: initx.GetEnv("SMTP_PASSWORD", ""),
		EmailFrom:    initx.GetEnv("EMAIL_FROM", ""),

		APIBaseURL: initx.GetEnv("API_BASE_URL", "http://localhost:8000"),
		WebURL:     initx.GetEnv("WEB_URL", "http://localhost:3000"),

		CookieDomain: initx.GetEnv("COOKIE_DOMAIN", "localhost"),
		CookieSecure: initx.GetEnvBool("COOKIE_SECURE", false),
	}
}
