// Package config читает конфигурацию приложения из переменных окружения.
// Принцип 12-factor app: конфигурация хранится в окружении, а не в коде.
package config

import (
	"fmt"
	"os"
)

// Config хранит всю конфигурацию приложения.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Load считывает конфигурацию. Возвращает ошибку если обязательное поле не задано.
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return &Config{
		Port:      getEnvOrDefault("PORT", "8080"),
		DatabaseURL: dbURL,
		// JWT_SECRET должен быть заменён в production — пустой дефолт специально
		// оставлен очевидно небезопасным, чтобы не уйти в прод незамеченным.
		JWTSecret: getEnvOrDefault("JWT_SECRET", "default-secret-change-me"),
	}, nil
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
