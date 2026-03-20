# mail-test-address-model

## purpose
- AI-only structured context for the test mail address data model document.
- Not an authoritative document.
- Source of truth remains `doc/mail/test-address-model.md`.

## source_documents
- `doc/mail/test-address-model.md`
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/01-storage-foundation.md`
- `doc/project/current-stack.md`
- `TECH_STACK.md`

## model_summary
- table purpose: store user-owned test mail addresses
- recommended table name: `test_mail_addresses`
- required fields: `id`, `owner_user_id`, `email`, `created_at`, `updated_at`, `deleted_at`
- deletion model: soft delete via nullable `deleted_at`
- uniqueness model: system-wide unique `email`

## repository_expectation
- repository should support create/get-by-id/get-by-email/list-by-owner/update/soft-delete
- default reads should prefer active rows where `deleted_at IS NULL`

## open_questions
- `id` type
- reuse policy for soft-deleted email values
- deleted-detail visibility
- email normalization strategy

## sync_note
- If the test mail address schema, uniqueness rule, or soft-delete rule changes, this file must be re-reviewed together with `doc/mail/test-address-model.md`.
