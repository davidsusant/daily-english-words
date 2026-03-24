package config

import (
	"fmt"
	"os"
)

// Config to holds all application configuration
type Config struct {
	DBHost			string
	DBPort			string
	DBUser			string
	DBPassword		string
	DBName			string
	ServerPort		string
	GeminiAPIKey	string
}

// Load reads configuration from environment variables
// Returns an error if any required variable is missing
func Load() (*Config, error) {
	cfg := &Config{
		DBHost: 	getEnv("DB_HOST", "localhost"),
		DBPort: 	getEnv("DB_PORT", "5432"),
		DBUser: 	getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName: 	getEnv("DB_NAME", "daily_english_words"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),		
	}

	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	return cfg, nil
}

// Reads an env var or returns a fallback default
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}