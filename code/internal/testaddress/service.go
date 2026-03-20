package testaddress

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"personal-mail-server/internal/auth"
)

const (
	generatedEmailDomain      = "mail.local"
	maxGenerateCandidateTries = 16
)

type UserReader interface {
	FindByLoginID(ctx context.Context, loginID string) (*auth.User, error)
}

type CandidateGenerator interface {
	Next() (string, error)
}

type GenerateCandidateResponse struct {
	Email string `json:"email"`
}

type randomCandidateGenerator struct{}

func (randomCandidateGenerator) Next() (string, error) {
	raw := make([]byte, 6)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate random candidate bytes: %w", err)
	}
	return "test-" + hex.EncodeToString(raw) + "@" + generatedEmailDomain, nil
}

type Service struct {
	repo      Repository
	users     UserReader
	issuer    auth.TokenIssuer
	generator CandidateGenerator
}

func NewService(repo Repository, users UserReader, issuer auth.TokenIssuer) *Service {
	return &Service{
		repo:      repo,
		users:     users,
		issuer:    issuer,
		generator: randomCandidateGenerator{},
	}
}

func newService(repo Repository, users UserReader, issuer auth.TokenIssuer, generator CandidateGenerator) *Service {
	if generator == nil {
		generator = randomCandidateGenerator{}
	}
	return &Service{repo: repo, users: users, issuer: issuer, generator: generator}
}

func (s *Service) GenerateCandidate(ctx context.Context, rawAccessToken string) (*GenerateCandidateResponse, *auth.AppError) {
	claims, err := s.issuer.VerifyAccessToken(rawAccessToken)
	if err != nil {
		return nil, auth.NewInvalidAccessToken()
	}

	user, repoErr := s.users.FindByLoginID(ctx, claims.Subject)
	if repoErr != nil {
		if errors.Is(repoErr, auth.ErrUserNotFound) {
			return nil, auth.NewInvalidAccessToken()
		}
		return nil, auth.NewInternalServerError()
	}
	if user.SessionVersion != claims.SessionVersion {
		return nil, auth.NewInvalidAccessToken()
	}

	for range maxGenerateCandidateTries {
		email, genErr := s.generator.Next()
		if genErr != nil {
			return nil, auth.NewInternalServerError()
		}

		_, findErr := s.repo.GetByEmail(ctx, email)
		if findErr == nil {
			continue
		}
		if errors.Is(findErr, ErrTestMailAddressNotFound) {
			return &GenerateCandidateResponse{Email: email}, nil
		}
		return nil, auth.NewInternalServerError()
	}

	return nil, auth.NewInternalServerError()
}
