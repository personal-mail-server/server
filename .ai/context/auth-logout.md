# auth-logout

## purpose
- AI-only structured context for logout-related source documents.
- Not an authoritative document.
- Source of truth remains `doc/api/v1/auth/logout/design.md`, `doc/api/v1/auth/logout/backend.md`, and `doc/api/v1/auth/logout/frontend.md`.

## source_documents
- `doc/api/v1/auth/logout/design.md`
- `doc/api/v1/auth/logout/backend.md`
- `doc/api/v1/auth/logout/frontend.md`
- `doc/api/v1/auth/login/design.md`
- `doc/api/v1/auth/login/backend.md`
- `doc/api/v1/auth/login/frontend.md`
- `TECH_STACK.md`
- `TRIGGERS.md`

## path_note
- Repository actual paths currently use `doc/`, `code/`, `.ai/` at root.
- Some governing docs describe these conceptually as `project/doc`, `project/code`, `project/.ai`.
- For current local file access, use actual root paths first.

## document_status
- logout feature design exists: `doc/api/v1/auth/logout/design.md`
- logout API spec exists: `doc/api/v1/auth/logout/backend.md`
- logout frontend behavior doc exists: `doc/api/v1/auth/logout/frontend.md`
- current codebase includes logout runtime path and frontend logout entry
- logout-specific API details have been confirmed by the user and reflected in source documents

## current_constraints
- current login docs explicitly exclude logout from their scope
- current code issues JWT access/refresh tokens only
- current code uses session versioning for account-wide invalidation
- frontend currently exposes login and authenticated state transitions on the root page

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
- dedicated protected API validation middleware is still not implemented

## frontend_status
- current root page includes logout entry after successful login
- logout action should live in authenticated screens only
- logout action placement: user menu
- successful logout must clear client auth state and return to the logged-out root-page state at `/`
- auth failure should also clear client auth state and return to the logged-out root-page state at `/`

## verification_checkpoints
- authenticated logout succeeds
- authentication failure returns `401 Unauthorized`
- logout invalidates all sessions for the account
- revoked session cannot access protected APIs
- revoked session cannot reissue tokens

## sync_note
- If logout scope, API contract, session invalidation policy, or frontend behavior changes, this file must be re-reviewed together with the source documents.
