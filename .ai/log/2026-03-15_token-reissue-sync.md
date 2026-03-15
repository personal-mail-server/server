# 작업 로그

## 작업 일시
- 2026-03-15

## 작업 유형
- 동기화

## 기능명
- 토큰 재발급 문서-코드 동기화

## 관련 키워드
- 인증
- 토큰 재발급
- 리프레시 토큰
- 세션 버전
- OpenAPI
- 테스트

## 변경 이유
- `doc/auth/token-reissue*.md` 기준 문서와 현재 코드 상태 사이에 구현 공백이 있어 동기화가 필요했기 때문

## 변경 대상
- `code/migrations/005_create_refresh_tokens.sql`
- `code/internal/auth/types.go`
- `code/internal/auth/errors.go`
- `code/internal/auth/validation.go`
- `code/internal/auth/token.go`
- `code/internal/auth/service.go`
- `code/internal/auth/postgres_repository.go`
- `code/internal/http/handlers/auth_handler.go`
- `code/internal/http/router/router.go`
- `code/internal/auth/service_test.go`
- `code/internal/auth/validation_test.go`
- `code/internal/http/handlers/auth_handler_test.go`
- `code/openapi/openapi.yaml`
- `code/README.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `.ai/context/auth-token-reissue.md`
- `.ai/context/project-current-stack.md`
- `.ai/context/project-current-testing.md`
- `.ai/log/2026-03-15_token-reissue-sync.md`

## 변경 내용
- 리프레시 토큰 `jti` 상태 저장용 테이블을 추가하고 1회성 consume-and-replace 방식의 토큰 재발급 구조를 구현
- 로그인 시 리프레시 토큰 상태 저장, 재발급 시 기존 토큰 즉시 무효화와 새 토큰 등록을 같은 흐름으로 반영
- 토큰 재발급 핸들러, 라우터, OpenAPI 계약, 단위 테스트를 추가
- 비밀번호 검증을 문서/프론트 규칙과 맞도록 ASCII 영문자 기준으로 정리
- 현재 구성/테스트 문서와 AI 컨텍스트를 실제 구현 상태에 맞게 갱신

## 영향 범위
- 인증 API
- JWT 발급 및 검증
- PostgreSQL 인증 상태 저장
- OpenAPI 문서
- 인증 테스트 체계

## 비고
- 로그아웃의 계정 전체 무효화는 기존 `session_version` 기준을 유지하고, 토큰 재발급만 개별 리프레시 토큰 상태 저장으로 보강함
