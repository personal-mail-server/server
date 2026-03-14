package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type fixedIssuer struct{}

func (fixedIssuer) IssuePair(_ time.Time, _ string) (string, string, error) {
	return "a", "r", nil
}

func TestAuthHandlerInvalidJSON(t *testing.T) {
	e := echo.New()
	repo := &testRepo{}
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
