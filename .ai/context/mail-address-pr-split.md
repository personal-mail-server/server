# mail-address-pr-split

## purpose
- AI-only structured context for the test mail address PR split documents.
- Not an authoritative document.
- Source of truth remains `doc/mail/pr-split/*.md`.

## source_documents
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/README.md`
- `doc/mail/pr-split/01-storage-foundation.md`
- `doc/mail/pr-split/02-generate-api.md`
- `doc/mail/pr-split/03-create-api.md`
- `doc/mail/pr-split/04-read-api.md`
- `doc/mail/pr-split/05-update-api.md`
- `doc/mail/pr-split/06-delete-api.md`
- `doc/mail/pr-split/07-frontend-read-create.md`
- `doc/mail/pr-split/08-frontend-update-delete.md`

## split_summary
- PR 1: storage foundation and logical-delete-ready schema
- PR 2: unique email candidate generation API
- PR 3: create API
- PR 4: list/detail read APIs
- PR 5: update API
- PR 6: logical delete API
- PR 7: frontend list/detail/create
- PR 8: frontend update/delete

## sequencing
- recommended order is storage -> generate -> create -> read -> update -> delete -> frontend read/create -> frontend update/delete

## guardrails
- do not mix storage, full CRUD APIs, and frontend in one PR
- keep logical delete policy isolated to delete-focused work
- each PR should update its own docs, tests, and OpenAPI if API surface changes

## sync_note
- If PR boundaries or implementation order change, this file must be re-reviewed together with `doc/mail/pr-split/*.md`.
