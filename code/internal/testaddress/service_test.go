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
	byID    map[int64]TestMailAddress
	err     error
}

func (m *memoryRepo) Create(_ context.Context, address TestMailAddress) (*TestMailAddress, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, ok := m.byEmail[address.Email]; ok {
		return nil, ErrDuplicateEmail
	}
	created := TestMailAddress{
		ID:          int64(len(m.byEmail) + 1),
		OwnerUserID: address.OwnerUserID,
		Email:       address.Email,
		CreatedAt:   time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
	}
	if m.byEmail == nil {
		m.byEmail = map[string]TestMailAddress{}
	}
	if m.byID == nil {
		m.byID = map[int64]TestMailAddress{}
	}
	m.byEmail[address.Email] = created
	m.byID[created.ID] = created
	return &created, nil
}

func (m *memoryRepo) GetByID(_ context.Context, id int64) (*TestMailAddress, error) {
	if m.err != nil {
		return nil, m.err
	}
	address, ok := m.byID[id]
	if !ok {
		return nil, ErrTestMailAddressNotFound
	}
	copy := address
	return &copy, nil
}

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

func (m *memoryRepo) ListByOwner(_ context.Context, ownerID int64) ([]TestMailAddress, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]TestMailAddress, 0)
	for _, address := range m.byEmail {
		if address.OwnerUserID == ownerID && address.DeletedAt == nil {
			result = append(result, address)
		}
	}
	return result, nil
}

func (m *memoryRepo) Update(_ context.Context, address TestMailAddress) (*TestMailAddress, error) {
	if m.err != nil {
		return nil, m.err
	}
	current, ok := m.byID[address.ID]
	if !ok || current.DeletedAt != nil {
		return nil, ErrTestMailAddressNotFound
	}
	if existing, exists := m.byEmail[address.Email]; exists && existing.ID != address.ID {
		return nil, ErrDuplicateEmail
	}
	delete(m.byEmail, current.Email)
	current.Email = address.Email
	current.UpdatedAt = time.Date(2026, 3, 20, 2, 0, 0, 0, time.UTC)
	m.byEmail[current.Email] = current
	m.byID[current.ID] = current
	copy := current
	return &copy, nil
}

func (m *memoryRepo) SoftDelete(context.Context, int64, time.Time) error { panic("not used") }

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

func TestCreateStoresAddressForAuthenticatedUser(t *testing.T) {
	repo := &memoryRepo{byEmail: map[string]TestMailAddress{}, byID: map[int64]TestMailAddress{}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 3, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 44, LoginID: "user-01", SessionVersion: 3}}
	service := newService(repo, users, issuer, nil)

	resp, appErr := service.Create(context.Background(), "valid-token", CreateRequest{Email: "created@mail.local"})
	if appErr != nil {
		t.Fatalf("expected success, got %+v", appErr)
	}
	if resp.Email != "created@mail.local" || resp.OwnerUserID != 44 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	stored, ok := repo.byEmail["created@mail.local"]
	if !ok || stored.OwnerUserID != 44 {
		t.Fatalf("expected owner-linked stored address, got %+v", stored)
	}
}

func TestCreateRejectsDuplicateEmail(t *testing.T) {
	repo := &memoryRepo{byEmail: map[string]TestMailAddress{"dup@mail.local": {ID: 1, Email: "dup@mail.local"}}, byID: map[int64]TestMailAddress{1: {ID: 1, Email: "dup@mail.local"}}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 1, LoginID: "user-01", SessionVersion: 1}}
	service := newService(repo, users, issuer, nil)

	_, appErr := service.Create(context.Background(), "valid-token", CreateRequest{Email: "dup@mail.local"})
	if appErr == nil || appErr.Status != 409 || appErr.Code != auth.CodeDuplicateEmail {
		t.Fatalf("expected duplicate conflict, got %+v", appErr)
	}
}

func TestCreateRejectsInvalidEmailFormat(t *testing.T) {
	service := newService(&memoryRepo{byEmail: map[string]TestMailAddress{}}, fakeUserReader{}, fakeTokenIssuer{}, nil)

	_, appErr := service.Create(context.Background(), "valid-token", CreateRequest{Email: "not-an-email"})
	if appErr == nil || appErr.Code != auth.CodeInvalidEmail {
		t.Fatalf("expected invalid email error, got %+v", appErr)
	}
}

func TestCreateRejectsMissingEmail(t *testing.T) {
	service := newService(&memoryRepo{byEmail: map[string]TestMailAddress{}}, fakeUserReader{}, fakeTokenIssuer{}, nil)

	_, appErr := service.Create(context.Background(), "valid-token", CreateRequest{})
	if appErr == nil || appErr.Code != auth.CodeMissingRequired {
		t.Fatalf("expected missing required error, got %+v", appErr)
	}
}

func TestListReturnsOnlyOwnersActiveAddresses(t *testing.T) {
	deletedAt := time.Date(2026, 3, 20, 1, 0, 0, 0, time.UTC)
	repo := &memoryRepo{byEmail: map[string]TestMailAddress{
		"one@mail.local":   {ID: 1, OwnerUserID: 7, Email: "one@mail.local"},
		"two@mail.local":   {ID: 2, OwnerUserID: 7, Email: "two@mail.local"},
		"other@mail.local": {ID: 3, OwnerUserID: 9, Email: "other@mail.local"},
		"gone@mail.local":  {ID: 4, OwnerUserID: 7, Email: "gone@mail.local", DeletedAt: &deletedAt},
	}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 7, LoginID: "user-01", SessionVersion: 1}}
	service := newService(repo, users, issuer, nil)

	resp, appErr := service.List(context.Background(), "valid-token")
	if appErr != nil {
		t.Fatalf("expected success, got %+v", appErr)
	}
	if len(resp.Addresses) != 2 {
		t.Fatalf("expected 2 addresses, got %+v", resp)
	}
}

func TestGetByIDReturnsAddressForOwner(t *testing.T) {
	repo := &memoryRepo{byID: map[int64]TestMailAddress{11: {ID: 11, OwnerUserID: 5, Email: "owned@mail.local"}}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 2, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 5, LoginID: "user-01", SessionVersion: 2}}
	service := newService(repo, users, issuer, nil)

	resp, appErr := service.GetByID(context.Background(), "valid-token", "11")
	if appErr != nil {
		t.Fatalf("expected success, got %+v", appErr)
	}
	if resp.ID != 11 || resp.Email != "owned@mail.local" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetByIDReturnsNotFoundForNonOwner(t *testing.T) {
	repo := &memoryRepo{byID: map[int64]TestMailAddress{11: {ID: 11, OwnerUserID: 99, Email: "hidden@mail.local"}}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 2, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 5, LoginID: "user-01", SessionVersion: 2}}
	service := newService(repo, users, issuer, nil)

	_, appErr := service.GetByID(context.Background(), "valid-token", "11")
	if appErr == nil || appErr.Status != 404 {
		t.Fatalf("expected not found, got %+v", appErr)
	}
}

func TestUpdateSucceedsWhenKeepingSameEmail(t *testing.T) {
	repo := &memoryRepo{
		byEmail: map[string]TestMailAddress{"same@mail.local": {ID: 8, OwnerUserID: 3, Email: "same@mail.local"}},
		byID:    map[int64]TestMailAddress{8: {ID: 8, OwnerUserID: 3, Email: "same@mail.local"}},
	}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 3, LoginID: "user-01", SessionVersion: 1}}
	service := newService(repo, users, issuer, nil)

	resp, appErr := service.Update(context.Background(), "valid-token", "8", UpdateRequest{Email: "same@mail.local"})
	if appErr != nil {
		t.Fatalf("expected success, got %+v", appErr)
	}
	if resp.Email != "same@mail.local" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestUpdateRejectsDuplicateEmailFromAnotherResource(t *testing.T) {
	repo := &memoryRepo{
		byEmail: map[string]TestMailAddress{
			"own@mail.local":   {ID: 8, OwnerUserID: 3, Email: "own@mail.local"},
			"other@mail.local": {ID: 9, OwnerUserID: 3, Email: "other@mail.local"},
		},
		byID: map[int64]TestMailAddress{
			8: {ID: 8, OwnerUserID: 3, Email: "own@mail.local"},
			9: {ID: 9, OwnerUserID: 3, Email: "other@mail.local"},
		},
	}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 3, LoginID: "user-01", SessionVersion: 1}}
	service := newService(repo, users, issuer, nil)

	_, appErr := service.Update(context.Background(), "valid-token", "8", UpdateRequest{Email: "other@mail.local"})
	if appErr == nil || appErr.Status != 409 || appErr.Code != auth.CodeDuplicateEmail {
		t.Fatalf("expected duplicate conflict, got %+v", appErr)
	}
}

func TestUpdateRejectsNonOwnerAsNotFound(t *testing.T) {
	repo := &memoryRepo{byID: map[int64]TestMailAddress{8: {ID: 8, OwnerUserID: 99, Email: "hidden@mail.local"}}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 3, LoginID: "user-01", SessionVersion: 1}}
	service := newService(repo, users, issuer, nil)

	_, appErr := service.Update(context.Background(), "valid-token", "8", UpdateRequest{Email: "new@mail.local"})
	if appErr == nil || appErr.Status != 404 {
		t.Fatalf("expected not found, got %+v", appErr)
	}
}

func TestUpdateRejectsDeletedResourceAsNotFound(t *testing.T) {
	deletedAt := time.Date(2026, 3, 20, 3, 0, 0, 0, time.UTC)
	repo := &memoryRepo{byID: map[int64]TestMailAddress{8: {ID: 8, OwnerUserID: 3, Email: "gone@mail.local", DeletedAt: &deletedAt}}}
	issuer := fakeTokenIssuer{claims: &auth.AuthTokenClaims{TokenUse: auth.TokenUseAccess, SessionVersion: 1, RegisteredClaims: jwt.RegisteredClaims{Subject: "user-01"}}}
	users := fakeUserReader{user: &auth.User{ID: 3, LoginID: "user-01", SessionVersion: 1}}
	service := newService(repo, users, issuer, nil)

	_, appErr := service.Update(context.Background(), "valid-token", "8", UpdateRequest{Email: "new@mail.local"})
	if appErr == nil || appErr.Status != 404 {
		t.Fatalf("expected not found, got %+v", appErr)
	}
}
