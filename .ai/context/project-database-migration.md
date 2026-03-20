# project-database-migration

## purpose
- AI-only structured context for database migration and rollback rules.
- Not an authoritative document.
- Source of truth remains `doc/project/database-migration.md`, `doc/project/current-stack.md`, and actual files under `code/`.

## source_documents
- `doc/project/database-migration.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `doc/README.md`
- `code/internal/db/migrate.go`
- `code/cmd/migrate/main.go`
- `code/migrations/`
- `Makefile`

## path_note
- Repository actual paths currently use `doc/`, `code/`, `.ai/` at root.
- Some governing docs describe these conceptually as `project/doc`, `project/code`, `project/.ai`.
- For current local file access, use actual root paths first.

## current_model
- database: PostgreSQL
- migration_runner: custom Go runner
- migration_tracking: `schema_migrations`
- startup_behavior: auto-apply up migrations only
- rollback_behavior: manual CLI-triggered rollback only
- up_file_format: `NNN_description.sql`
- down_file_format: `NNN_description.down.sql`

## coding_rules
- do not edit already-applied up migration files
- create up/down migration pairs in the same workstream
- keep one migration focused on one intent
- prefer corrective forward migrations over editing history
- sync docs, AI context, and AI logs when migration rules or commands change

## verification_checkpoints
- `go test ./...`
- `go build ./...`
- `go vet ./...`
- `make migrate-up`
- `make migrate-down STEPS=1`

## sync_note
- If migration naming, rollback policy, commands, or verification flow changes, this file must be re-reviewed together with source documents.
