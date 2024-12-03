package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	RedisURL      string
	RabbitMQURL   string
	AuthServiceURL string
}

func Load() *Config {
	return &Config{
		ServerAddress:  getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://testuser:testpass@postgres:5432/testingdb?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://redis:6379"),
		RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://auth:9000"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
