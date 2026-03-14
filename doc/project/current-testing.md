# 현재 테스트 문서

## 문서 목적
이 문서는 현재 저장소에 실제로 구현되어 있는 테스트와 검증 흐름을 정리한다.

본 문서는 테스트 원칙 자체를 정의하는 문서가 아니라, AI가 현재 어떤 테스트를 갖고 있고 어떤 방식으로 검증하고 있는지를 기록하는 현재 상태 문서이다.

---

## 기준 문서
- 최소 기술 기준: `TECH_STACK.md`
- 현재 구성 문서: `doc/project/current-stack.md`
- 현재 구현 기준 파일:
  - `code/internal/auth/service_test.go`
  - `code/internal/auth/validation_test.go`
  - `code/internal/http/handlers/auth_handler_test.go`
  - `code/README.md`

---

## 현재 테스트 구성 요약
- 언어 단위 테스트: Go `testing` 패키지 기반
- 현재 단위 테스트 대상:
  - 인증 서비스 로직
  - 입력 검증 로직
  - 로그인 HTTP 핸들러 일부
- 현재 실행 검증 대상:
  - `go test ./...`
  - `go build ./cmd/server`
  - `go vet ./...`
  - Docker Compose 기동 확인
  - 로그인 화면 HTTP 응답 확인
  - Swagger 접근 확인
  - 로그인 API 성공/실패/잠금 흐름 확인

---

## 단위 테스트

### 인증 서비스 테스트
위치:
- `code/internal/auth/service_test.go`

현재 검증 항목:
- 로그인 성공 시 실패 횟수 초기화
- 로그인 성공 시 액세스 토큰/리프레시 토큰 발급 메타데이터 확인
- 잘못된 비밀번호 반복 입력 시 실패 횟수 누적
- 5회째 실패 시 `ACCOUNT_LOCKED` 전환
- 잠금 상태에서는 로그인 차단
- 잠금 만료 후 다시 로그인 가능

### 입력 검증 테스트
위치:
- `code/internal/auth/validation_test.go`

현재 검증 항목:
- 로그인 ID 형식 검증
- 비밀번호 형식 검증
- 길이 제한 검증
- 허용 문자 및 공백 규칙 검증

### HTTP 핸들러 테스트
위치:
- `code/internal/http/handlers/auth_handler_test.go`

현재 검증 항목:
- 잘못된 JSON 요청 시 `400 Bad Request`
- 오류 응답 코드 매핑 확인

---

## 통합 성격 검증
현재는 별도의 통합 테스트 프레임워크가 아니라, 실행 기반 검증과 API 호출 검증 조합으로 통합 성격의 확인을 수행하고 있다.

현재 검증 항목:
- Docker Compose로 DB, 백엔드, 프론트엔드가 함께 기동되는지 확인
- 프론트엔드 루트 경로(`/`) 응답 확인
- `/docs` 경로 응답 확인
- `POST /api/v1/auth/login` 성공 응답 확인
- 잘못된 비밀번호 반복 시 잠금 정책이 실제 응답 코드에 반영되는지 확인

---

## 실행 명령

### 코드 레벨 검증
```bash
go test ./...
go build ./cmd/server
go vet ./...
```

### Docker 실행 검증
```bash
docker compose -f code/docker-compose.yml up --build -d
```

### 주요 확인 지점
- 로그인 화면: `http://localhost:3000`
- 백엔드 직접 호출: `http://localhost:18080/api/v1/auth/login`
- Swagger UI: `http://localhost:3000/docs`
- OpenAPI YAML: `http://localhost:3000/docs/openapi.yaml`

---

## 수동 검증 흐름

### 성공 로그인 확인
- 테스트 계정 `user-01` / `pass1234`로 로그인 요청
- `200 OK` 응답과 토큰 필드 확인

### 실패 로그인 확인
- 잘못된 비밀번호로 요청
- `401 Unauthorized` 응답 확인

### 잠금 정책 확인
- 잘못된 비밀번호를 반복 전송
- 5회째 실패 시 `423 Locked` 전환 확인
- 잠금 상태 중 로그인 차단 확인

### 화면 검증
- 브라우저에서 로그인 폼 표시 확인
- Swagger 링크 접근 확인
- 상태 영역에 응답 JSON 노출 확인

---

## 현재 미구현 테스트 항목
- 프론트엔드 전용 자동화 테스트
- 실제 브라우저 자동화 기반 E2E 테스트
- DB 격리 기반 통합 테스트 스위트
- 로그인 이후 보호 화면 흐름 테스트
- 토큰 재발급 흐름 테스트
- 성능 테스트

---

## 갱신 원칙
- 테스트 코드가 추가, 수정, 삭제되면 본 문서도 함께 갱신해야 한다.
- 검증 명령이나 실행 흐름이 바뀌면 본 문서도 함께 갱신해야 한다.
- 수동 검증 방식이 자동화로 전환되면 본 문서의 테스트 분류도 함께 갱신해야 한다.

---

## 정리
현재 프로젝트는 Go 단위 테스트와 Docker 기반 실행 검증을 함께 사용해 로그인 슬라이스를 검증하고 있다.

본 문서는 AI가 실제로 어떤 테스트를 작성하고 어떤 검증을 수행했는지 추적하기 위한 현재 상태 문서이며, 앞으로 테스트 체계가 바뀌면 지속적으로 함께 갱신되어야 한다.
