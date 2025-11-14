package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	DBHost        string
	DBUser        string
	DBName        string
	DBPassword    string
	AdminPassword string
	TLSCertFile   string
	TLSKeyFile    string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		DBHost:        os.Getenv("POSTGRES_HOST"),
		DBUser:        getEnvOrDefault("POSTGRES_USER", "challenge"),
		DBName:        getEnvOrDefault("POSTGRES_DB", "challenge"),
		DBPassword:    os.Getenv("POSTGRES_PASSWORD"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		TLSCertFile:   os.Getenv("TLS_CERT_FILE"),
		TLSKeyFile:    os.Getenv("TLS_KEY_FILE"),
	}

	// Validate required fields when using PostgreSQL
	if cfg.DBHost != "" && cfg.DBPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD is required when POSTGRES_HOST is set")
	}

	return cfg, nil
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
