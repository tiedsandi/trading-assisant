package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("HTTP_HOST", "")
	t.Setenv("HTTP_PORT", "")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environment != "development" || cfg.Address() != "0.0.0.0:8080" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("HTTP_HOST", "127.0.0.1")
	t.Setenv("HTTP_PORT", "9090")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environment != "test" || cfg.Address() != "127.0.0.1:9090" {
		t.Fatalf("environment overrides were not applied: %+v", cfg)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	for _, port := range []string{"0", "-1", "65536", "abc", "8080.5", " 8080 "} {
		t.Run(port, func(t *testing.T) {
			t.Setenv("HTTP_PORT", port)
			if _, err := Load(); err == nil {
				t.Fatalf("expected validation error for HTTP_PORT=%q", port)
			}
		})
	}
}

func TestLoadAcceptsPortBoundaries(t *testing.T) {
	for _, port := range []string{"1", "65535"} {
		t.Run(port, func(t *testing.T) {
			t.Setenv("HTTP_PORT", port)
			cfg, err := Load()
			if err != nil {
				t.Fatal(err)
			}
			if cfg.HTTPPort != port {
				t.Fatalf("HTTPPort = %q, want %q", cfg.HTTPPort, port)
			}
		})
	}
}

func TestAddressIPv6(t *testing.T) {
	cfg := Config{HTTPHost: "::1", HTTPPort: "8080"}
	if got := cfg.Address(); got != "[::1]:8080" {
		t.Fatalf("Address() = %q, want [::1]:8080", got)
	}
}
