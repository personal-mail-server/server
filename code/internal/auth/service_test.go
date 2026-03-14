package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time { return f.now }

type fakeIssuer struct{}

func (f fakeIssuer) IssuePair(_ time.Time, _ string) (string, string, error) {
	return "access-token", "refresh-token", nil
}

type memoryRepo struct {
	users map[string]*User
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{users: map[string]*User{}}
}

func (m *memoryRepo) FindByLoginID(_ context.Context, loginID string) (*User, error) {
	user, ok := m.users[loginID]
	if !ok {
		return nil, ErrUserNotFound
	}
	copy := *user
	return &copy, nil
}

func (m *memoryRepo) IncrementFailure(_ context.Context, userID int64, now time.Time) (int, *time.Time, error) {
	for _, user := range m.users {
		if user.ID == userID {
			user.FailedAttempts++
			if user.FailedAttempts >= MaxFailedAttempts {
				lockedUntil := now.Add(LockDuration)
				user.LockedUntil = &lockedUntil
			}
			return user.FailedAttempts, user.LockedUntil, nil
		}
	}
	return 0, nil, errors.New("user not found")
}

func (m *memoryRepo) ResetFailures(_ context.Context, userID int64) error {
	for _, user := range m.users {
		if user.ID == userID {
			user.FailedAttempts = 0
			user.LockedUntil = nil
			return nil
		}
	}
	return errors.New("user not found")
}

func TestLoginSuccessResetsFailures(t *testing.T) {
	hash, err := HashPassword("pass1234")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := newMemoryRepo()
	repo.users["user-1"] = &User{ID: 1, LoginID: "user-1", PasswordHash: hash, FailedAttempts: 4}
	service := NewService(repo, fakeIssuer{}, fakeClock{now: time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)})

	resp, appErr := service.Login(context.Background(), LoginRequest{LoginID: "user-1", Password: "pass1234"})
	if appErr != nil {
		t.Fatalf("expected success, got error: %+v", appErr)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatalf("expected both tokens")
	}
	if resp.AccessTokenExpiresIn != 1800 || resp.RefreshTokenExpiresIn != 604800 || resp.TokenType != "Bearer" {
		t.Fatalf("unexpected token metadata: %+v", resp)
	}

	stored := repo.users["user-1"]
	if stored.FailedAttempts != 0 {
		t.Fatalf("failed attempts should reset to 0, got %d", stored.FailedAttempts)
	}
	if stored.LockedUntil != nil {
		t.Fatalf("locked until should reset on success")
	}
}

func TestLoginInvalidCredentialFlowAndLockOnFifthFailure(t *testing.T) {
	hash, err := HashPassword("pass1234")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	now := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	repo := newMemoryRepo()
	repo.users["user-2"] = &User{ID: 2, LoginID: "user-2", PasswordHash: hash, FailedAttempts: 3}
	service := NewService(repo, fakeIssuer{}, fakeClock{now: now})

	_, err4 := service.Login(context.Background(), LoginRequest{LoginID: "user-2", Password: "wrong1234"})
	if err4 == nil || err4.Status != 401 || err4.Code != CodeInvalidCredentials {
		t.Fatalf("4th failure should be 401 INVALID_CREDENTIALS, got %+v", err4)
	}

	_, err5 := service.Login(context.Background(), LoginRequest{LoginID: "user-2", Password: "wrong1234"})
	if err5 == nil || err5.Status != 423 || err5.Code != CodeAccountLocked {
		t.Fatalf("5th failure should be 423 ACCOUNT_LOCKED, got %+v", err5)
	}

	stored := repo.users["user-2"]
	if stored.LockedUntil == nil {
		t.Fatalf("expected locked_until to be set on 5th failure")
	}
	if stored.LockedUntil.Sub(now) != 10*time.Minute {
		t.Fatalf("expected 10 minute lock, got %v", stored.LockedUntil.Sub(now))
	}
}

func TestLoginBlockedDuringLockAndAllowedAfterUnlock(t *testing.T) {
	hash, err := HashPassword("pass1234")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	now := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	lockedUntil := now.Add(2 * time.Minute)
	repo := newMemoryRepo()
	repo.users["user-3"] = &User{ID: 3, LoginID: "user-3", PasswordHash: hash, FailedAttempts: 5, LockedUntil: &lockedUntil}

	serviceLocked := NewService(repo, fakeIssuer{}, fakeClock{now: now})
	_, appErr := serviceLocked.Login(context.Background(), LoginRequest{LoginID: "user-3", Password: "pass1234"})
	if appErr == nil || appErr.Status != 423 {
		t.Fatalf("expected locked response while lock active, got %+v", appErr)
	}

	afterUnlock := lockedUntil.Add(time.Second)
	serviceAfter := NewService(repo, fakeIssuer{}, fakeClock{now: afterUnlock})
	_, appErr = serviceAfter.Login(context.Background(), LoginRequest{LoginID: "user-3", Password: "pass1234"})
	if appErr != nil {
		t.Fatalf("expected success after lock expired, got %+v", appErr)
	}
}
