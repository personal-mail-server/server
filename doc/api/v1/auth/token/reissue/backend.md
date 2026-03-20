# 토큰 재발급 API 명세

## 문서 목적
이 문서는 Personal Mail Server의 토큰 재발급 기능에 대한 API 계약을 정의한다.

본 문서는 `doc/api/v1/auth/token/reissue/design.md`의 기능 설계를 바탕으로 작성되며, 토큰 재발급 요청과 응답, 검증 규칙, 상태 코드, 에러 규칙을 명확히 하는 것을 목적으로 한다.

---

## 기준 문서
- 기능 설계 기준 문서: `doc/api/v1/auth/token/reissue/design.md`
- 로그인 API 기준 문서: `doc/api/v1/auth/login/backend.md`
- 로그아웃 API 기준 문서: `doc/api/v1/auth/logout/backend.md`
- 문서 저장소 기준: `doc/README.md`
- 기술 및 운영 기준: `TECH_STACK.md`

---

## 범위
이 문서는 토큰 재발급 API만 다룬다.

포함 범위:
- 토큰 재발급 요청 엔드포인트
- 요청 본문 구조
- 성공 응답 구조
- 실패 응답 구조
- 리프레시 토큰 검증 규칙
- 세션 버전 불일치 차단 규칙

제외 범위:
- 로그인 API
- 로그아웃 API
- 회원가입 API
- 비밀번호 재설정 API
- 내부 저장소 구현 구조
- 브라우저 저장 전략

---

## 엔드포인트 정의

### 토큰 재발급
- Method: `POST`
- Path: `/api/v1/auth/token/reissue`
- 설명: 유효한 리프레시 토큰을 검증하고, 성공 시 새 액세스 토큰과 새 리프레시 토큰을 발급한다.

---

## 요청 명세

### Request Headers
- `Content-Type: application/json`

### Request Body
```json
{
  "refreshToken": "string"
}
```

### Request Field Rules

#### `refreshToken`
- 타입: 문자열
- 필수 여부: 필수
- 의미: 로그인 성공 시 발급된 리프레시 토큰
- 공백만 있는 값: 허용하지 않음

### Request Validation Rules
- 요청 본문이 JSON 형식이 아니면 요청은 실패해야 한다.
- `refreshToken`이 누락되면 요청은 실패해야 한다.
- `refreshToken`이 문자열이 아니면 요청은 실패해야 한다.
- 공백만 있는 `refreshToken`은 요청 실패로 처리해야 한다.
- 액세스 토큰이나 기타 토큰을 `refreshToken` 자리에 제출하면 요청은 실패해야 한다.

---

## 성공 응답 명세

### Success Status
- `200 OK`

### Success Body
```json
{
  "accessToken": "string",
  "refreshToken": "string",
  "accessTokenExpiresIn": 1800,
  "refreshTokenExpiresIn": 604800,
  "tokenType": "Bearer"
}
```

### Success Field Rules
- `accessToken`
  - 새로 발급된 액세스 토큰 문자열
  - 보호된 API 호출에 사용한다.
- `refreshToken`
  - 새로 발급된 리프레시 토큰 문자열
  - 후속 재발급 흐름에 사용한다.
- `accessTokenExpiresIn`
  - 초 단위 만료 시간
  - 값은 `1800`이어야 한다.
- `refreshTokenExpiresIn`
  - 초 단위 만료 시간
  - 값은 `604800`이어야 한다.
- `tokenType`
  - 인증 스키마 식별값
  - 값은 `Bearer`여야 한다.

### Token Policy
- 액세스 토큰 만료 시간은 30분이다.
- 리프레시 토큰 만료 시간은 7일이다.
- 재발급 성공 시 새 액세스 토큰과 새 리프레시 토큰을 모두 응답 본문으로 전달한다.
- 재발급에 성공한 기존 리프레시 토큰은 즉시 무효화되어야 한다.

---

## 실패 응답 명세

### 공통 원칙
- 실패 응답은 내부 구현 정보나 민감한 보안 정보를 노출하지 않아야 한다.
- 실패 시 새 액세스 토큰과 새 리프레시 토큰은 반환하지 않아야 한다.
- 실패 응답은 계정 존재 여부나 세션 저장 구조를 과도하게 노출하지 않아야 한다.

### Error Response Body
```json
{
  "code": "string",
  "message": "string"
}
```

### 상태 코드 및 에러 코드

#### `400 Bad Request`
요청 형식 또는 입력값 검증 실패

에러 코드 예시:
- `INVALID_REQUEST_BODY`
- `MISSING_REQUIRED_FIELD`

응답 예시:
```json
{
  "code": "MISSING_REQUIRED_FIELD",
  "message": "입력값 형식이 올바르지 않습니다."
}
```

#### `401 Unauthorized`
리프레시 토큰 검증 실패 또는 무효 세션

에러 코드 예시:
- `INVALID_REFRESH_TOKEN`

응답 예시:
```json
{
  "code": "INVALID_REFRESH_TOKEN",
  "message": "인증 정보가 유효하지 않습니다. 다시 로그인해 주세요."
}
```

#### `500 Internal Server Error`
서버 내부 오류

에러 코드 예시:
- `INTERNAL_SERVER_ERROR`

응답 예시:
```json
{
  "code": "INTERNAL_SERVER_ERROR",
  "message": "요청을 처리할 수 없습니다. 잠시 후 다시 시도해 주세요."
}
```

---

## 리프레시 토큰 검증 규칙
- 서버는 제출된 토큰이 리프레시 토큰 전용 서명과 클레임 구조를 만족하는지 검증해야 한다.
- 토큰의 용도 값은 `refresh`여야 한다.
- 토큰이 만료되었거나 서명이 올바르지 않으면 `401 Unauthorized`를 반환해야 한다.
- 토큰의 사용자 식별자와 세션 버전을 기준으로 현재 계정 상태를 다시 확인해야 한다.
- 현재 계정의 세션 버전이 토큰의 세션 버전과 다르면 `401 Unauthorized`를 반환해야 한다.
- 이미 재발급에 한 번 사용된 리프레시 토큰이면 `401 Unauthorized`를 반환해야 한다.

---

## 인증 및 상태 전이 규칙
- 재발급 요청은 `Authorization` 헤더를 사용하지 않는다.
- 재발급 요청은 요청 본문의 `refreshToken`만으로 처리한다.
- 로그아웃 성공 후 기존 리프레시 토큰은 재발급에 사용할 수 없어야 한다.
- 재발급 성공 후 직전에 사용한 리프레시 토큰은 다시 사용할 수 없어야 한다.
- 재발급 성공 후 클라이언트는 새 액세스 토큰과 새 리프레시 토큰으로 인증 상태를 교체해야 한다.

---

## Swagger 반영 기준
- 이 명세의 엔드포인트, 요청 필드, 응답 필드, 상태 코드는 Swagger에 동일하게 반영 가능해야 한다.
- 필드명과 상태 코드는 구현 시 임의로 변경해서는 안 된다.
- 추가 응답 필드가 필요할 경우 기준 문서 갱신 후 반영해야 한다.

---

## 검증 기준

### 정상 케이스
- 유효한 `refreshToken`을 보내면 `200 OK`를 반환해야 한다.
- 성공 시 `accessToken`, `refreshToken`, `accessTokenExpiresIn`, `refreshTokenExpiresIn`, `tokenType`를 모두 반환해야 한다.
- 성공 직후 동일한 `refreshToken`으로 다시 요청하면 `401 Unauthorized`를 반환해야 한다.

### 입력 검증 케이스
- `refreshToken` 누락 시 `400 Bad Request`를 반환해야 한다.
- 요청 본문이 잘못되면 `400 Bad Request`를 반환해야 한다.

### 인증 실패 케이스
- 만료된 리프레시 토큰이면 `401 Unauthorized`를 반환해야 한다.
- 위조된 리프레시 토큰이면 `401 Unauthorized`를 반환해야 한다.
- 액세스 토큰을 제출하면 `401 Unauthorized`를 반환해야 한다.
- 로그아웃 이후 세션 버전이 달라진 토큰이면 `401 Unauthorized`를 반환해야 한다.
- 이미 사용 완료된 리프레시 토큰이면 `401 Unauthorized`를 반환해야 한다.

### 범위 제외 확인
- 본 문서에는 로그인 API가 포함되지 않아야 한다.
- 본 문서에는 로그아웃 API가 포함되지 않아야 한다.
- 본 문서에는 비밀번호 재설정 API가 포함되지 않아야 한다.

---

## 후속 문서화 필요 항목
- 보호 API 인증 실패 규칙
- 리프레시 토큰 저장 및 폐기 정책
- 감사 로그 상세 명세

---

## 정리
토큰 재발급 API는 `POST /api/v1/auth/token/reissue` 단일 엔드포인트를 통해 `refreshToken`을 받아 새 액세스 토큰과 새 리프레시 토큰을 응답 본문으로 반환한다.

또한 요청은 리프레시 토큰 전용 검증, 세션 버전 확인, 사용 완료 토큰 재사용 차단을 통과해야 하며, 로그아웃이나 세션 무효화 이후에는 반드시 `401 Unauthorized`로 차단되어야 한다.
