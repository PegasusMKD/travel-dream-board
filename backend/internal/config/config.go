package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port    string
	GinMode string

	DatabaseURL          string
	DatabaseMaxConns     int
	DatabaseMaxIdleConns int
	DatabaseConnLifetime time.Duration

	LogLevel       string
	LogLevelStdout string
	LogLevelFile   string
	LogFilePath    string
	LogFilePrefix  string

	FrontendDir string
	UploadsDir  string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	JWTSecret          string
	OpenRouterAPIKey   string
	ScrapingAntAPIKey  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DatabaseURL:          getEnv("DATABASE_URL", ""),
		DatabaseMaxConns:     getEnvInt("DATABASE_MAX_CONNS", 25),
		DatabaseMaxIdleConns: getEnvInt("DATABASE_MAX_IDLE_CONNS", 10),
		DatabaseConnLifetime: getEnvDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),

		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogLevelStdout: getEnv("LOG_LEVEL_STDOUT", "info"),
		LogLevelFile:   getEnv("LOG_LEVEL_FILE", "debug"),
		LogFilePath:    getEnv("LOG_FILE_PATH", "./logs"),
		LogFilePrefix:  getEnv("LOG_FILE_PREFIX", "travel-dream-board"),

		FrontendDir: getEnv("FRONTEND_DIR", "./frontend/dist"),
		UploadsDir:  getEnv("UPLOADS_DIR", "./uploads"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-default-key"),
		OpenRouterAPIKey:   getEnv("OPENROUTER_API_KEY", ""),
		ScrapingAntAPIKey:  getEnv("SCRAPINGANT_API_KEY", ""),
	}

	// Validate required fields based on service
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch strings.ToLower(value) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return defaultValue
}

func (c *Config) validate() error {
	// Service-specific validation
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
