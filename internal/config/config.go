package config

import (
	"os"
	"strconv"
)

type Config struct {
	MeilisearchURL  string
	MeilisearchAPIKey string
	ListenAddr      string
	LogLevel        string
}

func Load() *Config {
	return &Config{
		MeilisearchURL:   getEnv("MEILISEARCH_URL", "http://localhost:7700"),
		MeilisearchAPIKey: getEnv("MEILISEARCH_API_KEY", ""),
		ListenAddr:       getEnv("LISTEN_ADDR", ":9200"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
