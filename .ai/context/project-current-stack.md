# project-current-stack

## purpose
- AI-only structured context for the current implemented project stack.
- Not an authoritative document.
- Source of truth remains `doc/project/current-stack.md` and actual files under `code/`.

## source_documents
- `doc/project/current-stack.md`
- `TECH_STACK.md`
- `code/go.mod`
- `code/docker-compose.yml`
- `code/Dockerfile`
- `code/frontend/Dockerfile`
- `code/openapi/openapi.yaml`

## current_stack_summary
- backend_language: Go
- backend_framework: Echo
- database: PostgreSQL 16
- frontend: static HTML/CSS/JS served by nginx
- api_contract: OpenAPI 3.0.3 static yaml
- orchestration: Docker Compose
- operator_entrypoint: root Makefile
- ci_entrypoint: `.github/workflows/ci.yml`

## current_services
- db:
  - image: `postgres:16-alpine`
  - role: persistent auth state storage
- backend:
  - runtime: compiled Go binary
  - role: login API, OpenAPI file, Swagger page, health check
- frontend:
  - runtime: nginx
  - role: login page and reverse proxy to backend `/api` and `/docs`

## current_ports
- frontend_public: 3000
- backend_public: 18080
- backend_internal: 8080
- db_internal: 5432

## auth_slice_status
- implemented_feature: login only
- token_model: access + refresh
- lock_policy: 5 failures -> 10 minute lock
- seed_user:
  - loginId: `user-01`
  - password: `pass1234`

## sync_note
- If runtime stack, framework choice, database, docker layout, ports, or API documentation strategy changes, this file must be re-reviewed.
