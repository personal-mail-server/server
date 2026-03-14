# project-current-testing

## purpose
- AI-only structured context for the current implemented testing state.
- Not an authoritative document.
- Source of truth remains `doc/project/current-testing.md` and actual test/code files.

## source_documents
- `doc/project/current-testing.md`
- `code/internal/auth/service_test.go`
- `code/internal/auth/validation_test.go`
- `code/internal/http/handlers/auth_handler_test.go`
- `code/README.md`

## current_test_summary
- unit_tests:
  - auth service tests
  - auth validation tests
  - auth handler tests
- execution_checks:
  - `go test ./...`
  - `go build ./cmd/server`
  - `go vet ./...`
  - `make up` stack startup check
  - `make status` stack status check
  - `make down` stack shutdown check
  - login page http check
  - swagger route check
  - login success/failure/lock verification
  - github actions ci workflow verification


## current_gaps
- no frontend automated tests
- no browser automation E2E tests
- no isolated database integration test suite
- no refresh-token flow tests

## ci_status
- workflow_exists: true
- workflow_path: `.github/workflows/ci.yml`
- current_jobs:
  - go-checks
  - docker-smoke

## sync_note
- If tests, verification commands, or validation flow changes, this file must be re-reviewed.
