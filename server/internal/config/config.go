package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config holds all runtime configuration loaded from environment variables.
// The server performs a hard exit on startup if any required field is missing.
type Config struct {
	ServerPort       string
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	RedisURL         string
	JWTSecret        string
	AgentTokenSecret string
	RetentionDays    int
}

// MustLoad loads config from the environment and exits on validation failure.
func MustLoad() *Config {
	cfg := &Config{
		ServerPort:       getEnvOrDefault("SERVER_PORT", "8080"),
		PostgresHost:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnvOrDefault("POSTGRES_PORT", "5432"),
		PostgresDB:       getEnvOrDefault("POSTGRES_DB", "netwatch"),
		PostgresUser:     getEnvOrDefault("POSTGRES_USER", "netwatch"),
		PostgresPassword: requireEnv("POSTGRES_PASSWORD"),
		RedisURL:         getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:        requireEnv("JWT_SECRET"),
		AgentTokenSecret: requireEnv("AGENT_TOKEN_SECRET"),
		RetentionDays:    getEnvInt("RETENTION_DAYS", 30),
	}
	return cfg
}

// PostgresDSN builds the pgx connection string.
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PostgresUser, c.PostgresPassword,
	)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("config: required environment variable %q is not set", key)
	}
	return v
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("config: %q must be an integer, got %q", key, v)
		}
		return n
	}
	return fallback
}
