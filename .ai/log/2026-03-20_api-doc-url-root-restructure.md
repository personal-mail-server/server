# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 문서 구조 개편

## 기능명
- API 문서 URL 루트 기준 재구성 및 프론트/백엔드 문서 동반 배치

## 관련 키워드
- API 문서 구조
- URL 루트
- 가독성
- 프론트엔드 사양
- 백엔드 사양

## 변경 이유
- 사람이 직접 보는 문서 영역에서 API별 탐색성을 높이고, 같은 API의 프론트/백엔드 사양을 한 위치에서 확인할 수 있도록 정리하기 위해

## 변경 대상
- `doc/api/README.md`
- `doc/api/v1/auth/login/design.md`
- `doc/api/v1/auth/login/backend.md`
- `doc/api/v1/auth/login/frontend.md`
- `doc/api/v1/auth/logout/design.md`
- `doc/api/v1/auth/logout/backend.md`
- `doc/api/v1/auth/logout/frontend.md`
- `doc/api/v1/auth/token/reissue/design.md`
- `doc/api/v1/auth/token/reissue/backend.md`
- `doc/api/v1/auth/token/reissue/frontend.md`
- `doc/api/v1/test-addresses/frontend.md`
- `doc/api/v1/test-addresses/generate/backend.md`
- `doc/api/v1/test-addresses/generate/frontend.md`
- `doc/api/v1/test-addresses/create/backend.md`
- `doc/api/v1/test-addresses/create/frontend.md`
- `doc/api/v1/test-addresses/read/backend.md`
- `doc/api/v1/test-addresses/read/frontend.md`
- `doc/api/v1/test-addresses/update/backend.md`
- `doc/api/v1/test-addresses/update/frontend.md`
- `doc/api/v1/test-addresses/delete/backend.md`
- `doc/api/v1/test-addresses/delete/frontend.md`
- `doc/mail/server-requirements.md`
- `doc/project/current-stack.md`
- `.ai/context/auth-login.md`
- `.ai/context/auth-logout.md`
- `.ai/context/auth-token-reissue.md`
- `.ai/context/mail-server-requirements.md`
- `.ai/context/mail-test-address-generate-api.md`
- `.ai/context/mail-test-address-create-api.md`
- `.ai/context/mail-test-address-read-api.md`
- `.ai/context/mail-test-address-update-api.md`
- `.ai/context/mail-test-address-delete-api.md`
- `.ai/context/mail-test-address-frontend.md`
- `.ai/context/mail-test-address-frontend-update-delete.md`
- `.ai/log/2026-03-20_api-doc-url-root-restructure.md`

## 변경 내용
- 기존 `doc/auth/*`, `doc/mail/test-address-*-api.md`, `doc/mail/test-address-frontend.md` 문서를 `doc/api/v1/...` URL 루트 구조로 이동함
- auth 문서는 엔드포인트 폴더마다 `design.md`, `backend.md`, `frontend.md`를 함께 배치함
- test-addresses 문서는 엔드포인트 액션 폴더(`generate/create/read/update/delete`)로 분리하고 각 폴더에 `backend.md`, `frontend.md`를 배치함
- 문서 내부 참조 경로와 `.ai/context` 참조 경로를 새 구조로 동기화함

## 영향 범위
- API 문서 탐색 경로
- API 관련 문서 간 링크
- AI 컨텍스트 문서의 원본 문서 참조 경로

## 비고
- 코드/테스트 로직은 변경하지 않았고 문서 구조와 참조 경로만 개편함
