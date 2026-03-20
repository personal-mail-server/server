# 테스트용 메일 주소 논리 삭제 API 문서

## 문서 목적
이 문서는 테스트용 메일 주소 논리 삭제 API의 요청/응답 계약과 동작 규칙을 정의한다.

---

## 기준 문서
- `doc/mail/server-requirements.md`
- `doc/mail/pr-split/06-delete-api.md`
- `doc/mail/test-address-model.md`

---

## 엔드포인트
- 메서드: `DELETE`
- 경로: `/api/v1/test-addresses/{id}`

---

## 인증
- Bearer 액세스 토큰이 필요하다.

---

## 동작 규칙
- 현재 로그인한 사용자의 주소만 삭제할 수 있어야 한다.
- 삭제는 물리 삭제가 아니라 `deleted_at` 기록 방식이어야 한다.
- 비소유, 미존재, 이미 삭제된 자원은 모두 `404` 로 응답해야 한다.
- 삭제 후 기본 목록 조회에서는 더 이상 보이면 안 된다.

---

## 성공 응답
- 상태 코드: `204 No Content`
- 응답 본문: 없음

---

## 실패 응답

### 401
- 코드: `INVALID_ACCESS_TOKEN`

### 404
- 코드: `RESOURCE_NOT_FOUND`

### 500
- 코드: `INTERNAL_SERVER_ERROR`

---

## 정리
이 API는 테스트용 메일 주소를 즉시 물리 삭제하지 않고, 인증 사용자 범위 안에서만 논리 삭제하기 위한 보호 API이다.
