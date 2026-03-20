# 테스트용 메일 주소 생성 API 문서

## 문서 목적
이 문서는 테스트용 메일 주소 생성 API의 요청/응답 계약과 동작 규칙을 정의한다.

---

## 기준 문서
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/03-create-api.md`
- `doc/mail/test-address-model.md`
- `doc/mail/test-address-generate-api.md`

---

## 엔드포인트
- 메서드: `POST`
- 경로: `/api/v1/test-addresses`

---

## 인증
- Bearer 액세스 토큰이 필요하다.
- 인증되지 않았거나 세션 버전이 맞지 않으면 실패해야 한다.

---

## 요청 본문
```json
{
  "email": "created@mail.local"
}
```

## 요청 필드
- `email`: 사용자가 직접 입력하거나 후보 생성 API 응답을 복사해 넣는 테스트용 메일 주소

---

## 성공 응답
```json
{
  "id": 1,
  "ownerUserId": 1,
  "email": "created@mail.local",
  "createdAt": "2026-03-20T00:00:00Z",
  "updatedAt": "2026-03-20T00:00:00Z",
  "deletedAt": null
}
```

---

## 동작 규칙
- 서버는 생성 전에 입력 이메일 형식을 검증해야 한다.
- 서버는 저장 직전 동일 이메일 존재 여부를 확인해야 한다.
- 최종 저장 단계에서도 DB 유니크 제약으로 중복을 방지해야 한다.
- 생성된 주소는 인증된 사용자 소유로 저장되어야 한다.
- 현재 구현은 이메일 정규화 없이 입력값 그대로 저장한다.

---

## 실패 응답

### 400
- 코드: `INVALID_REQUEST_BODY`
- 의미: JSON 형식이 잘못됨

### 400
- 코드: `MISSING_REQUIRED_FIELD`
- 의미: `email` 이 비어 있음

### 400
- 코드: `INVALID_EMAIL_FORMAT`
- 의미: 이메일 형식이 유효하지 않음

### 401
- 코드: `INVALID_ACCESS_TOKEN`
- 의미: 액세스 토큰이 없거나 유효하지 않거나 세션 버전이 맞지 않음

### 409
- 코드: `DUPLICATE_TEST_ADDRESS_EMAIL`
- 의미: 이미 사용 중인 메일 주소임

### 500
- 코드: `INTERNAL_SERVER_ERROR`
- 의미: 저장 또는 조회 중 내부 오류 발생

---

## 제외 범위
- 목록/상세 조회
- 수정/삭제
- 프론트엔드 UI

---

## 정리
이 API는 인증된 사용자가 직접 입력하거나 후보 API에서 받은 이메일을 사용해 테스트용 메일 주소를 실제로 저장하기 위한 보호 API이다.
