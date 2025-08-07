package constants

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment           string
	Port                  string
	Hostname              string
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
	GoogleRedirectPath    string
	FEResetPassword       string
	FEOAuthRedirect       string
}

func NewConfig() *Config {
	godotenv.Load()

	tokenExpiryHours := 1
	expiryStr := os.Getenv("TOKEN_EXPIRY_HOURS")
	tokenExpiryHours, _ = strconv.Atoi(expiryStr)

	resetTokenExpiryHours := 24
	resetExpiryStr := os.Getenv("RESET_TOKEN_EXPIRY_HOURS")
	resetTokenExpiryHours, _ = strconv.Atoi(resetExpiryStr)

	return &Config{
		Environment:           os.Getenv("ENVIRONMENT"),
		Port:                  os.Getenv("PORT"),
		Hostname:              os.Getenv("HOSTNAME"),
		MongoURI:              os.Getenv("MONGO_URI"),
		MongoDB:               os.Getenv("MONGO_DB"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		TokenExpiryHours:      tokenExpiryHours,
		SMTPHost:              os.Getenv("SMTP_HOST"),
		SMTPPort:              os.Getenv("SMTP_PORT"),
		SMTPUsername:          os.Getenv("SMTP_USERNAME"),
		SMTPPassword:          os.Getenv("SMTP_PASSWORD"),
		EmailFrom:             os.Getenv("EMAIL_FROM"),
		ResetTokenExpiryHours: resetTokenExpiryHours,
		GoogleClientID:        os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectPath:    os.Getenv("GOOGLE_REDIRECT_PATH"),
		FEResetPassword:       os.Getenv("FE_RESET_PASSWORD"),
		FEOAuthRedirect:       os.Getenv("FE_OAUTH_REDIRECT"),
	}
}
