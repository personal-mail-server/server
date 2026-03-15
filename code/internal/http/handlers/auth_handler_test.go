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

type testRefreshTokenRecord struct {
	userID         int64
	sessionVersion int
	expiresAt      time.Time
	usedAt         *time.Time
	replacedBy     string
}

type testRepo struct {
	user          *auth.User
	refreshTokens map[string]testRefreshTokenRecord
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

func (r *testRepo) StoreRefreshToken(_ context.Context, userID int64, tokenID string, sessionVersion int, expiresAt time.Time) error {
	if r.refreshTokens == nil {
		r.refreshTokens = map[string]testRefreshTokenRecord{}
	}
	r.refreshTokens[tokenID] = testRefreshTokenRecord{userID: userID, sessionVersion: sessionVersion, expiresAt: expiresAt}
	return nil
}

func (r *testRepo) ConsumeRefreshTokenAndStoreReplacement(_ context.Context, userID int64, currentTokenID, replacementTokenID string, sessionVersion int, now, replacementExpiresAt time.Time) (bool, error) {
	if r.refreshTokens == nil {
		return false, nil
	}
	record, ok := r.refreshTokens[currentTokenID]
	if !ok || record.userID != userID || record.sessionVersion != sessionVersion || record.usedAt != nil || !record.expiresAt.After(now) {
		return false, nil
	}
	record.usedAt = &now
	record.replacedBy = replacementTokenID
	r.refreshTokens[currentTokenID] = record
	r.refreshTokens[replacementTokenID] = testRefreshTokenRecord{userID: userID, sessionVersion: sessionVersion, expiresAt: replacementExpiresAt}
	return true, nil
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type fixedIssuer struct{}

func (fixedIssuer) IssuePair(_ time.Time, _ string, _ int, refreshTokenID string) (*auth.IssuedTokenPair, error) {
	return &auth.IssuedTokenPair{
		AccessToken:           "a",
		RefreshToken:          "r",
		RefreshTokenID:        refreshTokenID,
		RefreshTokenExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour),
	}, nil
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

func (fixedIssuer) VerifyRefreshToken(raw string) (*auth.AuthTokenClaims, error) {
	if raw != "valid-refresh-token" {
		return nil, auth.ErrInvalidToken
	}
	return &auth.AuthTokenClaims{
		TokenUse:       auth.TokenUseRefresh,
		SessionVersion: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "user-01",
			ID:      "refresh-token-id",
		},
	}, nil
}

func TestAuthHandlerInvalidJSON(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}, refreshTokens: map[string]testRefreshTokenRecord{}}
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
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}, refreshTokens: map[string]testRefreshTokenRecord{}}
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
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}, refreshTokens: map[string]testRefreshTokenRecord{}}
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

func TestAuthHandlerReissueInvalidJSON(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}, refreshTokens: map[string]testRefreshTokenRecord{}}
	service := auth.NewService(repo, fixedIssuer{}, fixedClock{now: time.Now().UTC()})
	h := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token/reissue", bytes.NewBufferString("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Reissue(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestAuthHandlerReissueMissingRefreshToken(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}, refreshTokens: map[string]testRefreshTokenRecord{}}
	service := auth.NewService(repo, fixedIssuer{}, fixedClock{now: time.Now().UTC()})
	h := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token/reissue", bytes.NewBufferString(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Reissue(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestAuthHandlerReissueSuccess(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}, refreshTokens: map[string]testRefreshTokenRecord{"refresh-token-id": {userID: 1, sessionVersion: 1, expiresAt: time.Now().UTC().Add(time.Hour)}}}
	service := auth.NewService(repo, fixedIssuer{}, fixedClock{now: time.Now().UTC()})
	h := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token/reissue", bytes.NewBufferString(`{"refreshToken":"valid-refresh-token"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Reissue(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body auth.LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.AccessToken == "" || body.RefreshToken == "" {
		t.Fatalf("expected token pair response, got %+v", body)
	}
}

func TestAuthHandlerReissueInvalidRefreshToken(t *testing.T) {
	e := echo.New()
	repo := &testRepo{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}, refreshTokens: map[string]testRefreshTokenRecord{}}
	service := auth.NewService(repo, fixedIssuer{}, fixedClock{now: time.Now().UTC()})
	h := NewAuthHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token/reissue", bytes.NewBufferString(`{"refreshToken":"bad-refresh-token"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Reissue(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}
