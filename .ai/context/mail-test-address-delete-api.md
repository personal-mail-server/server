# mail-test-address-delete-api

## purpose
- AI-only structured context for the test mail address delete API document.
- Source of truth remains `doc/api/v1/mails/delete/backend.md`.

## source_documents
- `doc/api/v1/mails/delete/backend.md`
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/06-delete-api.md`

## api_summary
- endpoint: `DELETE /api/v1/mails/{id}`
- auth: bearer access token required
- delete_model: soft delete via `deleted_at`
- success_status: `204`
- non_owner_strategy: return `404`

## sync_note
- If delete API contract or soft-delete policy changes, re-review this file with `doc/api/v1/mails/delete/backend.md`.
