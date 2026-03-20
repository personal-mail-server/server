# 테스트용 메일 주소 관리 PR 분할 계획

## 문서 목적
이 문서는 `doc/domain/mail/server-requirements.md` 기준 기능을 여러 개의 독립적인 PR로 분할하기 위한 계획 문서이다.

각 PR은 가능한 한 단독 리뷰와 단독 머지가 가능해야 하며, 다음 PR의 기반이 되는 최소 변경만 포함해야 한다.

---

## 기준 문서
- `doc/domain/mail/server-requirements.md`
- `doc/README.md`
- `TECH_STACK.md`

---

## 분할 원칙
- 하나의 PR은 하나의 책임만 가진다.
- DB 기반, API, 프론트엔드를 한 PR에 과도하게 섞지 않는다.
- 생성/수정/삭제는 각각 독립 검토가 가능하면 분리한다.
- 논리 삭제 정책은 삭제 PR에서만 확정한다.
- 각 PR은 대응 문서, OpenAPI, 테스트를 자기 범위 안에서 함께 갱신한다.

---

## PR 목록
- `01-storage-foundation.md` - 데이터 모델과 저장 기반
- `02-generate-api.md` - 유니크 메일 후보 생성 API
- `03-create-api.md` - 테스트용 메일 주소 생성 API
- `04-read-api.md` - 목록/상세 조회 API
- `05-update-api.md` - 수정 API
- `06-delete-api.md` - 논리 삭제 API
- `07-frontend-read-create.md` - 프론트엔드 조회/생성
- `08-frontend-update-delete.md` - 프론트엔드 수정/삭제

---

## 권장 순서
1. `01-storage-foundation.md`
2. `02-generate-api.md`
3. `03-create-api.md`
4. `04-read-api.md`
5. `05-update-api.md`
6. `06-delete-api.md`
7. `07-frontend-read-create.md`
8. `08-frontend-update-delete.md`

---

## 정리
이 분할안은 테스트용 메일 주소 관리 기능을 작은 PR 단위로 쪼개어, 저장 기반부터 API와 프론트엔드까지 순차적으로 안전하게 머지하기 위한 계획이다.
