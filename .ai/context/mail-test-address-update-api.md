# mail-test-address-update-api

## purpose
- AI-only structured context for the test mail address update API document.
- Source of truth remains `doc/api/v1/test-addresses/update/backend.md`.

## source_documents
- `doc/api/v1/test-addresses/update/backend.md`
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/05-update-api.md`

## api_summary
- endpoint: `PUT /api/v1/test-addresses/{id}`
- auth: bearer access token required
- request: `{ "email": string }`
- self_value_retention: allowed
- duplicate_error_code: `DUPLICATE_TEST_ADDRESS_EMAIL`
- non_owner_strategy: return `404`

## sync_note
- If update API contract, duplicate policy, or non-owner visibility rule changes, re-review this file with `doc/api/v1/test-addresses/update/backend.md`.
