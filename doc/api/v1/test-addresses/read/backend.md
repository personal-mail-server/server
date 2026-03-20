# 테스트용 메일 주소 목록/상세 조회 API 문서

## 문서 목적
이 문서는 테스트용 메일 주소 목록 조회 API와 상세 조회 API의 요청/응답 계약과 동작 규칙을 정의한다.

---

## 기준 문서
- `doc/domain/mail/server-requirements.md`
- `doc/domain/mail/pr-split/04-read-api.md`
- `doc/domain/mail/test-address-model.md`

---

## 목록 조회 API
- 메서드: `GET`
- 경로: `/api/v1/test-addresses`

### 인증
- Bearer 액세스 토큰이 필요하다.

### 성공 응답
```json
{
  "addresses": [
    {
      "id": 1,
      "ownerUserId": 1,
      "email": "one@mail.local",
      "createdAt": "2026-03-20T00:00:00Z",
      "updatedAt": "2026-03-20T00:00:00Z",
      "deletedAt": null
    }
  ]
}
```

### 동작 규칙
- 현재 로그인한 사용자의 주소만 조회해야 한다.
- 논리 삭제된 주소는 기본 목록에서 제외해야 한다.

---

## 상세 조회 API
- 메서드: `GET`
- 경로: `/api/v1/test-addresses/{id}`

### 인증
- Bearer 액세스 토큰이 필요하다.

### 동작 규칙
- 현재 로그인한 사용자의 주소만 상세 조회할 수 있어야 한다.
- 다른 사용자의 주소이거나 존재하지 않는 주소는 `404` 로 응답해야 한다.

---

## 실패 응답

### 401
- 코드: `INVALID_ACCESS_TOKEN`
- 의미: 액세스 토큰이 없거나 유효하지 않거나 세션 버전이 맞지 않음

### 404
- 코드: `RESOURCE_NOT_FOUND`
- 의미: 해당 주소가 없거나 현재 사용자 소유가 아님

### 500
- 코드: `INTERNAL_SERVER_ERROR`
- 의미: 조회 중 내부 오류 발생

---

## 정리
이 문서는 테스트용 메일 주소 읽기 흐름을 목록 조회와 상세 조회로 분리해 정의하며, 두 API 모두 인증 사용자 범위 안에서만 데이터를 반환해야 한다.
