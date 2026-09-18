package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Environment        string
	HTTPHost           string
	HTTPPort           string
	DatabaseURL        string
	AuthAllowedOrigins []string
}

func Load() (Config, error) {
	cfg := Config{
		Environment: envOrDefault("APP_ENV", "development"),
		HTTPHost:    envOrDefault("HTTP_HOST", "0.0.0.0"),
		HTTPPort:    envOrDefault("HTTP_PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	port, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("HTTP_PORT must be an integer between 1 and 65535")
	}
	if cfg.Environment != "development" && cfg.Environment != "test" && cfg.Environment != "production" {
		return Config{}, fmt.Errorf("APP_ENV must be development, test, or production")
	}
	origins := os.Getenv("AUTH_ALLOWED_ORIGINS")
	if origins == "" && cfg.Environment != "production" {
		origins = "http://localhost:3000"
	}
	if strings.TrimSpace(origins) == "" {
		return Config{}, fmt.Errorf("AUTH_ALLOWED_ORIGINS is required in production")
	}
	for _, value := range strings.Split(origins, ",") {
		origin := strings.TrimSpace(value)
		parsed, err := url.Parse(origin)
		if err != nil || parsed.Hostname() == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || parsed.Opaque != "" ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") || strings.Contains(parsed.Host, "*") ||
			(cfg.Environment == "production" && parsed.Scheme != "https") {
			return Config{}, fmt.Errorf("AUTH_ALLOWED_ORIGINS must contain exact HTTP(S) origins; production requires HTTPS")
		}
		cfg.AuthAllowedOrigins = append(cfg.AuthAllowedOrigins, origin)
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
