# mail-test-address-frontend-update-delete

## purpose
- AI-only structured context for PR 8 frontend update/delete additions.
- Source of truth remains `doc/api/v1/test-addresses/frontend.md`.

## source_documents
- `doc/api/v1/test-addresses/frontend.md`
- `doc/mail/pr-split/08-frontend-update-delete.md`
- `doc/api/v1/test-addresses/update/backend.md`
- `doc/api/v1/test-addresses/delete/backend.md`

## ui_summary
- detail panel now includes inline update form
- detail panel now includes delete action
- duplicate update errors surface in status panel
- delete success removes item from list and reselects another item when available

## sync_note
- If update/delete UI behavior changes, re-review this file with `doc/api/v1/test-addresses/frontend.md`.
