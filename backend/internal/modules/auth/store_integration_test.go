//go:build integration

package auth

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"trading-assistant/backend/db"
	"trading-assistant/backend/internal/platform/database"
)

func TestPostgresRegistrationAtomicityAndConcurrentDuplicate(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Fatal("TEST_DATABASE_URL is required; use compose.test.yml")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	provider, sqlDB, err := db.NewProvider(pool)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	if _, err := provider.Up(ctx); err != nil {
		t.Fatal(err)
	}
	repository := postgresStore{pool: pool}
	email := fmt.Sprintf("atomic-%d@example.com", time.Now().UnixNano())
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _, err := pool.Exec(cleanupCtx, `DELETE FROM users WHERE email=$1`, email); err != nil {
			t.Error("test user cleanup failed")
		}
	}()
	// A session constraint failure must roll back the user inserted earlier in the transaction.
	if _, err := repository.register(ctx, email, "hash-fixture-not-a-password", make([]byte, 31), time.Now().Add(SessionLifetime)); err == nil {
		t.Fatal("invalid session digest should fail")
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM users WHERE email=$1`, email).Scan(&count); err != nil || count != 0 {
		t.Fatal("failed registration left a user without a session")
	}
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, digest, err := newSessionToken()
			if err != nil {
				results <- err
				return
			}
			_, err = repository.register(ctx, email, "hash-fixture-not-a-password", digest, time.Now().Add(SessionLifetime))
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	var successful, duplicate int
	for err := range results {
		switch err {
		case nil:
			successful++
		case errDuplicate:
			duplicate++
		default:
			t.Fatalf("unexpected registration result: %v", err)
		}
	}
	if successful != 1 || duplicate != 1 {
		t.Fatalf("concurrent registration results success=%d duplicate=%d", successful, duplicate)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM sessions s JOIN users u ON s.user_id=u.id WHERE u.email=$1`, email).Scan(&count); err != nil || count != 1 {
		t.Fatal("concurrent duplicate created an extra or missing session")
	}
}
