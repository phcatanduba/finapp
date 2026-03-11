package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	DatabaseURL string

	JWTSecret           string
	JWTAccessTTLMinutes int
	JWTRefreshTTLDays   int

	PluggyClientID      string
	PluggyClientSecret  string
	PluggyBaseURL       string
	PluggyWebhookSecret string

	CORSAllowedOrigins string
}

func Load() (*Config, error) {
	// load .env if present (ignored in production if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		Env:                 getEnv("ENV", "development"),
		DatabaseURL:         mustGetEnv("DATABASE_URL"),
		JWTSecret:           mustGetEnv("JWT_SECRET"),
		JWTAccessTTLMinutes: getEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
		JWTRefreshTTLDays:   getEnvInt("JWT_REFRESH_TTL_DAYS", 7),
		PluggyClientID:      mustGetEnv("PLUGGY_CLIENT_ID"),
		PluggyClientSecret:  mustGetEnv("PLUGGY_CLIENT_SECRET"),
		PluggyBaseURL:       getEnv("PLUGGY_BASE_URL", "https://api.pluggy.ai"),
		PluggyWebhookSecret: getEnv("PLUGGY_WEBHOOK_SECRET", ""),
		CORSAllowedOrigins:  getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}

	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
