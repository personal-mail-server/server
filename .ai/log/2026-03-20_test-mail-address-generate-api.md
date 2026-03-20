# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 기능 추가

## 기능명
- 테스트용 메일 주소 후보 생성 API 구현

## 관련 키워드
- 테스트용 메일 주소
- 후보 생성 API
- Bearer 인증
- OpenAPI
- mail.local

## 변경 이유
- PR 2 범위에 따라 사용자가 생성 화면에서 유니크 메일 주소 후보를 요청할 수 있는 보호 API가 필요했기 때문

## 변경 대상
- `code/internal/testaddress/service.go`
- `code/internal/testaddress/service_test.go`
- `code/internal/http/handlers/testaddress_handler.go`
- `code/internal/http/handlers/testaddress_handler_test.go`
- `code/internal/http/router/router.go`
- `code/internal/app/server.go`
- `code/openapi/openapi.yaml`
- `doc/mail/test-address-generate-api.md`
- `.ai/context/mail-test-address-generate-api.md`
- `.ai/log/2026-03-20_test-mail-address-generate-api.md`

## 변경 내용
- 액세스 토큰 검증 후 유니크 메일 주소 후보를 반환하는 서비스와 핸들러를 추가함
- 테스트용 메일 주소 후보 생성 엔드포인트를 라우터와 OpenAPI에 반영함
- 현재 구현 기본 도메인을 `mail.local` 로 두고 API 문서에 기록함

## 영향 범위
- 향후 테스트용 메일 주소 생성 화면
- OpenAPI 계약
- 인증된 보호 API 범위

## 비고
- 현재는 후보 생성만 포함하며 실제 주소 저장은 다음 PR 범위임
