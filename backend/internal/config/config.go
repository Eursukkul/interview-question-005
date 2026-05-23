package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	CORSAllowedOrigins []string
}

func Load() Config {
	_ = godotenv.Load()

	port := getEnv("APP_PORT", "8080")
	origins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173")

	return Config{
		Port:               port,
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://queue_user:queue_password@localhost:5432/queue_db?sslmode=disable"),
		CORSAllowedOrigins: splitCSV(origins),
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
