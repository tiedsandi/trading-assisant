//go:build integration

package integration

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"trading-assistant/backend/internal/app"
	"trading-assistant/backend/internal/platform/database"
)

func TestDatabaseReadiness(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Fatal("TEST_DATABASE_URL is required; use compose.test.yml")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	var value int
	if err := pool.QueryRow(ctx, "SELECT 1").Scan(&value); err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("SELECT 1 returned %d", value)
	}
	handler := app.NewHandler(pool.Ping)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/ready", nil))
	if response.Code != 200 {
		t.Fatalf("connected readiness = %d", response.Code)
	}
	pool.Close()
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/ready", nil))
	if response.Code != 503 {
		t.Fatalf("closed pool readiness = %d", response.Code)
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/health", nil))
	if response.Code != 200 {
		t.Fatalf("liveness must not depend on database: %d", response.Code)
	}
}
