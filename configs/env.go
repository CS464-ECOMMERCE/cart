package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	GRPCPort        string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	RedisDefaultTTL time.Duration
}

var envConfig EnvConfig

// InitEnv initializes environment variables
func InitEnv() {
	// Load .env file if it exists
	godotenv.Load("./secrets/.env")

	// Set default values
	envConfig = EnvConfig{
		GRPCPort:        getEnv("GRPC_PORT", "50051"),
		RedisAddr:       getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", "redis_password"),
		RedisDB:         getEnvAsInt("REDIS_DB", 0),
		RedisDefaultTTL: getEnvAsDuration("REDIS_DEFAULT_TTL", 30*time.Minute),
	}

	fmt.Println("Environment variables initialized")
}

// GetEnvConfig returns the current environment configuration
func GetEnvConfig() EnvConfig {
	return envConfig
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to get environment variable as int
func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// Helper function to get environment variable as duration
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
