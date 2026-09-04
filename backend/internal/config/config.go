package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment    string
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AllowedOrigins string

	RateLimitRPS              int
	RateLimitBurst            int
	AIRequestsPerMinute       int

	LLMAPIURL  string
	LLMAPIKey  string
	LLMModel   string

	ResendAPIKey string
	EmailFrom    string
	AppBaseURL   string

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
}

func Load() Config {
	loadDotEnv()

	return Config{
		Environment:         getEnv("APP_ENV", "development"),
		Port:                getEnv("PORT", "8080"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/heritage_weaver?sslmode=disable"),
		JWTSecret:           getEnv("JWT_SECRET", "dev-secret-do-not-use-in-production"),
		AllowedOrigins:      getEnv("CORS_ALLOWED_ORIGINS", "*"),
		RateLimitRPS:        getEnvInt("RATE_LIMIT_RPS", 10),
		RateLimitBurst:      getEnvInt("RATE_LIMIT_BURST", 20),
		AIRequestsPerMinute: getEnvInt("AI_RATE_LIMIT_PER_MINUTE", 10),

		LLMAPIURL: getEnv("LLM_API_URL", "https://opencode.ai/zen/v1/chat/completions"),
		LLMAPIKey: os.Getenv("LLM_API_KEY"),
		LLMModel:     getEnv("LLM_MODEL", "laguna-s-2.1-free"),
		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		EmailFrom:    getEnv("EMAIL_FROM", "Hearthside <onboarding@resend.dev>"),
		AppBaseURL:   getEnv("APP_BASE_URL", "http://localhost:3000"),
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPass:     getEnv("SMTP_PASS", ""),
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

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return val
}
