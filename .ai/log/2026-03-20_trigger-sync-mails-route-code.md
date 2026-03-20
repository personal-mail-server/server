# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 동기화

## 기능명
- `동기화` 트리거 실행: mails 사양과 실제 라우트/호출 경로 일치화

## 관련 키워드
- 동기화 트리거
- mails API
- 라우터
- 프론트엔드 호출 경로
- 핸들러 테스트

## 변경 이유
- 사양 문서와 OpenAPI가 `/api/v1/mails`로 변경된 상태에서 실제 코드 경로가 `/api/v1/test-addresses`로 남아 있어 불일치를 해소하기 위해

## 변경 대상
- `code/internal/http/router/router.go`
- `code/frontend/app.js`
- `code/internal/http/handlers/testaddress_handler_test.go`
- `.ai/log/2026-03-20_trigger-sync-mails-route-code.md`

## 변경 내용
- Echo 라우터 그룹 경로를 `/test-addresses`에서 `/mails`로 변경
- 프론트엔드 API 호출/상태 메시지 경로를 `/api/v1/mails...`로 변경
- 핸들러 테스트의 요청 URL 및 `SetPath` 값을 `/api/v1/mails...`로 변경

## 영향 범위
- 테스트용 메일 주소 API 실제 진입 URL
- 프론트엔드 테스트용 메일 주소 관리 동작
- 관련 핸들러 테스트 경로 검증

## 비고
- 검증: `go test ./...` 통과, `go build ./cmd/server` 통과
