# 작업 로그

## 작업 일시
- 2026-03-14

## 작업 유형
- 문서 수정

## 기능명
- `.ai/context` 우선 참조 규칙 강화

## 관련 키워드
- AGENTS
- TRIGGERS
- doc/README
- .ai/context
- 우선 참조
- 교차 검증

## 변경 이유
- 작업 중 `.ai/context`가 갱신만 되고 실제 시작 단계에서 우선 참조되지 않는 문제가 반복되어, 작업 시작 절차에 우선 조회 및 원본 문서 교차 검증 규칙을 명시할 필요가 있었기 때문

## 변경 대상
- `AGENTS.md`
- `TRIGGERS.md`
- `doc/README.md`
- `.ai/context/triggers.md`
- `.ai/context/READEME.md`
- `.ai/log/2026-03-14_context-priority-reference-rule.md`

## 변경 내용
- `AGENTS.md` 문서 기반 작업 원칙에 `.ai/context` 우선 조회와 원본 문서 교차 검증 규칙 추가
- `TRIGGERS.md` 기본 원칙 및 각 트리거 동작에 `.ai/context` 선조회 규칙 추가
- `doc/README.md` 기본 원칙과 문서-코드 관계에 `.ai/context` 우선 참조 및 원본 우선 원칙 추가
- 관련 AI context 문서(`triggers`, `READEME`)에 동일한 우선 참조 절차 반영

## 영향 범위
- 트리거 실행 시작 절차
- 문서 기반 동기화 절차
- AI 문서 탐색 우선순위

## 비고
- `.ai/context`는 우선 조회 대상이지만 최종 기준은 원본 문서라는 원칙을 명시적으로 유지했다.
