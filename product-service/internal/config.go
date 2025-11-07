package internal

import (
	"github.com/instrlabs/shared/functionx"
	"github.com/joho/godotenv"
)

type Config struct {
	// Environment
	Environment string
	Port        string

	// Database
	MongoURI string
	MongoDB  string

	// API
	ApiUrl string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: functionx.GetEnvString("ENVIRONMENT", "development"),
		Port:        functionx.GetEnvString("PORT", ":3005"),

		MongoURI: functionx.GetEnvString("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  functionx.GetEnvString("MONGO_DB", "instrlabs-apps"),

		ApiUrl: functionx.GetEnvString("API_URL", ""),
	}
}
