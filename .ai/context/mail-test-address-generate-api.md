# mail-test-address-generate-api

## purpose
- AI-only structured context for the test mail address generate API document.
- Source of truth remains `doc/mail/test-address-generate-api.md`.

## source_documents
- `doc/mail/test-address-generate-api.md`
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/02-generate-api.md`

## api_summary
- endpoint: `POST /api/v1/test-addresses/generate`
- auth: bearer access token required
- request_body: none
- response: `{ "email": string }`
- current_default_domain: `mail.local`

## sync_note
- If the generate API contract or candidate email domain changes, re-review this file with `doc/mail/test-address-generate-api.md`.
