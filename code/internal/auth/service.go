package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

type Service struct {
	repo   Repository
	issuer TokenIssuer
	clock  Clock
}

func NewService(repo Repository, issuer TokenIssuer, clock Clock) *Service {
	if clock == nil {
		clock = RealClock{}
	}
	return &Service{repo: repo, issuer: issuer, clock: clock}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, *AppError) {
	if validationErr := ValidateLoginRequest(req); validationErr != nil {
		return nil, validationErr
	}

	now := s.clock.Now()
	user, err := s.repo.FindByLoginID(ctx, req.LoginID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, NewUnauthorized()
		}
		return nil, NewInternalServerError()
	}

	if user.LockedUntil != nil && now.Before(*user.LockedUntil) {
		return nil, NewLocked()
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		failedAttempts, lockedUntil, incrementErr := s.repo.IncrementFailure(ctx, user.ID, now)
		if incrementErr != nil {
			return nil, NewInternalServerError()
		}

		if failedAttempts >= MaxFailedAttempts {
			if lockedUntil != nil && now.Before(*lockedUntil) {
				return nil, NewLocked()
			}
			return nil, NewLocked()
		}
		return nil, NewUnauthorized()
	}

	if err := s.repo.ResetFailures(ctx, user.ID); err != nil {
		return nil, NewInternalServerError()
	}

	refreshTokenID, err := NewRefreshTokenID()
	if err != nil {
		return nil, NewInternalServerError()
	}

	pair, err := s.issuer.IssuePair(now, user.LoginID, user.SessionVersion, refreshTokenID)
	if err != nil {
		return nil, NewInternalServerError()
	}
	if err := s.repo.StoreRefreshToken(ctx, user.ID, pair.RefreshTokenID, user.SessionVersion, pair.RefreshTokenExpiresAt); err != nil {
		return nil, NewInternalServerError()
	}

	return &LoginResponse{
		AccessToken:           pair.AccessToken,
		RefreshToken:          pair.RefreshToken,
		AccessTokenExpiresIn:  AccessTokenExpiresInSeconds,
		RefreshTokenExpiresIn: RefreshTokenExpiresInSeconds,
		TokenType:             TokenTypeBearer,
	}, nil
}

func (s *Service) Reissue(ctx context.Context, req ReissueRequest) (*LoginResponse, *AppError) {
	if validationErr := ValidateReissueRequest(req); validationErr != nil {
		return nil, validationErr
	}

	claims, err := s.issuer.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, NewInvalidRefreshToken()
	}

	user, repoErr := s.repo.FindByLoginID(ctx, claims.Subject)
	if repoErr != nil {
		if errors.Is(repoErr, ErrUserNotFound) {
			return nil, NewInvalidRefreshToken()
		}
		return nil, NewInternalServerError()
	}
	if user.SessionVersion != claims.SessionVersion {
		return nil, NewInvalidRefreshToken()
	}

	now := s.clock.Now()
	replacementTokenID, err := NewRefreshTokenID()
	if err != nil {
		return nil, NewInternalServerError()
	}

	pair, err := s.issuer.IssuePair(now, user.LoginID, user.SessionVersion, replacementTokenID)
	if err != nil {
		return nil, NewInternalServerError()
	}

	rotated, err := s.repo.ConsumeRefreshTokenAndStoreReplacement(ctx, user.ID, claims.ID, pair.RefreshTokenID, user.SessionVersion, now, pair.RefreshTokenExpiresAt)
	if err != nil {
		return nil, NewInternalServerError()
	}
	if !rotated {
		return nil, NewInvalidRefreshToken()
	}

	return &LoginResponse{
		AccessToken:           pair.AccessToken,
		RefreshToken:          pair.RefreshToken,
		AccessTokenExpiresIn:  AccessTokenExpiresInSeconds,
		RefreshTokenExpiresIn: RefreshTokenExpiresInSeconds,
		TokenType:             TokenTypeBearer,
	}, nil
}

func (s *Service) Logout(ctx context.Context, rawAccessToken string) *AppError {
	claims, err := s.issuer.VerifyAccessToken(rawAccessToken)
	if err != nil {
		return NewInvalidAccessToken()
	}

	user, repoErr := s.repo.FindByLoginID(ctx, claims.Subject)
	if repoErr != nil {
		if errors.Is(repoErr, ErrUserNotFound) {
			return NewInvalidAccessToken()
		}
		return NewInternalServerError()
	}

	if user.SessionVersion != claims.SessionVersion {
		return NewInvalidAccessToken()
	}

	updated, err := s.repo.IncrementSessionVersion(ctx, user.ID, user.SessionVersion)
	if err != nil {
		return NewInternalServerError()
	}
	if !updated {
		return NewInvalidAccessToken()
	}

	return nil
}

func HashPassword(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password hash: %w", err)
	}
	return string(hash), nil
}
