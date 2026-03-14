# 로그인 API 명세

## 문서 목적
이 문서는 Personal Mail Server의 로그인 기능에 대한 API 계약을 정의한다.

본 문서는 `doc/auth/login.md`의 기능 설계를 바탕으로 작성되며, 로그인 요청과 응답, 검증 규칙, 상태 코드, 에러 규칙을 명확히 하는 것을 목적으로 한다.

---

## 기준 문서
- 기능 설계 기준 문서: `doc/auth/login.md`
- 문서 저장소 기준: `doc/README.md`
- 기술 및 운영 기준: `TECH_STACK.md`

---

## 범위
이 문서는 로그인 API만 다룬다.

포함 범위:
- 로그인 요청 엔드포인트
- 요청 필드와 검증 규칙
- 성공 응답 구조
- 실패 응답 구조
- 계정 잠금 처리 규칙
- 토큰 전달 방식과 만료 시간

제외 범위:
- 로그아웃 API
- 토큰 재발급 API
- 비밀번호 재설정 API
- 회원가입 API
- 이메일 인증 API
- 2차 인증 API
- 내부 구현 구조
- 데이터베이스 스키마

---

## 엔드포인트 정의

### 로그인
- Method: `POST`
- Path: `/api/v1/auth/login`
- 설명: 로그인 ID와 비밀번호를 검증하고, 성공 시 액세스 토큰과 리프레시 토큰을 발급한다.

---

## 요청 명세

### Request Headers
- `Content-Type: application/json`

### Request Body
```json
{
  "loginId": "string",
  "password": "string"
}
```

### Request Field Rules

#### `loginId`
- 타입: 문자열
- 필수 여부: 필수
- 길이: 4자 이상 32자 이하
- 허용 문자: 영문 소문자, 숫자, 하이픈(`-`)
- 공백 허용: 불가
- 대문자 허용: 불가
- 이메일 형식 사용: 불가

#### `password`
- 타입: 문자열
- 필수 여부: 필수
- 길이: 8자 이상 64자 이하
- 포함 규칙: 영문자 1개 이상, 숫자 1개 이상 필수
- 특수문자: 허용
- 공백: 허용하지 않음

### Request Validation Rules
- 요청 본문이 JSON 형식이 아니면 요청은 실패해야 한다.
- `loginId` 또는 `password`가 누락되면 요청은 실패해야 한다.
- `loginId` 형식이 규칙에 맞지 않으면 요청은 실패해야 한다.
- `password` 형식이 규칙에 맞지 않으면 요청은 실패해야 한다.

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
  - 액세스 토큰 문자열
  - 보호된 API 호출에 사용한다.
- `refreshToken`
  - 리프레시 토큰 문자열
  - 후속 토큰 재발급 흐름에 사용한다.
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
- 액세스 토큰과 리프레시 토큰은 모두 응답 본문으로 전달한다.

---

## 실패 응답 명세

### 공통 원칙
- 실패 응답은 내부 구현 정보나 민감한 보안 정보를 노출하지 않아야 한다.
- 인증 실패 응답은 계정 존재 여부를 과도하게 추론할 수 없도록 설계해야 한다.
- 실패 시 액세스 토큰과 리프레시 토큰은 반환하지 않아야 한다.

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
- `INVALID_LOGIN_ID_FORMAT`
- `INVALID_PASSWORD_FORMAT`
- `MISSING_REQUIRED_FIELD`

응답 예시:
```json
{
  "code": "INVALID_LOGIN_ID_FORMAT",
  "message": "입력값 형식이 올바르지 않습니다."
}
```

#### `401 Unauthorized`
인증 실패

에러 코드 예시:
- `INVALID_CREDENTIALS`

응답 예시:
```json
{
  "code": "INVALID_CREDENTIALS",
  "message": "로그인 ID 또는 비밀번호가 올바르지 않습니다."
}
```

#### `423 Locked`
계정 잠금 상태

에러 코드 예시:
- `ACCOUNT_LOCKED`

응답 예시:
```json
{
  "code": "ACCOUNT_LOCKED",
  "message": "계정이 일시적으로 잠겼습니다. 잠시 후 다시 시도해 주세요."
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

## 계정 잠금 처리 규칙
- 로그인 실패 5회째 요청 시 계정은 10분 동안 잠금 상태로 전환되어야 한다.
- 연속 실패 횟수만 누적한다.
- 로그인 성공 시 실패 횟수는 즉시 초기화되어야 한다.
- 5회째 실패 요청과 잠금 기간 중 로그인 요청은 `423 Locked`로 응답해야 한다.
- 잠금 상태에서는 비밀번호가 올바르더라도 토큰을 발급하지 않아야 한다.

---

## 인증 헤더 사용 규칙
- 보호된 API 호출 시 액세스 토큰은 `Authorization` 헤더에 `Bearer <token>` 형식으로 전달하는 것을 기준으로 한다.
- 본 문서는 로그인 API 응답 계약을 다루며, 실제 보호 API 목록은 별도 문서에서 정의한다.

---

## Swagger 반영 기준
- 이 명세의 엔드포인트, 요청 필드, 응답 필드, 상태 코드는 Swagger에 동일하게 반영 가능해야 한다.
- 필드명과 상태 코드는 구현 시 임의로 변경해서는 안 된다.
- 추가 응답 필드가 필요할 경우 기준 문서 갱신 후 반영해야 한다.

---

## 검증 기준

### 정상 케이스
- 유효한 `loginId`와 `password`를 보내면 `200 OK`를 반환해야 한다.
- 성공 시 `accessToken`, `refreshToken`, `accessTokenExpiresIn`, `refreshTokenExpiresIn`, `tokenType`를 모두 반환해야 한다.

### 입력 검증 케이스
- `loginId` 누락 시 `400 Bad Request`를 반환해야 한다.
- `password` 누락 시 `400 Bad Request`를 반환해야 한다.
- `loginId`가 길이 또는 문자 규칙을 위반하면 `400 Bad Request`를 반환해야 한다.
- `password`가 길이, 조합 규칙, 공백 규칙을 위반하면 `400 Bad Request`를 반환해야 한다.

### 인증 실패 케이스
- 존재하지 않는 계정 또는 잘못된 비밀번호는 `401 Unauthorized`를 반환해야 한다.
- 인증 실패 응답은 계정 존재 여부를 노출하지 않아야 한다.

### 잠금 케이스
- 연속 5회 실패째 요청은 `423 Locked`를 반환해야 한다.
- 잠금 시간 10분이 지나면 다시 로그인 시도가 가능해야 한다.

### 범위 제외 확인
- 본 문서에는 로그아웃 API가 포함되지 않아야 한다.
- 본 문서에는 토큰 재발급 API가 포함되지 않아야 한다.
- 본 문서에는 비밀번호 재설정 API가 포함되지 않아야 한다.
- 본 문서에는 이메일 인증 API가 포함되지 않아야 한다.
- 본 문서에는 2차 인증 API가 포함되지 않아야 한다.

---

## 후속 문서화 필요 항목
- 토큰 재발급 API 명세
- 비밀번호 저장 정책
- 감사 로그 상세 명세
- 운영자 계정 잠금 해제 절차

---

## 정리
로그인 API는 `POST /api/v1/auth/login` 단일 엔드포인트를 통해 로그인 ID와 비밀번호를 검증하고, 성공 시 액세스 토큰과 리프레시 토큰을 응답 본문으로 반환한다.

또한 입력 검증, 인증 실패, 계정 잠금 상태를 명확한 상태 코드와 에러 코드로 구분하며, 액세스 토큰 30분, 리프레시 토큰 7일 기준을 사용한다.
