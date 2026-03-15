# 작업 로그

## 작업 일시
- 2026-03-15

## 작업 유형
- 문서 동기화

## 기능명
- 전수 조사 기반 사양 문서 정합성 정리

## 관련 키워드
- 사양
- 인증
- 토큰 재발급
- 로그아웃
- 현재 테스트
- OpenAPI

## 변경 이유
- 저장소 전체 전수 조사 결과, 현재 구현 및 CI 흐름과 일부 사양/현재 상태 문서 표현이 어긋난 부분이 확인되었기 때문

## 변경 대상
- `doc/project/current-testing.md`
- `doc/auth/token-reissue.md`
- `doc/auth/token-reissue-api.md`
- `doc/auth/logout.md`
- `doc/auth/logout-frontend.md`
- `code/openapi/openapi.yaml`
- `.ai/context/project-current-testing.md`
- `.ai/context/auth-logout.md`
- `.ai/log/2026-03-15_full-spec-audit-doc-sync.md`

## 변경 내용
- 현재 CI 검증 범위를 실제 GitHub Actions 워크플로와 일치하도록 정리
- 토큰 재발급 문서의 stale 현재 상태 문구를 실제 구현 반영 상태로 수정
- 토큰 재발급 API 명세의 `MISSING_REQUIRED_FIELD` 예시 메시지를 실제 응답 계약과 일치하도록 조정
- 로그아웃 프론트엔드 문서의 화면 전이 표현을 실제 구현인 same-page logged-out state 기준으로 정리
- OpenAPI 상단 설명을 현재 인증 슬라이스 범위에 맞게 갱신
- 대응하는 AI 컨텍스트 문서를 원본 문서 기준으로 동기화
- 토큰 재발급 400 검증 계약에서 존재하지 않는 에러 코드 예시를 제거하고 공백 전용 입력 거부 규칙을 명시
- 현재 테스트 문서에서 토큰 재발급 검증을 자동 실행 검증이 아닌 수동 검증 기준으로 구분해 정리

## 영향 범위
- 인증 관련 기준 문서
- 현재 테스트 상태 문서
- OpenAPI 계약 문서 표현
- AI 정리 문서 검색 정합성

## 비고
- 이번 작업은 문서 및 문서성 계약 파일 정합성 갱신에 한정하며, 별도 코드 동작 변경은 포함하지 않음
