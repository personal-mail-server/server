# 작업 로그

## 작업 일시
- 2026-03-14

## 작업 유형
- CI/CD 수정

## 기능명
- GitHub Actions heredoc 들여쓰기 오류 수정

## 관련 키워드
- GitHub Actions
- CI
- heredoc
- Python
- 들여쓰기
- docker-smoke

## 변경 이유
- `main` 브랜치 GitHub Actions에서 Python heredoc 종료가 깨지며 `IndentationError`와 here-document 경고가 발생해 CI가 실패했기 때문

## 변경 대상
- `.github/workflows/ci.yml`
- `.ai/log/2026-03-14_ci-heredoc-fix.md`

## 변경 내용
- `Check successful login` 단계의 Python heredoc 블록에서 잘못 들어간 추가 공백 제거
- 종료 토큰 `PY`가 정상적으로 인식되도록 정렬 복구
- 동일 블록을 로컬 셸에서 재실행해 문법 오류가 사라졌는지 확인

## 영향 범위
- GitHub Actions `docker-smoke` 잡
- 로그인 성공 검증 단계

## 비고
- 문제 원인은 `.github/workflows/ci.yml`의 assert 두 줄과 heredoc 종료 줄 들여쓰기 불일치였다.
