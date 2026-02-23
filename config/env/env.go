package env

import (
	"log"
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
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

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
