# mail-test-address-create-api

## purpose
- AI-only structured context for the test mail address create API document.
- Source of truth remains `doc/mail/test-address-create-api.md`.

## source_documents
- `doc/mail/test-address-create-api.md`
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/03-create-api.md`
- `doc/mail/test-address-model.md`

## api_summary
- endpoint: `POST /api/v1/test-addresses`
- auth: bearer access token required
- request: `{ "email": string }`
- success_status: `201`
- response_shape: created test mail address resource
- duplicate_error_code: `DUPLICATE_TEST_ADDRESS_EMAIL`

## sync_note
- If the create API contract, validation rule, or response shape changes, re-review this file with `doc/mail/test-address-create-api.md`.
