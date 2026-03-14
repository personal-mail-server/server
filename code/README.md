# Backend Login Sync Slice

## 실행

```bash
docker compose -f code/docker-compose.yml up --build
```

## 확인 포인트
- 로그인 화면: `http://localhost:3000`
- API: `http://localhost:18080/api/v1/auth/login`
- OpenAPI: `http://localhost:3000/docs/openapi.yaml`
- Swagger UI: `http://localhost:3000/docs`

## 기본 테스트 계정
- 로그인 ID: `user-01`
- 비밀번호: `pass1234`

위 계정은 Docker 초기화 시 자동으로 주입된다.
