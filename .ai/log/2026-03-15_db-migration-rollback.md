# 작업 로그

## 작업 일시
- 2026-03-15

## 작업 유형
- 기능 추가

## 기능명
- DB 마이그레이션 롤백 대응 및 동기화 규칙 문서화

## 관련 키워드
- 데이터베이스
- 마이그레이션
- 롤백
- PostgreSQL
- Makefile
- 문서 동기화

## 변경 이유
- 현재 저장소는 업 마이그레이션만 지원하고 롤백 대응과 DB 변경 작업 규칙 문서가 없어서, 수동 롤백 경로와 기준 문서를 함께 마련해야 했기 때문

## 변경 대상
- `code/internal/config/config.go`
- `code/internal/db/migrate.go`
- `code/internal/db/migrate_test.go`
- `code/cmd/migrate/main.go`
- `code/migrations/001_create_users.down.sql`
- `code/migrations/002_seed_login_user.down.sql`
- `code/migrations/003_update_seed_login_user_password.down.sql`
- `code/migrations/004_add_users_session_version.down.sql`
- `code/migrations/005_create_refresh_tokens.down.sql`
- `code/Dockerfile`
- `Makefile`
- `README.md`
- `code/README.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `doc/project/database-migration.md`
- `.ai/context/project-current-stack.md`
- `.ai/context/project-current-testing.md`
- `.ai/context/project-database-migration.md`
- `.ai/log/2026-03-15_db-migration-rollback.md`

## 변경 내용
- 커스텀 SQL 마이그레이션 러너에 업 적용 함수와 단계별 롤백 함수 추가
- 전용 마이그레이션 CLI 엔트리포인트 추가
- 기존 업 마이그레이션에 대응하는 다운 마이그레이션 파일 추가
- 백엔드 Docker 이미지에 마이그레이션 바이너리를 포함하고 Make 명령이 해당 이미지를 사용하도록 조정
- 마이그레이션 파일 탐색 및 다운 경로 계산 로직에 대한 단위 테스트 추가
- 현재 구성 문서와 현재 테스트 문서에 롤백/검증 흐름 반영
- DB 마이그레이션 및 동기화 시 코드 작성 규칙을 별도 기준 문서로 추가
- 대응하는 AI 컨텍스트 문서를 원본 문서 기준으로 동기화

## 영향 범위
- 데이터베이스 스키마 변경 절차
- 로컬 수동 마이그레이션/롤백 실행 방식
- 백엔드 Docker 빌드 산출물
- 프로젝트 기준 문서와 AI 검색용 컨텍스트

## 비고
- 검증은 `go test ./...`, `go build ./...`, `go vet ./...`, `make migrate-up`, `make migrate-down STEPS=1`, 재적용 복구까지 수행함
