# 작업 로그

## 작업 일시
- 2026-03-14

## 작업 유형
- 문서 수정

## 기능명
- `.ai/context` 우선 참조 규칙 문구 조정

## 관련 키워드
- AGENTS
- TRIGGERS
- doc/README
- .ai/context
- 우선 참조

## 변경 이유
- 사용자의 의도는 교차 검증 의무 강화가 아니라 `.ai/context`를 작업 시작 기준으로 우선 참조하는 것이므로, 관련 규칙 문구를 그 의도에 맞게 조정할 필요가 있었기 때문

## 변경 대상
- `AGENTS.md`
- `TRIGGERS.md`
- `doc/README.md`
- `.ai/context/triggers.md`
- `.ai/context/READEME.md`
- `.ai/log/2026-03-14_context-priority-wording-adjustment.md`

## 변경 내용
- 교차 검증 중심 문구를 우선 참조 중심 문구로 변경
- `.ai/context`를 기본 작업 해석/진행 기준으로 사용하는 규칙 반영
- 충돌 또는 해석 불가 시에만 원본 문서로 fallback하는 조건 명시

## 영향 범위
- 문서 기반 작업 시작 절차
- 트리거 해석 절차
- AI context 사용 우선순위

## 비고
- 최종 기준이 원본 문서라는 원칙은 충돌/해석 불가 조건에서 유지한다.
