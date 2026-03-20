package testaddress

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"personal-mail-server/internal/auth"
)

type fakeUserReader struct {
	user *auth.User
	err  error
}

func (f fakeUserReader) FindByLoginID(_ context.Context, loginID string) (*auth.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.user == nil || f.user.LoginID != loginID {
		return nil, auth.ErrUserNotFound
	}
	copy := *f.user
	return &copy, nil
}

type fakeTokenIssuer struct {
	claims *auth.AuthTokenClaims
	err    error
}

func (f fakeTokenIssuer) IssuePair(_ time.Time, _ string, _ int, _ string) (*auth.IssuedTokenPair, error) {
	panic("not used")
}

func (f fakeTokenIssuer) VerifyAccessToken(_ string) (*auth.AuthTokenClaims, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.claims == nil {
		return nil, auth.ErrInvalidToken
	}
	return f.claims, nil
}

func (f fakeTokenIssuer) VerifyRefreshToken(_ string) (*auth.AuthTokenClaims, error) {
	panic("not used")
}

type fakeCandidateGenerator struct {
	values []string
	err    error
	idx    int
}

func (f *fakeCandidateGenerator) Next() (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if f.idx >= len(f.values) {
		return "", errors.New("no candidate")
	}
	value := f.values[f.idx]
	f.idx++
	return value, nil
}

type memoryRepo struct {
	byEmail map[string]TestMailAddress
	err     error
}

func (m *memoryRepo) Create(context.Context, TestMailAddress) (*TestMailAddress, error) {
	panic("not used")
}
func (m *memoryRepo) GetByID(context.Context, int64) (*TestMailAddress, error) { panic("not used") }
func (m *memoryRepo) ListByOwner(context.Context, int64) ([]TestMailAddress, error) {
	panic("not used")
}
func (m *memoryRepo) Update(context.Context, TestMailAddress) (*TestMailAddress, error) {
	panic("not used")
}
func (m *memoryRepo) SoftDelete(context.Context, int64, time.Time) error { panic("not used") }

func (m *memoryRepo) GetByEmail(_ context.Context, email string) (*TestMailAddress, error) {
	if m.err != nil {
		return nil, m.err
	}
	address, ok := m.byEmail[email]
	if !ok {
		return nil, ErrTestMailAddressNotFound
	}
	copy := address
	return &copy, nil
}

func TestGenerateCandidateReturnsAvailableEmail(t *testing.T) {
	repo := &memoryRepo{byEmail: map[string]TestMailAddress{"taken@mail.local": {ID: 1, Email: "taken@mail.local"}}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 2, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 10, LoginID: "user-01", SessionVersion: 2}}
	gen := &fakeCandidateGenerator{values: []string{"taken@mail.local", "free@mail.local"}}
	service := newService(repo, users, issuer, gen)

	resp, appErr := service.GenerateCandidate(context.Background(), "valid-access-token")
	if appErr != nil {
		t.Fatalf("expected success, got %+v", appErr)
	}
	if resp.Email != "free@mail.local" {
		t.Fatalf("expected free candidate, got %+v", resp)
	}
}

func TestGenerateCandidateRejectsInvalidAccessToken(t *testing.T) {
	service := newService(&memoryRepo{byEmail: map[string]TestMailAddress{}}, fakeUserReader{}, fakeTokenIssuer{err: auth.ErrInvalidToken}, &fakeCandidateGenerator{})

	_, appErr := service.GenerateCandidate(context.Background(), "bad-token")
	if appErr == nil || appErr.Code != auth.CodeInvalidAccessToken {
		t.Fatalf("expected invalid access token, got %+v", appErr)
	}
}

func TestGenerateCandidateRejectsStaleSessionVersion(t *testing.T) {
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 10, LoginID: "user-01", SessionVersion: 2}}
	service := newService(&memoryRepo{byEmail: map[string]TestMailAddress{}}, users, issuer, &fakeCandidateGenerator{})

	_, appErr := service.GenerateCandidate(context.Background(), "stale-token")
	if appErr == nil || appErr.Code != auth.CodeInvalidAccessToken {
		t.Fatalf("expected invalid access token, got %+v", appErr)
	}
}

func TestGenerateCandidateReturnsInternalErrorOnRepositoryFailure(t *testing.T) {
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 10, LoginID: "user-01", SessionVersion: 1}}
	service := newService(&memoryRepo{err: errors.New("db down")}, users, issuer, &fakeCandidateGenerator{values: []string{"free@mail.local"}})

	_, appErr := service.GenerateCandidate(context.Background(), "valid-token")
	if appErr == nil || appErr.Code != auth.CodeInternalServerError {
		t.Fatalf("expected internal server error, got %+v", appErr)
	}
}
