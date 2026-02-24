package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port        int
	DatabaseURL string
	LogLevel    string
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Port:        8080,
		DatabaseURL: "postgresql://chameleon:chameleon@localhost:5432/chameleon",
		LogLevel:    "info",
	}

	// Override with env vars
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return cfg
}
