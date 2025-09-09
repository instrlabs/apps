package initx

import (
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		switch v {
		case "true", "1", "yes", "TRUE", "Y", "y":
			return true
		case "false", "0", "no", "FALSE", "N", "n":
			return false
		}
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}

	return fallback
}
