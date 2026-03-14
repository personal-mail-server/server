# 작업 로그

## 작업 일시
- 2026-03-14

## 작업 유형
- 기능 추가

## 기능명
- 로그아웃 동기화 구현

## 관련 키워드
- 로그아웃
- 세션 버전 관리
- 액세스 토큰
- 리프레시 토큰
- OpenAPI
- 프론트엔드
- 동기화

## 변경 이유
- 예약 트리거 `동기화`에 따라 기준 문서에 정의된 로그아웃 API, 세션 무효화, 프론트엔드 로그아웃 흐름을 현재 코드와 설정에 실제로 반영할 필요가 있었기 때문

## 변경 대상
- `code/internal/auth/types.go`
- `code/internal/auth/errors.go`
- `code/internal/auth/token.go`
- `code/internal/auth/service.go`
- `code/internal/auth/postgres_repository.go`
- `code/internal/auth/service_test.go`
- `code/internal/http/handlers/auth_handler.go`
- `code/internal/http/handlers/auth_handler_test.go`
- `code/internal/http/router/router.go`
- `code/migrations/004_add_users_session_version.sql`
- `code/openapi/openapi.yaml`
- `code/frontend/index.html`
- `code/frontend/styles.css`
- `code/frontend/app.js`
- `.github/workflows/ci.yml`
- `doc/auth/login-frontend.md`
- `doc/auth/logout.md`
- `doc/auth/logout-api.md`
- `doc/auth/logout-frontend.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `.ai/context/auth-login.md`
- `.ai/context/auth-logout.md`
- `.ai/context/project-current-stack.md`
- `.ai/context/project-current-testing.md`
- `.ai/log/2026-03-14_logout-sync-implementation.md`

## 변경 내용
- 로그아웃 API와 세션 버전 기반 전체 세션 무효화 로직 추가
- JWT 클레임에 세션 버전 정보를 포함하고 로그아웃 시 구세션 토큰 재사용 차단
- 로그아웃 핸들러, 라우터, 단위 테스트 추가
- OpenAPI와 CI 검증 흐름에 로그아웃 시나리오 반영
- 프론트엔드에 인증 상태 요약과 사용자 메뉴 기반 로그아웃 흐름 추가
- 현재구성, 현재테스트, 인증 관련 문서와 AI context 동기화

## 영향 범위
- 인증 백엔드 구현
- 인증 프론트엔드 동작
- OpenAPI 계약
- CI Docker smoke 검증
- 현재 구성 및 테스트 문서

## 비고
- 현재 구현은 계정 단위 전체 세션 무효화를 세션 버전 관리 방식으로 처리한다.
- 토큰 재발급과 별도 보호 API 미들웨어는 아직 후속 구현 대상이다.
