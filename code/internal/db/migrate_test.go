package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverMigrationsIgnoresDownFilesAndSortsAscending(t *testing.T) {
	dir := t.TempDir()
	files := []string{
		"002_seed_login_user.sql",
		"001_create_users.sql",
		"001_create_users.down.sql",
		"005_create_refresh_tokens.down.sql",
		"005_create_refresh_tokens.sql",
		"badname.sql",
		"README.txt",
	}

	for _, name := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("SELECT 1;"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	migrations, err := discoverMigrations(dir)
	if err != nil {
		t.Fatalf("discover migrations: %v", err)
	}

	if len(migrations) != 3 {
		t.Fatalf("expected 3 up migrations, got %d", len(migrations))
	}

	wantVersions := []string{
		"001_create_users.sql",
		"002_seed_login_user.sql",
		"005_create_refresh_tokens.sql",
	}

	for i, want := range wantVersions {
		if migrations[i].version != want {
			t.Fatalf("migration[%d] version = %s, want %s", i, migrations[i].version, want)
		}
	}
}

func TestDiscoverMigrationsDerivesDownPathFromUpFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "004_add_users_session_version.sql"), []byte("SELECT 1;"), 0o644); err != nil {
		t.Fatalf("write migration: %v", err)
	}

	migrations, err := discoverMigrations(dir)
	if err != nil {
		t.Fatalf("discover migrations: %v", err)
	}

	if len(migrations) != 1 {
		t.Fatalf("expected 1 migration, got %d", len(migrations))
	}

	want := filepath.Join(dir, "004_add_users_session_version.down.sql")
	if migrations[0].downPath != want {
		t.Fatalf("down path = %s, want %s", migrations[0].downPath, want)
	}
}
