# 작업 로그

## 작업 일시
- 2026-03-20

## 작업 유형
- 설정 수정

## 기능명
- 로컬 PostgreSQL 포트 공개 설정

## 관련 키워드
- PostgreSQL
- Docker Compose
- localhost
- 5432
- DB 접속

## 변경 이유
- 사용자가 로컬 DBMS 클라이언트에서 `localhost:5432` 로 직접 접속하려 했지만, DB 컨테이너 포트가 호스트에 공개되지 않아 접속 거부 오류가 발생했기 때문

## 변경 대상
- `code/docker-compose.yml`
- `doc/project/current-stack.md`
- `.ai/context/project-current-stack.md`
- `.ai/log/2026-03-20_local-db-port-exposure.md`

## 변경 내용
- PostgreSQL 컨테이너에 `5432:5432` 포트 매핑을 추가함
- 현재 구성 문서와 AI 컨텍스트에 데이터베이스 외부 노출 포트를 반영함

## 영향 범위
- 로컬 DB 클라이언트 접속 방식
- Docker Compose 포트 바인딩

## 비고
- 같은 머신에서 이미 5432 포트를 사용하는 다른 PostgreSQL 인스턴스가 있으면 컨테이너 기동 시 충돌할 수 있음
