package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment    string
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AllowedOrigins string
}

func Load() Config {
	loadDotEnv()

	return Config{
		Environment:    getEnv("APP_ENV", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/heritage_weaver?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-do-not-use-in-production"),
		AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}
}

func loadDotEnv() {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("load .env: %v", err)
		}
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
