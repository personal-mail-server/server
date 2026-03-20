# mail-test-address-generate-api

## purpose
- AI-only structured context for the test mail address generate API document.
- Source of truth remains `doc/api/v1/mails/generate/backend.md`.

## source_documents
- `doc/api/v1/mails/generate/backend.md`
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/02-generate-api.md`

## api_summary
- endpoint: `POST /api/v1/mails/generate`
- auth: bearer access token required
- request_body: none
- response: `{ "email": string }`
- current_default_domain: `mail.local`

## sync_note
- If the generate API contract or candidate email domain changes, re-review this file with `doc/api/v1/mails/generate/backend.md`.
