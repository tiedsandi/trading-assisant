package app

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReadiness(t *testing.T) {
	for _, tc := range []struct {
		name  string
		check func(context.Context) error
		code  int
		body  string
	}{
		{"healthy", func(ctx context.Context) error {
			deadline, ok := ctx.Deadline()
			if !ok || time.Until(deadline) > 2*time.Second {
				t.Error("database check must have a bounded deadline")
			}
			return nil
		}, 200, "{\"status\":\"ready\"}\n"},
		{"unavailable", func(context.Context) error { return errors.New("secret connection details") }, 503, "{\"status\":\"not_ready\"}\n"},
		{"missing checker", nil, 503, "{\"status\":\"not_ready\"}\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			NewHandler(tc.check).ServeHTTP(response, httptest.NewRequest("GET", "/ready", nil))
			if response.Code != tc.code || response.Body.String() != tc.body {
				t.Fatalf("got %d %q, want %d %q", response.Code, response.Body.String(), tc.code, tc.body)
			}
			if response.Header().Get("Content-Type") != "application/json" || response.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("missing readiness response headers")
			}
		})
	}
}

func TestReadinessPropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	response := httptest.NewRecorder()
	NewHandler(func(ctx context.Context) error {
		if ctx.Err() != context.Canceled {
			t.Error("request cancellation was not propagated")
		}
		return ctx.Err()
	}).ServeHTTP(response, httptest.NewRequest("GET", "/ready", nil).WithContext(ctx))
	if response.Code != 503 {
		t.Fatalf("status = %d, want 503", response.Code)
	}
}
