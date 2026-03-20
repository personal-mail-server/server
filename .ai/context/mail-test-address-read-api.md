# mail-test-address-read-api

## purpose
- AI-only structured context for the test mail address read API document.
- Source of truth remains `doc/mail/test-address-read-api.md`.

## source_documents
- `doc/mail/test-address-read-api.md`
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/04-read-api.md`

## api_summary
- list_endpoint: `GET /api/v1/test-addresses`
- detail_endpoint: `GET /api/v1/test-addresses/{id}`
- auth: bearer access token required
- list_scope: current owner only, soft-deleted excluded
- detail_scope: current owner only, non-owner hidden as 404

## sync_note
- If read API contract, ownership rule, or list/detail response shape changes, re-review this file with `doc/mail/test-address-read-api.md`.
