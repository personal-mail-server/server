# auth-login

## purpose
- AI-only structured context for login-related source documents.
- Not an authoritative document.
- Source of truth remains `doc/auth/login.md` and `doc/auth/login-api.md`.

## source_documents
- `doc/auth/login.md`
- `doc/auth/login-api.md`
- `doc/auth/login-frontend.md`
- `TRIGGERS.md`
- `AGENTS.md`
- `TECH_STACK.md`

## path_note
- Repository actual paths currently use `doc/`, `code/`, `.ai/` at root.
- Some governing docs describe these conceptually as `project/doc`, `project/code`, `project/.ai`.
- For current local file access, use actual root paths first.

## document_status
- login feature design exists: `doc/auth/login.md`
- login API spec exists: `doc/auth/login-api.md`
- backend login slice exists under `code/`
- frontend login slice exists under `code/frontend/`
- docker runnable slice exists under `code/docker-compose.yml` and `code/Dockerfile`

## implementation_status
- backend stack: Go + Echo
- frontend stack: static HTML/CSS/JS via nginx container
- database: PostgreSQL via Docker Compose
- swagger/openapi: static OpenAPI served at `/docs/openapi.yaml`, proxied Swagger page at `/docs`
- login page: available at frontend root with a single-card login form
- seed user:
  - loginId: `user-01`
  - password: `pass1234`

## frontend_document_status
- frontend behavior doc exists: `doc/auth/login-frontend.md`
- current frontend scope:
  - login form only
  - client-side validation
  - status panel with raw JSON response output
  - swagger link exposure
  - no redirect after login
  - no token persistence UI

## login_feature_scope

### included
- login ID and password authentication
- login success and failure handling
- failed-attempt counting
- account lock policy
- access token issuance
- refresh token issuance
- login persistence via refresh token model

### excluded
- logout
- password reset
- email verification
- 2FA
- signup
- UI design
- implementation details

## login_identifier_rules
- type: non-email plain identifier
- charset: lowercase letters, digits, hyphen
- min_length: 4
- max_length: 32
- disallow:
  - spaces
  - uppercase letters
  - underscore
  - special characters except hyphen

## password_rules
- min_length: 8
- max_length: 64
- required:
  - at least 1 letter
  - at least 1 digit
- optional:
  - special characters
- disallow:
  - spaces

## authentication_model
- token_pair_required: true
- access_token: used for protected API access
- refresh_token: used for login persistence and token reissue flow
- token_delivery: response body
- token_type: Bearer
- access_token_expires_in_seconds: 1800
- refresh_token_expires_in_seconds: 604800

## login_endpoint
- method: POST
- path: `/api/v1/auth/login`
- request_body:
  - `loginId: string`
  - `password: string`
- success_status: 200
- success_fields:
  - `accessToken`
  - `refreshToken`
  - `accessTokenExpiresIn`
  - `refreshTokenExpiresIn`
  - `tokenType`

## error_contract
- 400:
  - `INVALID_REQUEST_BODY`
  - `INVALID_LOGIN_ID_FORMAT`
  - `INVALID_PASSWORD_FORMAT`
  - `MISSING_REQUIRED_FIELD`
- 401:
  - `INVALID_CREDENTIALS`
- 423:
  - `ACCOUNT_LOCKED`
- 500:
  - `INTERNAL_SERVER_ERROR`

## lock_policy
- consecutive_failed_attempts_limit: 5
- lock_duration_minutes: 10
- lock_on: fifth failed request
- during_lock:
  - login must fail
  - token must not be issued
- on_success:
  - failed_attempt_count resets

## security_requirements
- input validation required
- credential failure responses must avoid account enumeration
- sensitive internal details must not be exposed
- token issuance, auth failure, and lock transitions must be traceable

## acceptance_checkpoints
- valid login succeeds
- unknown account fails
- wrong password fails
- invalid loginId format fails
- invalid password format fails
- fifth consecutive failure locks account
- locked account remains blocked during lock window
- success resets failure count

## follow_up_documents_needed
- token reissue API spec
- password storage policy
- audit log detail spec
- operator unlock procedure

## sync_note
- This file was generated to sync `.ai/context` with newly added login source documents.
- If `doc/auth/login.md` or `doc/auth/login-api.md` changes, this file must be re-reviewed.
