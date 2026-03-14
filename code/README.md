# Backend Login Sync Slice

## 실행

```bash
make up
```

직접 Docker Compose를 사용할 경우에는 아래 명령을 사용할 수 있다.

```bash
docker compose -f code/docker-compose.yml up --build -d
```

## 루트 Make 명령
- `make up` - DB, 백엔드, 프론트엔드 전체 스택 기동
- `make down` - 전체 스택 종료
- `make status` - 현재 컨테이너 상태 확인
- `make logs` - 전체 스택 로그 확인
- `make push-trigger` - 현재 브랜치를 push하고 `main` 대상 PR 생성 또는 재사용 후 auto-merge 예약

## PUSH 자동화
- 예약 트리거 `PUSH`의 로컬 실행 진입점은 `make push-trigger`이다.
- 이 명령은 현재 브랜치가 `main` 또는 `master`가 아니고 작업 트리가 깨끗할 때만 동작한다.
- 실행 전 필요 조건:
  - `git` 사용 가능
  - `gh` CLI 설치 및 로그인 완료
  - `origin` remote 설정 완료
  - 저장소에서 auto-merge 기능 사용 가능
- 현재 버전은 Copilot 리뷰 자동 수집/해석은 포함하지 않는다.
- PR 본문은 한국어 기준의 사람이 읽기 쉬운 형식으로 자동 작성한다.

## 확인 포인트
- 로그인 화면: `http://localhost:3000`
- API: `http://localhost:18080/api/v1/auth/login`
- OpenAPI: `http://localhost:3000/docs/openapi.yaml`
- Swagger UI: `http://localhost:3000/docs`

## 기본 테스트 계정
- 로그인 ID: `user-01`
- 비밀번호: `pass1234`

위 계정은 Docker 초기화 시 자동으로 주입된다.
