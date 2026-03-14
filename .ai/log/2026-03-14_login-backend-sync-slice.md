# 작업 로그

## 작업 일시
- 2026-03-14

## 작업 유형
- 기능 추가

## 기능명
- 로그인 백엔드 동기화 1차 슬라이스

## 관련 키워드
- 로그인
- Echo
- Go
- PostgreSQL
- 잠금 정책
- OpenAPI
- Docker

## 변경 이유
- `doc/auth/login.md`, `doc/auth/login-api.md`에 정의된 로그인 기능/계약을 `code/` 기준으로 최초 실행 가능한 백엔드로 동기화하기 위해

## 변경 대상
- `project/code/cmd/server/main.go`
- `project/code/internal/app/server.go`
- `project/code/internal/config/config.go`
- `project/code/internal/db/postgres.go`
- `project/code/internal/db/migrate.go`
- `project/code/internal/auth/*.go`
- `project/code/internal/http/handlers/*.go`
- `project/code/internal/http/router/router.go`
- `project/code/migrations/001_create_users.sql`
- `project/code/openapi/openapi.yaml`
- `project/code/Dockerfile`
- `project/code/docker-compose.yml`
- `project/code/.dockerignore`
- `project/code/go.mod`
- `project/code/go.sum`
- `project/code/README.md`

## 변경 내용
- `POST /api/v1/auth/login` 구현 및 문서 기준 상태 코드/에러 코드 매핑 적용
- 로그인 입력 검증 규칙 구현(`loginId`, `password` 규칙)
- 연속 실패 카운트, 5회째 10분 잠금, 잠금 중 차단, 성공 시 실패 횟수 초기화 구현
- 액세스 토큰 30분, 리프레시 토큰 7일 토큰 발급 구현
- PostgreSQL 테이블/마이그레이션 추가 및 서버 시작 시 마이그레이션 실행
- OpenAPI 계약 파일 및 Swagger UI 접근 경로(`/docs`, `/docs/openapi.yaml`) 추가
- 백엔드+DB 실행용 Dockerfile/Compose 구성 추가
- 검증/인증/잠금 동작 테스트 추가

## 영향 범위
- 인증(로그인) API
- 인증 상태 저장 스키마
- 백엔드 실행/테스트 환경

## 비고
- LSP 진단 도구(`gopls`)가 실행 PATH에 없어 `lsp_diagnostics` 확인은 불가했으며, 대신 `go test ./...`, `go build ./cmd/server`, `go vet ./...`로 검증 완료
