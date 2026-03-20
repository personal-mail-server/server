package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"personal-mail-server/internal/config"
	"personal-mail-server/internal/db"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up or down")
	migrationsDir := flag.String("dir", "migrations", "migration directory path")
	steps := flag.Int("steps", 1, "number of rollback steps for down direction")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, config.LoadDatabaseURL())
	if err != nil {
		fmt.Fprintf(os.Stderr, "migration connection failed: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	switch *direction {
	case "up":
		if err := db.MigrateUp(ctx, pool, *migrationsDir); err != nil {
			fmt.Fprintf(os.Stderr, "migration up failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("migration up completed")
	case "down":
		rolledBack, err := db.RollbackSteps(ctx, pool, *migrationsDir, *steps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "migration down failed: %v\n", err)
			os.Exit(1)
		}
		if len(rolledBack) == 0 {
			fmt.Println("no migrations rolled back")
			return
		}
		fmt.Printf("rolled back %d migration(s)\n", len(rolledBack))
		for _, version := range rolledBack {
			fmt.Printf("- %s\n", version)
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid direction %q\n", *direction)
		os.Exit(1)
	}
}
