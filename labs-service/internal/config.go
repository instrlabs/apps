package internal

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	MongoURI    string
	MongoDB     string

	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool

	IMAGEServiceURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Environment: os.Getenv("ENVIRONMENT"),
		Port:        os.Getenv("PORT"),
		MongoURI:    os.Getenv("MONGO_URI"),
		MongoDB:     os.Getenv("MONGO_DB"),

		S3Endpoint:  os.Getenv("S3_ENDPOINT"),
		S3Region:    os.Getenv("S3_REGION"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3UseSSL:    getEnvBool("S3_USE_SSL", false),

		IMAGEServiceURL: os.Getenv("IMAGE_SERVICE_URL"),
	}
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if value == "true" || value == "1" || value == "yes" {
			return true
		} else if value == "false" || value == "0" || value == "no" {
			return false
		}
	}
	return fallback
}
