# auth-logout

## purpose
- AI-only structured context for logout-related source documents.
- Not an authoritative document.
- Source of truth remains `doc/auth/logout.md`, `doc/auth/logout-api.md`, and `doc/auth/logout-frontend.md`.

## source_documents
- `doc/auth/logout.md`
- `doc/auth/logout-api.md`
- `doc/auth/logout-frontend.md`
- `doc/auth/login.md`
- `doc/auth/login-api.md`
- `doc/auth/login-frontend.md`
- `TECH_STACK.md`
- `TRIGGERS.md`

## path_note
- Repository actual paths currently use `doc/`, `code/`, `.ai/` at root.
- Some governing docs describe these conceptually as `project/doc`, `project/code`, `project/.ai`.
- For current local file access, use actual root paths first.

## document_status
- logout feature design exists: `doc/auth/logout.md`
- logout API spec exists: `doc/auth/logout-api.md`
- logout frontend behavior doc exists: `doc/auth/logout-frontend.md`
- current codebase still has login-only implementation and no logout runtime path yet
- logout-specific API details have been confirmed by the user and reflected in source documents

## current_constraints
- current login docs explicitly exclude logout from their scope
- current code issues JWT access/refresh tokens only
- current code has no token persistence, revocation store, or logout endpoint
- frontend currently exposes login screen only

## logout_scope

### included
- authenticated user logout
- server-side account-wide session invalidation
- access-token and refresh-token reuse blocking after logout
- logout API contract
- frontend logout action behavior

### excluded
- admin-forced logout
- signup
- password reset
- token reissue API details
- detailed protected-screen UI design

## logout_model
- server_side_required: true
- revoked_scope: all_sessions
- endpoint:
  - method: POST
  - path: `/api/v1/auth/logout`
- auth_header: `Authorization: Bearer <accessToken>`
- success_response: `204 No Content`
- auth_failure_status: `401 Unauthorized`
- internal_error_status: `500 Internal Server Error`
- access_token_invalidation: session versioning

## implementation_gap
- current auth service only supports login issuance
- current repository interface has no session invalidation methods
- actual implementation will require server-side state capable of account-wide session invalidation

## frontend_status
- current login page has no logout UI
- logout action should live in authenticated screens only
- logout action placement: user menu
- successful logout must clear client auth state and return to `/`
- auth failure should also clear client auth state and return to `/`

## verification_checkpoints
- authenticated logout succeeds
- authentication failure returns `401 Unauthorized`
- logout invalidates all sessions for the account
- revoked session cannot access protected APIs
- revoked session cannot reissue tokens

## sync_note
- If logout scope, API contract, session invalidation policy, or frontend behavior changes, this file must be re-reviewed together with the source documents.
