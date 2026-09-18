package config

import (
	"strings"
	"testing"
)

func TestAuthOrigins(t *testing.T) {
	for _, tc := range []struct {
		name, environment, origins string
		valid                      bool
	}{
		{"development default", "development", "", true},
		{"explicit development", "development", "http://localhost:3000, http://127.0.0.1:3000", true},
		{"production missing", "production", "", false},
		{"production HTTPS", "production", "https://app.example.com", true},
		{"production insecure", "production", "http://app.example.com", false},
		{"production mixed", "production", "https://app.example.com,http://localhost:3000", false},
		{"wildcard", "development", "*", false},
		{"subdomain wildcard", "development", "https://*.example.com", false},
		{"credentials", "development", "https://user:secret@example.com", false},
		{"path", "development", "https://example.com/login", false},
		{"trailing slash", "development", "https://example.com/", false},
		{"query", "development", "https://example.com?anything", false},
		{"fragment", "development", "https://example.com#fragment", false},
		{"empty item", "development", "https://example.com,", false},
		{"unknown environment", "prod", "https://example.com", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tc.environment)
			t.Setenv("HTTP_PORT", "8080")
			t.Setenv("AUTH_ALLOWED_ORIGINS", tc.origins)
			cfg, err := Load()
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v error=%v", tc.valid, err)
			}
			if err != nil && strings.Contains(err.Error(), "secret") {
				t.Fatal("configuration error exposes credentials")
			}
			if err == nil && len(cfg.AuthAllowedOrigins) == 0 {
				t.Fatal("no allowed origins configured")
			}
		})
	}
}
