# project-current-testing

## purpose
- AI-only structured context for the current implemented testing state.
- Not an authoritative document.
- Source of truth remains `doc/project/current-testing.md` and actual test/code files.

## source_documents
- `doc/project/current-testing.md`
- `code/internal/auth/service_test.go`
- `code/internal/auth/validation_test.go`
- `code/internal/db/migrate_test.go`
- `code/internal/http/handlers/auth_handler_test.go`
- `code/README.md`

## current_test_summary
- unit_tests:
  - auth service tests
  - auth validation tests
  - auth handler tests
  - logout service and handler tests
  - token reissue service and handler tests
  - migration utility tests for discovery and down-path derivation
  - test mail address repository tests for create/read/list/update-delete behaviors
  - test mail address candidate generation service and handler tests
  - test mail address create service and handler tests
  - test mail address read service and handler tests
- execution_checks:
  - `go test ./...`
  - `go build ./...`
  - `go vet ./...`
  - `make migrate-up`
  - `make migrate-down STEPS=1`
  - `make up` stack startup check
  - `make status` stack status check
  - `make down` stack shutdown check
  - login page http check
  - swagger route check
  - openapi yaml route check
  - login success/failure/lock verification
  - logout success and token reuse rejection verification
  - github actions ci workflow verification
  - push-trigger unit test verification

## manual_verification_notes
- token reissue flow is currently tracked as a manual verification scenario, not a persisted CI or docker-smoke check
- automated coverage for token reissue currently exists at unit-test level in auth service and handler tests


## current_gaps
- no frontend automated tests
- no browser automation E2E tests
- no isolated database integration test suite
- no real Postgres migration up/down integration test suite
- no dedicated protected API middleware coverage yet

## mail_address_test_status
- repository_tests_exist: true
- package_path: `code/internal/testaddress/postgres_repository_test.go`
- coverage_focus:
  - duplicate email handling
  - active-owner listing
  - soft delete behavior
  - not-found mapping
- candidate_api_tests_exist: true
- candidate_api_paths:
  - `code/internal/testaddress/service_test.go`
  - `code/internal/http/handlers/testaddress_handler_test.go`
- create_api_tests_exist: true
- create_api_paths:
  - `code/internal/testaddress/service_test.go`
  - `code/internal/http/handlers/testaddress_handler_test.go`
- read_api_tests_exist: true
- read_api_paths:
  - `code/internal/testaddress/service_test.go`
  - `code/internal/http/handlers/testaddress_handler_test.go`

## ci_status
- workflow_exists: true
- workflow_path: `.github/workflows/ci.yml`
- current_jobs:
  - go-checks
  - docker-smoke

## push_trigger_status
- workflow_exists: false
- local_entrypoint: `make push-trigger`
- implementation_path: `code/cmd/push-trigger/main.go`

## sync_note
- If tests, verification commands, or validation flow changes, this file must be re-reviewed.
