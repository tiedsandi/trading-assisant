package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"trading-assistant/backend/db"
	"trading-assistant/backend/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "status") {
		return errors.New("usage: go run ./cmd/migrate [up|status]")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	pool, err := database.Open(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()
	provider, sqlDB, err := db.NewProvider(pool)
	if err != nil {
		return errors.New("migration initialization failed")
	}
	defer sqlDB.Close()
	if os.Args[1] == "up" {
		results, err := provider.Up(ctx)
		if err != nil {
			return errors.New("migration failed; inspect the schema and pending migration files")
		}
		fmt.Printf("Applied %d migration(s).\n", len(results))
		return nil
	}
	statuses, err := provider.Status(ctx)
	if err != nil {
		return errors.New("migration status failed")
	}
	for _, status := range statuses {
		fmt.Printf("%s %s\n", status.Source.Path, status.State)
	}
	return nil
}
