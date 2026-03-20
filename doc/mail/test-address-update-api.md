# 테스트용 메일 주소 수정 API 문서

## 문서 목적
이 문서는 테스트용 메일 주소 수정 API의 요청/응답 계약과 동작 규칙을 정의한다.

---

## 기준 문서
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/05-update-api.md`
- `doc/mail/test-address-model.md`

---

## 엔드포인트
- 메서드: `PUT`
- 경로: `/api/v1/test-addresses/{id}`

---

## 인증
- Bearer 액세스 토큰이 필요하다.

---

## 요청 본문
```json
{
  "email": "updated@mail.local"
}
```

---

## 동작 규칙
- 현재 로그인한 사용자의 주소만 수정할 수 있어야 한다.
- 자기 자신과 동일한 현재 이메일 값으로는 수정 가능해야 한다.
- 다른 자원이 이미 사용 중인 이메일로는 수정할 수 없어야 한다.
- 존재하지 않거나 타 사용자 소유인 주소는 `404` 로 응답해야 한다.
- 논리 삭제된 주소는 수정할 수 없어야 한다.

---

## 성공 응답
- 상태 코드: `200 OK`
- 응답 본문: `TestAddressResponse`

---

## 실패 응답

### 400
- 코드: `INVALID_REQUEST_BODY`, `MISSING_REQUIRED_FIELD`, `INVALID_EMAIL_FORMAT`

### 401
- 코드: `INVALID_ACCESS_TOKEN`

### 404
- 코드: `RESOURCE_NOT_FOUND`

### 409
- 코드: `DUPLICATE_TEST_ADDRESS_EMAIL`

### 500
- 코드: `INTERNAL_SERVER_ERROR`

---

## 정리
이 API는 인증된 사용자가 자신의 테스트용 메일 주소를 수정하되, 자기 자신 값 유지 허용과 타 자원 중복 차단을 명확히 보장하기 위한 보호 API이다.
