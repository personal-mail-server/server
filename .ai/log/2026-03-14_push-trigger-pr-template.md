# 작업 로그

## 작업 일시
- 2026-03-14

## 작업 유형
- 자동화 수정

## 기능명
- PUSH 트리거 PR 작성 형식 한국어화

## 관련 키워드
- PUSH
- pull request
- 한국어
- 템플릿
- 자동화

## 변경 이유
- 자동 생성되는 PR 내용을 사람이 읽기 쉬운 한국어 형식으로 정리할 필요가 있었기 때문

## 변경 대상
- `code/internal/automation/pushtrigger/service.go`
- `code/internal/automation/pushtrigger/service_test.go`
- `code/README.md`
- `.ai/context/triggers.md`
- `.ai/log/2026-03-14_push-trigger-pr-template.md`

## 변경 내용
- PUSH 트리거의 PR 본문 섹션을 `## 요약`, `## 기준 정보` 한국어 템플릿으로 변경
- 대상 브랜치와 작업 브랜치를 한국어 문구로 함께 표시하도록 수정
- 관련 테스트와 문서, AI context 동기화

## 영향 범위
- 자동 생성 PR 가독성
- PUSH 트리거 결과물 형식

## 비고
- PR 제목은 최신 커밋 제목을 기본으로 유지하고, 본문 중심으로 한국어 보기 형식을 적용함
