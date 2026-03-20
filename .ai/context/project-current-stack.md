# project-current-stack

## purpose
- AI-only structured context for the current implemented project stack.
- Not an authoritative document.
- Source of truth remains `doc/project/current-stack.md` and actual files under `code/`.
- Migration and rollback details are additionally governed by `doc/project/database-migration.md`.

## source_documents
- `doc/project/database-migration.md`
- `doc/project/current-stack.md`
- `TECH_STACK.md`
- `code/cmd/migrate/main.go`
- `code/go.mod`
- `code/internal/db/migrate.go`
- `code/docker-compose.yml`
- `code/Dockerfile`
- `code/frontend/Dockerfile`
- `code/openapi/openapi.yaml`

## current_stack_summary
- backend_language: Go
- backend_framework: Echo
- database: PostgreSQL 16
- database_migration: custom SQL up/down runner with manual rollback CLI
- frontend: static HTML/CSS/JS served by nginx
- api_contract: OpenAPI 3.0.3 static yaml
- orchestration: Docker Compose
- operator_entrypoint: root Makefile
- ci_entrypoint: `.github/workflows/ci.yml`
- pr_automation_entrypoint: `make push-trigger`

## current_services
- db:
  - image: `postgres:16-alpine`
  - role: persistent auth state storage
- backend:
  - runtime: compiled Go binary
  - role: login API, logout API, token reissue API, OpenAPI file, Swagger page, health check
- frontend:
  - runtime: nginx
  - role: auth page with login/logout interaction and reverse proxy to backend `/api` and `/docs`

## current_ports
- db_public: 5432
- frontend_public: 3000
- backend_public: 18080
- backend_internal: 8080
- db_internal: 5432

## auth_slice_status
- implemented_feature:
  - login
  - logout
  - token reissue
- token_model: access + refresh
- logout_invalidation: session versioning
- refresh_rotation: persisted refresh-token `jti` with one-time consume-and-replace
- lock_policy: 5 failures -> 10 minute lock
- seed_user:
  - loginId: `user-01`
  - password: `pass1234`

## mail_address_foundation_status
- implemented:
  - test mail address migration pair
  - test mail address model document
  - test mail address repository package
  - test mail address candidate generation API
  - test mail address create API
  - test mail address list/detail read APIs
- not_yet_implemented:
  - test mail address update/delete APIs
  - test mail address frontend UI

## sync_note
- If runtime stack, framework choice, database, docker layout, ports, or API documentation strategy changes, this file must be re-reviewed.
- If migration or rollback workflow changes, `project-database-migration.md` should be reviewed together.
