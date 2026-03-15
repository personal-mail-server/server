# auth-token-reissue

## purpose
- AI-only structured context for token reissue-related source documents.
- Not an authoritative document.
- Source of truth remains `doc/auth/token-reissue.md`, `doc/auth/token-reissue-api.md`, and `doc/auth/token-reissue-frontend.md`.

## source_documents
- `doc/auth/token-reissue.md`
- `doc/auth/token-reissue-api.md`
- `doc/auth/token-reissue-frontend.md`
- `doc/auth/login.md`
- `doc/auth/login-api.md`
- `doc/auth/login-frontend.md`
- `doc/auth/logout.md`
- `doc/auth/logout-api.md`
- `doc/auth/logout-frontend.md`
- `TECH_STACK.md`
- `TRIGGERS.md`

## path_note
- Repository actual paths currently use `doc/`, `code/`, `.ai/` at root.
- Some governing docs describe these conceptually as `project/doc`, `project/code`, `project/.ai`.
- For current local file access, use actual root paths first.

## document_status
- token reissue feature design exists: `doc/auth/token-reissue.md`
- token reissue API spec exists: `doc/auth/token-reissue-api.md`
- token reissue frontend behavior doc exists: `doc/auth/token-reissue-frontend.md`
- current codebase now exposes a documented token reissue endpoint and persists refresh-token state for rotation

## current_constraints
- login currently issues access and refresh tokens with 30-minute and 7-day expiry windows
- JWT claims currently include token use and session version information
- logout currently invalidates account sessions via session version increment
- refresh token state is persisted in a dedicated table keyed by refresh-token `jti`
- frontend currently keeps auth state in memory only
- protected-screen and auto-reissue flows are not yet implemented
- current documented draft now requires one-time use invalidation for a successfully used refresh token

## token_reissue_scope

### included
- refresh-token based token reissue
- refresh-token validation and token-use enforcement
- session-version verification against current account state
- issuance of a new access token and a new refresh token
- immediate invalidation of the previously used refresh token after successful reissue
- frontend state replacement rules after successful reissue

### excluded
- login
- logout
- signup
- password reset
- 2FA
- persistent browser storage policy
- detailed protected-screen UI

## token_reissue_model
- endpoint:
  - method: POST
  - path: `/api/v1/auth/token/reissue`
- request_body:
  - `refreshToken: string`
- success_status: `200 OK`
- success_fields:
  - `accessToken`
  - `refreshToken`
  - `accessTokenExpiresIn`
  - `refreshTokenExpiresIn`
  - `tokenType`
- auth_failure_status: `401 Unauthorized`
- validation_failure_status: `400 Bad Request`
- internal_error_status: `500 Internal Server Error`
- implementation_status:
  - backend endpoint implemented
  - openapi contract implemented
  - service-level refresh rotation implemented
  - frontend automatic reissue flow still not implemented

## verification_checkpoints
- valid refresh token reissue succeeds
- missing request body field fails with validation error
- access token misuse fails
- expired or forged refresh token fails
- logged-out session cannot reissue tokens
- already-used refresh token cannot reissue again
- successful reissue replaces client auth state with new token pair

## sync_note
- If token reissue scope, API contract, session-version validation rules, or frontend behavior changes, this file must be re-reviewed together with the source documents.
