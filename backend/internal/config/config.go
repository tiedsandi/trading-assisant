package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Config struct {
	Environment string
	HTTPHost    string
	HTTPPort    string
}

func Load() (Config, error) {
	cfg := Config{
		Environment: envOrDefault("APP_ENV", "development"),
		HTTPHost:    envOrDefault("HTTP_HOST", "0.0.0.0"),
		HTTPPort:    envOrDefault("HTTP_PORT", "8080"),
	}
	port, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("HTTP_PORT must be an integer between 1 and 65535")
	}
	return cfg, nil
}

func (c Config) Address() string {
	return net.JoinHostPort(c.HTTPHost, c.HTTPPort)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
