# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 문서 사양 변경

## 기능명
- 테스트용 메일 주소 API URL 루트 `test-addresses` -> `mails` 전환

## 관련 키워드
- API URL
- mails
- OpenAPI
- 문서 경로 정리
- AI 컨텍스트 동기화

## 변경 이유
- API 리소스 명을 복수형 `mails`로 통일해 URL 명명 규칙을 일관되게 유지하기 위해

## 변경 대상
- `doc/api/README.md`
- `doc/api/v1/mails/frontend.md`
- `doc/api/v1/mails/generate/backend.md`
- `doc/api/v1/mails/generate/frontend.md`
- `doc/api/v1/mails/create/backend.md`
- `doc/api/v1/mails/create/frontend.md`
- `doc/api/v1/mails/read/backend.md`
- `doc/api/v1/mails/read/frontend.md`
- `doc/api/v1/mails/update/backend.md`
- `doc/api/v1/mails/update/frontend.md`
- `doc/api/v1/mails/delete/backend.md`
- `doc/api/v1/mails/delete/frontend.md`
- `doc/project/current-stack.md`
- `code/openapi/openapi.yaml`
- `.ai/context/mail-test-address-generate-api.md`
- `.ai/context/mail-test-address-create-api.md`
- `.ai/context/mail-test-address-read-api.md`
- `.ai/context/mail-test-address-update-api.md`
- `.ai/context/mail-test-address-delete-api.md`
- `.ai/context/mail-test-address-frontend.md`
- `.ai/context/mail-test-address-frontend-update-delete.md`
- `.ai/log/2026-03-20_mails-url-spec-rename.md`

## 변경 내용
- 문서 폴더를 `doc/api/v1/test-addresses/`에서 `doc/api/v1/mails/`로 이동함
- 문서 본문의 엔드포인트 경로를 `/api/v1/test-addresses...`에서 `/api/v1/mails...`로 일괄 변경함
- 문서 간 참조 경로를 `doc/api/v1/mails/...` 기준으로 동기화함
- OpenAPI 계약 경로를 `/api/v1/mails`, `/api/v1/mails/{id}`, `/api/v1/mails/generate`로 갱신함
- `.ai/context`의 원본 문서 참조 및 엔드포인트 요약 경로를 함께 동기화함

## 영향 범위
- API 문서 탐색 경로
- OpenAPI 계약 경로 정의
- AI 컨텍스트 문서의 참조 경로 및 API 요약

## 비고
- 이번 변경은 사양/문서 기준 정리이며, 백엔드 라우터 및 프론트엔드 실제 호출 경로는 아직 `test-addresses`를 사용한다.
