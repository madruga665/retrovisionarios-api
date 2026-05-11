package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	Port           string
	TrustedProxies string
	LogLevel       string
	AllowedOrigins string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		Port:           getEnv("PORT", "5000"),
		TrustedProxies: getEnv("TRUSTED_PROXIES", "localhost"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
