# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 기능 추가

## 기능명
- 테스트용 메일 주소 논리 삭제 API 구현

## 관련 키워드
- 테스트용 메일 주소
- 논리 삭제 API
- deleted_at
- Bearer 인증
- 소유권 검증

## 변경 이유
- PR 6 범위에 따라 인증된 사용자가 자신의 테스트용 메일 주소를 물리 삭제가 아닌 논리 삭제로 처리할 수 있어야 했기 때문

## 변경 대상
- `code/internal/testaddress/service.go`
- `code/internal/testaddress/service_test.go`
- `code/internal/http/handlers/testaddress_handler.go`
- `code/internal/http/handlers/testaddress_handler_test.go`
- `code/internal/http/router/router.go`
- `code/openapi/openapi.yaml`
- `doc/mail/test-address-delete-api.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `.ai/context/project-current-stack.md`
- `.ai/context/project-current-testing.md`
- `.ai/context/mail-test-address-delete-api.md`
- `.ai/log/2026-03-20_test-mail-address-delete-api.md`

## 변경 내용
- 테스트용 메일 주소 논리 삭제 서비스와 핸들러를 추가함
- 소유권 확인 후 `deleted_at` 기록으로 삭제하도록 반영함
- 삭제 라우트와 OpenAPI 계약을 추가하고, 현재 구성/테스트 문서와 AI 컨텍스트를 실제 구현 상태에 맞게 갱신함

## 영향 범위
- 테스트용 메일 주소 삭제 기능
- OpenAPI 계약
- 테스트용 메일 주소 보호 API 범위

## 비고
- 현재 단계는 삭제 API까지 포함하며 프론트엔드는 아직 구현되지 않음
