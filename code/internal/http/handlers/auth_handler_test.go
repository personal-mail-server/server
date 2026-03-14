package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"personal-mail-server/internal/auth"
)

type testRepo struct {
	user *auth.User
}

func (r *testRepo) FindByLoginID(_ context.Context, loginID string) (*auth.User, error) {
	if r.user == nil || r.user.LoginID != loginID {
		return nil, auth.ErrUserNotFound
	}
	copy := *r.user
	return &copy, nil
}

func (r *testRepo) IncrementFailure(_ context.Context, _ int64, now time.Time) (int, *time.Time, error) {
	r.user.FailedAttempts++
	if r.user.FailedAttempts >= auth.MaxFailedAttempts {
		lockedUntil := now.Add(auth.LockDuration)
		r.user.LockedUntil = &lockedUntil
	}
	return r.user.FailedAttempts, r.user.LockedUntil, nil
}

func (r *testRepo) ResetFailures(_ context.Context, _ int64) error {
	r.user.FailedAttempts = 0
	r.user.LockedUntil = nil
	return nil
}

func (r *testRepo) IncrementSessionVersion(_ context.Context, _ int64, currentVersion int) (bool, error) {
	if r.user == nil {
		return false, auth.ErrUserNotFound
	}
	if r.user.SessionVersion != currentVersion {
		return false, nil
	}
	r.user.SessionVersion++
	return true, nil
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type fixedIssuer struct{}

func (fixedIssuer) IssuePair(_ time.Time, _ string, _ int) (string, string, error) {
	return "a", "r", nil
}

func (fixedIssuer) VerifyAccessToken(raw string) (*auth.AuthTokenClaims, error) {
	if raw != "valid-access-token" {
		return nil, auth.ErrInvalidToken
	}
	return &auth.AuthTokenClaims{
		TokenUse:       auth.TokenUseAccess,
		SessionVersion: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "user-01",
		},
	}, nil
}

func TestAuthHandlerInvalidJSON(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}
	service := auth.NewService(repo, fixedIssuer{}, fixedClock{now: time.Now().UTC()})
	h := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Login(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	var body auth.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code != auth.CodeInvalidRequestBody {
		t.Fatalf("expected INVALID_REQUEST_BODY, got %s", body.Code)
	}
}

func TestAuthHandlerLogoutMissingAuthorization(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}
	service := auth.NewService(repo, fixedIssuer{}, fixedClock{now: time.Now().UTC()})
	h := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Logout(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	var body auth.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code != auth.CodeInvalidAccessToken {
		t.Fatalf("expected INVALID_ACCESS_TOKEN, got %s", body.Code)
	}
}

func TestAuthHandlerLogoutSuccess(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}
	service := auth.NewService(repo, fixedIssuer{}, fixedClock{now: time.Now().UTC()})
	h := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid-access-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Logout(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", rec.Body.String())
	}
	if repo.user.SessionVersion != 2 {
		t.Fatalf("expected session version incremented to 2, got %d", repo.user.SessionVersion)
	}
}
