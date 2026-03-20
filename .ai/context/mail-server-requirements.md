# mail-server-requirements

## purpose
- AI-only structured context for the mail server MVP requirements document.
- Not an authoritative document.
- Source of truth remains `doc/mail/server-requirements.md`.

## source_documents
- `doc/mail/server-requirements.md`
- `README.md`
- `TECH_STACK.md`
- `doc/README.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `doc/auth/login.md`
- `doc/auth/logout.md`
- `doc/auth/token-reissue.md`
- `TRIGGERS.md`

## path_note
- Repository actual paths currently use `doc/`, `code/`, `.ai/` at root.
- Some governing docs describe these conceptually as `project/doc`, `project/code`, `project/.ai`.
- For current local file access, use actual root paths first.

## document_status
- narrowed mail-address management requirements doc exists: `doc/mail/server-requirements.md`
- current repository already has auth and Docker runtime foundations but does not yet have mail-domain CRUD docs

## product_direction
- target: test mail address management as the first mail-domain slice
- first_scope: CRUD for test mail addresses
- explicit_non_goal: mail receiving, sending, and policy simulation in this step
- primary_value: let users manage unique test mail addresses without manual uniqueness burden

## mvp_scope

### included
- authenticated access to mail-server features
- test mail address list/detail/create/update/delete
- logical deletion for test mail addresses
- unique-mail validation on create and update
- unique-mail candidate generation API
- web UI flow for list/detail/create/update/delete

### excluded
- mail receiving
- mail storage and message viewing
- policy assignment and scenario simulation
- outbound sending feature
- SMTP infrastructure
- admin console
- mobile app

## core_requirements
- one user can manage multiple test mail addresses
- list and detail views must exist for owned addresses
- create and update must enforce unique email values
- delete must be logical, not physical
- server must expose an API that returns a unique candidate email
- users must still be able to type their own email value directly
- users must only access their own mail-address resources
- Docker-based web verification flow must exist

## open_questions
- domain strategy for generated test addresses
- whether soft-deleted addresses can be reused
- exact editable fields beyond email value
- candidate email generation format and constraints
- whether deleted addresses need filtered list access later

## sync_note
- If mail-address CRUD scope, logical-delete policy, unique-email rules, or later mail-domain design docs change, this file must be re-reviewed together with `doc/mail/server-requirements.md`.
