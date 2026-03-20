package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"personal-mail-server/internal/auth"
	"personal-mail-server/internal/testaddress"
)

type testAddressUserReader struct {
	user *auth.User
}

func (r testAddressUserReader) FindByLoginID(_ context.Context, loginID string) (*auth.User, error) {
	if r.user == nil || r.user.LoginID != loginID {
		return nil, auth.ErrUserNotFound
	}
	copy := *r.user
	return &copy, nil
}

type testAddressTokenIssuer struct {
	claims *auth.AuthTokenClaims
	err    error
}

func (t testAddressTokenIssuer) IssuePair(_ time.Time, _ string, _ int, _ string) (*auth.IssuedTokenPair, error) {
	panic("not used")
}

func (t testAddressTokenIssuer) VerifyAccessToken(_ string) (*auth.AuthTokenClaims, error) {
	if t.err != nil {
		return nil, t.err
	}
	return t.claims, nil
}

func (t testAddressTokenIssuer) VerifyRefreshToken(_ string) (*auth.AuthTokenClaims, error) {
	panic("not used")
}

type testAddressRepo struct {
	used map[string]testaddress.TestMailAddress
}

func (r *testAddressRepo) Create(_ context.Context, address testaddress.TestMailAddress) (*testaddress.TestMailAddress, error) {
	if r.used == nil {
		r.used = map[string]testaddress.TestMailAddress{}
	}
	if _, ok := r.used[address.Email]; ok {
		return nil, testaddress.ErrDuplicateEmail
	}
	created := testaddress.TestMailAddress{
		ID:          int64(len(r.used) + 1),
		OwnerUserID: address.OwnerUserID,
		Email:       address.Email,
		CreatedAt:   time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
	}
	if r.used == nil {
		r.used = map[string]testaddress.TestMailAddress{}
	}
	if _, ok := r.used[address.Email]; ok {
		return nil, testaddress.ErrDuplicateEmail
	}
	r.used[address.Email] = created
	return &created, nil
}
func (r *testAddressRepo) GetByID(context.Context, int64) (*testaddress.TestMailAddress, error) {
	panic("not used")
}
func (r *testAddressRepo) ListByOwner(context.Context, int64) ([]testaddress.TestMailAddress, error) {
	panic("not used")
}
func (r *testAddressRepo) Update(context.Context, testaddress.TestMailAddress) (*testaddress.TestMailAddress, error) {
	panic("not used")
}
func (r *testAddressRepo) SoftDelete(context.Context, int64, time.Time) error { panic("not used") }
func (r *testAddressRepo) GetByEmail(_ context.Context, email string) (*testaddress.TestMailAddress, error) {
	if address, ok := r.used[email]; ok {
		copy := address
		return &copy, nil
	}
	return nil, testaddress.ErrTestMailAddressNotFound
}

type fixedCandidateGenerator struct{ email string }

func (g fixedCandidateGenerator) Next() (string, error) { return g.email, nil }

type failingTokenIssuer struct{}

func (failingTokenIssuer) IssuePair(_ time.Time, _ string, _ int, _ string) (*auth.IssuedTokenPair, error) {
	panic("not used")
}

func (failingTokenIssuer) VerifyAccessToken(_ string) (*auth.AuthTokenClaims, error) {
	return nil, auth.ErrInvalidToken
}

func (failingTokenIssuer) VerifyRefreshToken(_ string) (*auth.AuthTokenClaims, error) {
	panic("not used")
}

func TestTestAddressHandlerGenerateCandidateMissingAuthorization(t *testing.T) {
	e := echo.New()
	service := testaddress.NewService(&testAddressRepo{used: map[string]testaddress.TestMailAddress{}}, testAddressUserReader{}, testAddressTokenIssuer{})
	h := NewTestAddressHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/test-addresses/generate", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.GenerateCandidate(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestTestAddressHandlerGenerateCandidateSuccess(t *testing.T) {
	e := echo.New()
	service := testaddress.NewService(&testAddressRepo{used: map[string]testaddress.TestMailAddress{}}, testAddressUserReader{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}, testAddressTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}})
	h := NewTestAddressHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/test-addresses/generate", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid-access-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.GenerateCandidate(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body testaddress.GenerateCandidateResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Email == "" {
		t.Fatalf("expected generated email response")
	}
}

func TestTestAddressHandlerCreateInvalidJSON(t *testing.T) {
	e := echo.New()
	service := testaddress.NewService(&testAddressRepo{used: map[string]testaddress.TestMailAddress{}}, testAddressUserReader{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}, testAddressTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}})
	h := NewTestAddressHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/test-addresses", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid-access-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Create(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestTestAddressHandlerCreateSuccess(t *testing.T) {
	e := echo.New()
	service := testaddress.NewService(&testAddressRepo{used: map[string]testaddress.TestMailAddress{}}, testAddressUserReader{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}, testAddressTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}})
	h := NewTestAddressHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/test-addresses", strings.NewReader(`{"email":"new@mail.local"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid-access-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Create(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	var body testaddress.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Email != "new@mail.local" || body.OwnerUserID != 1 {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestTestAddressHandlerCreateDuplicateEmail(t *testing.T) {
	e := echo.New()
	service := testaddress.NewService(&testAddressRepo{used: map[string]testaddress.TestMailAddress{"dup@mail.local": {ID: 1, Email: "dup@mail.local"}}}, testAddressUserReader{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}, testAddressTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}})
	h := NewTestAddressHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/test-addresses", strings.NewReader(`{"email":"dup@mail.local"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid-access-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Create(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rec.Code)
	}
}
