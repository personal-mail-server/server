package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port               int
	DatabaseURL        string
	AccessTokenSecret  string
	RefreshTokenSecret string
	EnableCORS         bool
	AllowedOrigins     []string
	RequestTimeout     time.Duration
}

func Load() (Config, error) {
	port := envInt("PORT", 8080)
	databaseURL := LoadDatabaseURL()

	accessSecret := strings.TrimSpace(os.Getenv("ACCESS_TOKEN_SECRET"))
	if accessSecret == "" {
		return Config{}, fmt.Errorf("ACCESS_TOKEN_SECRET is required")
	}

	refreshSecret := strings.TrimSpace(os.Getenv("REFRESH_TOKEN_SECRET"))
	if refreshSecret == "" {
		return Config{}, fmt.Errorf("REFRESH_TOKEN_SECRET is required")
	}

	allowedOrigins := envCSV("ALLOWED_ORIGINS")

	return Config{
		Port:               port,
		DatabaseURL:        databaseURL,
		AccessTokenSecret:  accessSecret,
		RefreshTokenSecret: refreshSecret,
		EnableCORS:         len(allowedOrigins) > 0,
		AllowedOrigins:     allowedOrigins,
		RequestTimeout:     time.Duration(envInt("REQUEST_TIMEOUT_SECONDS", 10)) * time.Second,
	}, nil
}

func LoadDatabaseURL() string {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return "postgres://postgres:postgres@db:5432/mail_server?sslmode=disable"
	}
	return databaseURL
}

func envInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func envCSV(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
