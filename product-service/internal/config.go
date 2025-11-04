package internal

import (
	initx "github.com/instrlabs/shared/init"
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
		Environment: initx.GetEnv("ENVIRONMENT", "development"),
		Port:        initx.GetEnv("PORT", "3005"),

		MongoURI: initx.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  initx.GetEnv("MONGO_DB", "instrlabs"),

		ApiUrl: initx.GetEnv("API_URL", "http://localhost:3000"),
	}
}
