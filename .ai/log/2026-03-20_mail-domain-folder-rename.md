# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 문서 구조 변경

## 기능명
- 메일 도메인 문서 폴더를 `doc/mail`에서 `doc/domain/mail`로 이관

## 관련 키워드
- 문서 구조
- 도메인 문서
- URL 루트 API 문서
- 참조 경로 동기화

## 변경 이유
- API 문서(`doc/api/v1/...`)와 도메인 문서의 역할을 분리해 사람이 문서를 볼 때 목적별 탐색이 더 명확하도록 정리하기 위해

## 변경 대상
- `doc/domain/mail/server-requirements.md`
- `doc/domain/mail/test-address-model.md`
- `doc/domain/mail/pr-split/README.md`
- `doc/domain/mail/pr-split/01-storage-foundation.md`
- `doc/domain/mail/pr-split/02-generate-api.md`
- `doc/domain/mail/pr-split/03-create-api.md`
- `doc/domain/mail/pr-split/04-read-api.md`
- `doc/domain/mail/pr-split/05-update-api.md`
- `doc/domain/mail/pr-split/06-delete-api.md`
- `doc/domain/mail/pr-split/07-frontend-read-create.md`
- `doc/domain/mail/pr-split/08-frontend-update-delete.md`
- `doc/api/v1/test-addresses/create/backend.md`
- `doc/api/v1/test-addresses/generate/backend.md`
- `doc/api/v1/test-addresses/read/backend.md`
- `doc/api/v1/test-addresses/update/backend.md`
- `doc/api/v1/test-addresses/delete/backend.md`
- `doc/api/v1/test-addresses/frontend.md`
- `doc/project/current-stack.md`
- `.ai/log/2026-03-20_mail-domain-folder-rename.md`

## 변경 내용
- 기존 `doc/mail` 하위 도메인 문서를 동일한 상대 구조로 `doc/domain/mail`로 이동함
- 문서 내부 및 API 문서에서 참조하던 `doc/mail/...` 경로를 `doc/domain/mail/...`로 갱신함
- 현재 구성 문서에 도메인 문서 루트(`doc/domain/mail/`)를 반영함

## 영향 범위
- 메일 도메인 요구사항/모델/계획 문서 접근 경로
- API 문서의 기준 문서 참조 경로

## 비고
- 코드 동작 변경 없음 (문서 경로/참조만 변경)
