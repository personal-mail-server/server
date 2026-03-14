# 현재 프로젝트 구성 문서

## 문서 목적
이 문서는 현재 저장소에 실제로 반영되어 있는 프로젝트 구성을 정리한다.

`TECH_STACK.md`가 사람이 정한 최소 기술 기준이라면, 본 문서는 그 기준 위에서 AI가 실제로 선택하고 구현한 현재 구성을 설명하는 문서이다.

---

## 기준 문서
- 최소 기술 기준: `TECH_STACK.md`
- 프로젝트 개요: `README.md`
- 문서 저장소 원칙: `doc/README.md`
- 현재 구현 기준 파일:
  - `code/go.mod`
  - `code/Dockerfile`
  - `code/docker-compose.yml`
  - `code/frontend/Dockerfile`
  - `code/openapi/openapi.yaml`
  - `code/internal/app/server.go`

---

## 현재 구성 요약
- 백엔드 언어: Go
- 백엔드 프레임워크: Echo Framework
- 데이터베이스: PostgreSQL 16 (Docker image: `postgres:16-alpine`)
- 프론트엔드: 정적 HTML/CSS/JavaScript + nginx
- API 계약 문서: OpenAPI 3.0.3 정적 YAML
- 실행 방식: Docker Compose 기반 다중 서비스 실행
- 루트 수동 실행 명령: `Makefile` 기반
- CI/CD 검증: GitHub Actions 기반
- PR 자동화 진입점: `PUSH` 트리거 및 `make push-trigger`

---

## 백엔드 구성
- 엔트리포인트: `code/cmd/server/main.go`
- 애플리케이션 초기화: `code/internal/app/server.go`
- 라우팅: Echo 라우터 기반
- 미들웨어:
  - recover
  - request id
  - secure headers
  - body limit
- 현재 구현 범위:
  - 로그인 API
  - Swagger/OpenAPI 파일 노출
  - 헬스 체크 엔드포인트

현재 백엔드는 `Go + Echo` 조합을 사용하며, 단일 프로세스에서 HTTP API와 OpenAPI 문서 노출을 담당한다.

---

## 데이터베이스 구성
- 데이터베이스 종류: PostgreSQL
- 버전 계열: 16-alpine 이미지 사용
- 연결 방식: Docker Compose 네트워크 내부 연결
- 접속 문자열 기본값: `postgres://postgres:postgres@db:5432/mail_server?sslmode=disable`

현재 로그인 슬라이스 기준 DB 역할:
- 사용자 계정 저장
- 비밀번호 해시 저장
- 로그인 실패 횟수 저장
- 계정 잠금 시간 저장
- 마이그레이션 버전 관리

---

## 프론트엔드 구성
- 제공 방식: 정적 파일 기반 프론트엔드
- 런타임 컨테이너: nginx
- 진입 화면: 로그인 단일 페이지
- 현재 기술 선택:
  - `index.html`
  - `styles.css`
  - `app.js`
  - nginx reverse proxy 설정

현재 프론트엔드는 SPA 프레임워크를 사용하지 않고, 로그인 단일 화면을 빠르게 검증 가능하게 만들기 위한 경량 구조를 사용한다.

---

## API 문서 구성
- 형식: OpenAPI 3.0.3
- 저장 위치: `code/openapi/openapi.yaml`
- 노출 경로:
  - `/docs/openapi.yaml`
  - `/docs`

현재는 코드 주석 기반 생성기가 아니라 정적 OpenAPI YAML을 기준으로 계약을 관리한다.

---

## Docker 구성
- 오케스트레이션: `code/docker-compose.yml`
- 백엔드 빌드: `code/Dockerfile`
- 프론트엔드 빌드: `code/frontend/Dockerfile`
- 루트 실행 진입점: `Makefile`

현재 서비스 구성:
- `db`: PostgreSQL
- `backend`: Go + Echo 서버
- `frontend`: nginx 정적 프론트엔드

현재 기본 실행 명령:
- `make up`
- `make down`
- `make status`
- `make logs`
- `make push-trigger`

현재 CI 진입점:
- `.github/workflows/ci.yml`

현재 PR 자동화 진입점:
- `code/cmd/push-trigger/main.go`
- `make push-trigger`

현재 외부 노출 포트:
- 프론트엔드: `3000`
- 백엔드 직접 접근: `18080`

---

## 인증 구현 상태
- 로그인 방식: 로그인 ID + 비밀번호
- 토큰 방식: 액세스 토큰 + 리프레시 토큰
- 액세스 토큰 만료: 30분
- 리프레시 토큰 만료: 7일
- 로그인 실패 정책: 연속 5회 실패 시 10분 잠금

현재는 로그인 API와 로그인 화면만 구현되어 있으며, 로그아웃/재발급/회원가입/비밀번호 재설정/2차 인증은 포함되지 않는다.

---

## 디렉토리 구성 요약
- `code/cmd/server/` - 서버 엔트리포인트
- `code/internal/` - 내부 애플리케이션 로직
- `code/migrations/` - DB 마이그레이션
- `code/openapi/` - OpenAPI 계약 파일
- `code/frontend/` - 로그인 프론트엔드 정적 자산 및 nginx 설정
- `doc/auth/` - 로그인 관련 기준 문서
- `doc/project/` - 프로젝트 구성/운영 관점 기준 문서

---

## 갱신 원칙
- 본 문서는 현재 실제 구성 기준 문서이다.
- 설계 문서, 구현 코드, Docker 구성, 실행 방식, 프레임워크 선택, DB 선택, API 문서 방식이 변경되면 함께 갱신해야 한다.
- 최소 기준이 바뀌는 경우는 `TECH_STACK.md`를 수정하고, 실제 구현 구성이 바뀌는 경우는 본 문서를 수정한다.

---

## 정리
현재 프로젝트는 `Go + Echo + PostgreSQL + 정적 프론트엔드(nginx) + OpenAPI YAML + Docker Compose` 조합으로 로그인 슬라이스가 구현되어 있다.

본 문서는 AI가 실제로 선택한 현재 구성을 설명하는 문서이며, 앞으로 프로젝트 구성이 바뀔 때 지속적으로 함께 갱신되어야 한다.
