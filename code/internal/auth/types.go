package auth

import (
	"context"
	"time"
)

const (
	AccessTokenExpiresInSeconds  = 1800
	RefreshTokenExpiresInSeconds = 604800
	TokenTypeBearer              = "Bearer"
	TokenUseAccess               = "access"
	TokenUseRefresh              = "refresh"

	MaxFailedAttempts = 5
	LockDuration      = 10 * time.Minute
)

type LoginRequest struct {
	LoginID  string `json:"loginId"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExpiresIn  int    `json:"accessTokenExpiresIn"`
	RefreshTokenExpiresIn int    `json:"refreshTokenExpiresIn"`
	TokenType             string `json:"tokenType"`
}

type ReissueRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type User struct {
	ID             int64
	LoginID        string
	PasswordHash   string
	FailedAttempts int
	LockedUntil    *time.Time
	SessionVersion int
}

type IssuedTokenPair struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenID        string
	RefreshTokenExpiresAt time.Time
}

type Repository interface {
	FindByLoginID(ctx context.Context, loginID string) (*User, error)
	IncrementFailure(ctx context.Context, userID int64, now time.Time) (int, *time.Time, error)
	ResetFailures(ctx context.Context, userID int64) error
	IncrementSessionVersion(ctx context.Context, userID int64, currentVersion int) (bool, error)
	StoreRefreshToken(ctx context.Context, userID int64, tokenID string, sessionVersion int, expiresAt time.Time) error
	ConsumeRefreshTokenAndStoreReplacement(ctx context.Context, userID int64, currentTokenID, replacementTokenID string, sessionVersion int, now, replacementExpiresAt time.Time) (bool, error)
}

type TokenIssuer interface {
	IssuePair(now time.Time, loginID string, sessionVersion int, refreshTokenID string) (*IssuedTokenPair, error)
	VerifyAccessToken(rawToken string) (*AuthTokenClaims, error)
	VerifyRefreshToken(rawToken string) (*AuthTokenClaims, error)
}
