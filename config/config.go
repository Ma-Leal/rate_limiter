package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr      string
	RateLimitIP    int
	RateLimitToken int
	Window         time.Duration
	BlockDuration  time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load("cmd/server/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RateLimitIP:    getEnvAsInt("RATE_LIMIT_IP", 5),
		RateLimitToken: getEnvAsInt("RATE_LIMIT_TOKEN", 10),
		Window:         time.Duration(getEnvAsInt("WINDOW_SECONDS", 60)) * time.Second,
		BlockDuration:  time.Duration(getEnvAsInt("BLOCK_DURATION_SECONDS", 300)) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := getEnv(name, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
