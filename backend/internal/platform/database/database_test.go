package database

import (
	"context"
	"testing"
)

func TestOpenRejectsInvalidConfiguration(t *testing.T) {
	for _, url := range []string{"", " ", "postgres://user:secret@localhost:invalid/db"} {
		pool, err := Open(context.Background(), url)
		if pool != nil {
			pool.Close()
			t.Fatal("unexpected pool")
		}
		if err == nil {
			t.Fatal("expected configuration error")
		}
		if err.Error() != "DATABASE_URL is required" && err.Error() != "invalid DATABASE_URL" {
			t.Fatal("expected sanitized configuration error")
		}
	}
}
