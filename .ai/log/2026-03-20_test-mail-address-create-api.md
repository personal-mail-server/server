# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 기능 추가

## 기능명
- 테스트용 메일 주소 생성 API 구현

## 관련 키워드
- 테스트용 메일 주소
- 생성 API
- Bearer 인증
- 유니크 검증
- OpenAPI

## 변경 이유
- PR 3 범위에 따라 사용자가 직접 입력하거나 후보 API에서 받은 이메일로 테스트용 메일 주소를 실제 생성할 수 있어야 했기 때문

## 변경 대상
- `code/internal/auth/errors.go`
- `code/internal/testaddress/types.go`
- `code/internal/testaddress/validation.go`
- `code/internal/testaddress/service.go`
- `code/internal/testaddress/service_test.go`
- `code/internal/http/handlers/testaddress_handler.go`
- `code/internal/http/handlers/testaddress_handler_test.go`
- `code/internal/http/router/router.go`
- `code/openapi/openapi.yaml`
- `doc/mail/test-address-create-api.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `.ai/context/project-current-stack.md`
- `.ai/context/project-current-testing.md`
- `.ai/context/mail-test-address-create-api.md`
- `.ai/log/2026-03-20_test-mail-address-create-api.md`

## 변경 내용
- 테스트용 메일 주소 생성 요청 DTO와 이메일 검증 로직을 추가함
- 액세스 토큰 검증 후 사용자 소유로 테스트용 메일 주소를 생성하는 서비스를 추가함
- 생성 핸들러, 라우터, OpenAPI 계약을 추가하고 중복 생성 시 충돌 응답을 반환하도록 반영함
- 현재 구성/테스트 문서와 AI 컨텍스트를 실제 구현 상태에 맞게 갱신함

## 영향 범위
- 테스트용 메일 주소 생성 기능
- OpenAPI 계약
- 테스트용 메일 주소 관련 보호 API 범위

## 비고
- 현재 단계는 생성 API까지 포함하며 목록/상세/수정/삭제는 아직 구현되지 않음
