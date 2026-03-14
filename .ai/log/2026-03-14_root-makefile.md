# 작업 로그

## 작업 일시
- 2026-03-14

## 작업 유형
- 실행 도구 추가

## 기능명
- 루트 Makefile 기반 수동 실행 명령 추가

## 관련 키워드
- Makefile
- make up
- Docker Compose
- 실행 명령
- 운영 편의성

## 변경 이유
- 프로젝트 루트에서 `make up` 같은 명령으로 전체 스택을 수동 실행하고 상태를 확인할 수 있는 진입점이 필요했기 때문

## 변경 대상
- `Makefile`
- `README.md`
- `code/README.md`
- `doc/project/current-stack.md`
- `doc/project/current-testing.md`
- `.ai/context/project-current-stack.md`
- `.ai/context/project-current-testing.md`
- `.ai/log/2026-03-14_root-makefile.md`

## 변경 내용
- 루트 `Makefile`에 `up`, `down`, `status`, `ps`, `logs` 명령 추가
- 루트 README와 실행 문서, 현재 구성/테스트 문서에 루트 Make 사용 흐름 반영
- 관련 AI context 문서에 루트 실행 진입점 정보 반영

## 영향 범위
- 프로젝트 수동 실행 방법
- Docker 기반 운영 진입점
- 현재 구성 및 테스트 문서

## 비고
- 실제 서비스 오케스트레이션 소스는 계속 `code/docker-compose.yml`이며, `Makefile`은 루트 편의 래퍼 역할만 수행함
