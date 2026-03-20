package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const downMigrationSuffix = ".down.sql"

var upMigrationPattern = regexp.MustCompile(`^\d{3}_.+\.sql$`)

type migrationFile struct {
	version  string
	upPath   string
	downPath string
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	return MigrateUp(ctx, pool, migrationsDir)
}

func MigrateUp(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	if err := ensureSchemaMigrationsTable(ctx, pool); err != nil {
		return err
	}

	migrations, err := discoverMigrations(migrationsDir)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		alreadyApplied, err := migrationExists(ctx, pool, migration.version)
		if err != nil {
			return err
		}
		if alreadyApplied {
			continue
		}

		statement, err := os.ReadFile(migration.upPath)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", migration.version, err)
		}

		if err := applyMigration(ctx, pool, migration.version, string(statement)); err != nil {
			return err
		}
	}

	return nil
}

func RollbackSteps(ctx context.Context, pool *pgxpool.Pool, migrationsDir string, steps int) ([]string, error) {
	if steps < 1 {
		return nil, fmt.Errorf("steps must be at least 1")
	}

	if err := ensureSchemaMigrationsTable(ctx, pool); err != nil {
		return nil, err
	}

	migrations, err := discoverMigrations(migrationsDir)
	if err != nil {
		return nil, err
	}

	byVersion := make(map[string]migrationFile, len(migrations))
	for _, migration := range migrations {
		byVersion[migration.version] = migration
	}

	rolledBack := make([]string, 0, steps)
	for range steps {
		version, err := lastAppliedMigration(ctx, pool)
		if err != nil {
			return rolledBack, err
		}
		if version == "" {
			break
		}

		migration, ok := byVersion[version]
		if !ok {
			return rolledBack, fmt.Errorf("rollback migration %s: up migration file not found", version)
		}

		statement, err := os.ReadFile(migration.downPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return rolledBack, fmt.Errorf("rollback migration %s: down migration file not found", version)
			}
			return rolledBack, fmt.Errorf("read rollback file %s: %w", filepath.Base(migration.downPath), err)
		}

		if err := rollbackMigration(ctx, pool, version, string(statement)); err != nil {
			return rolledBack, err
		}

		rolledBack = append(rolledBack, version)
	}

	return rolledBack, nil
}

func ensureSchemaMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}
	return nil
}

func discoverMigrations(migrationsDir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	migrations := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") || strings.HasSuffix(name, downMigrationSuffix) || !upMigrationPattern.MatchString(name) {
			continue
		}

		migrations = append(migrations, migrationFile{
			version:  name,
			upPath:   filepath.Join(migrationsDir, name),
			downPath: filepath.Join(migrationsDir, strings.TrimSuffix(name, ".sql")+downMigrationSuffix),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return migrations, nil
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, version string, statement string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration tx %s: %w", version, err)
	}

	if _, err := tx.Exec(ctx, statement); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("execute migration %s: %w", version, err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, version); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("record migration %s: %w", version, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", version, err)
	}

	return nil
}

func rollbackMigration(ctx context.Context, pool *pgxpool.Pool, version string, statement string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin rollback tx %s: %w", version, err)
	}

	if _, err := tx.Exec(ctx, statement); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("execute rollback %s: %w", version, err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, version); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("delete migration record %s: %w", version, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit rollback %s: %w", version, err)
	}

	return nil
}

func migrationExists(ctx context.Context, pool *pgxpool.Pool, version string) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}
	return exists, nil
}

func lastAppliedMigration(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	var version string
	err := pool.QueryRow(ctx, `
		SELECT version
		FROM schema_migrations
		ORDER BY version DESC
		LIMIT 1
	`).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("select last applied migration: %w", err)
	}
	return version, nil
}
