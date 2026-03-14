# triggers

## purpose
- AI-only structured context for reserved trigger behavior in this repository.
- Not an authoritative document.
- Source of truth remains `TRIGGERS.md`.

## source_documents
- `TRIGGERS.md`

## reserved_triggers
- `동기화`
- `도큐멘트`
- `현재구성`
- `현재테스트`
- `PUSH`

## push_trigger_summary
- purpose: push current feature branch, create or reuse PR to `main`, and enable auto-merge
- local_entrypoint: `make push-trigger`
- implementation_path: `code/cmd/push-trigger/main.go`
- guardrails:
  - refuse on `main` or `master`
  - refuse on dirty worktree
  - require `origin`
  - require `gh` CLI and authenticated session
  - no force-push
  - no admin bypass merge
- out_of_scope:
  - Copilot review automation
  - automatic review comment fixing
  - automatic reply/resolve of review comments
- pr_template:
  - language: Korean
  - sections:
    - `## 요약`
    - `## 기준 정보`

## sync_note
- If trigger semantics change, this file must be re-reviewed together with `TRIGGERS.md`.
