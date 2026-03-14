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

## 확인 포인트
- 로그인 화면: `http://localhost:3000`
- API: `http://localhost:18080/api/v1/auth/login`
- OpenAPI: `http://localhost:3000/docs/openapi.yaml`
- Swagger UI: `http://localhost:3000/docs`

## 기본 테스트 계정
- 로그인 ID: `user-01`
- 비밀번호: `pass1234`

위 계정은 Docker 초기화 시 자동으로 주입된다.
