# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 기능 추가

## 기능명
- 테스트용 메일 주소 목록/상세 조회 API 구현

## 관련 키워드
- 테스트용 메일 주소
- 목록 조회 API
- 상세 조회 API
- Bearer 인증
- 소유권 검증

## 변경 이유
- PR 4 범위에 따라 인증 사용자가 자신의 테스트용 메일 주소 목록과 상세를 조회할 수 있어야 했기 때문

## 변경 대상
- `code/internal/auth/errors.go`
- `code/internal/testaddress/types.go`
- `code/internal/testaddress/service.go`
- `code/internal/testaddress/service_test.go`
- `code/internal/http/handlers/testaddress_handler.go`
- `code/internal/http/handlers/testaddress_handler_test.go`
- `code/internal/http/router/router.go`
- `code/openapi/openapi.yaml`
- `doc/mail/test-address-read-api.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `.ai/context/project-current-stack.md`
- `.ai/context/project-current-testing.md`
- `.ai/context/mail-test-address-read-api.md`
- `.ai/log/2026-03-20_test-mail-address-read-api.md`

## 변경 내용
- 테스트용 메일 주소 목록 조회와 상세 조회 서비스 로직을 추가함
- 상세 조회에서 타 사용자 리소스를 `404` 로 숨기도록 소유권 검증을 반영함
- 조회 핸들러, 라우터, OpenAPI 계약을 추가함
- 현재 구성/테스트 문서와 AI 컨텍스트를 실제 구현 상태에 맞게 갱신함

## 영향 범위
- 테스트용 메일 주소 읽기 API
- OpenAPI 계약
- 테스트용 메일 주소 관리 보호 API 범위

## 비고
- 현재 단계는 목록/상세 조회까지 포함하며 수정/삭제와 프론트엔드는 아직 구현되지 않음
