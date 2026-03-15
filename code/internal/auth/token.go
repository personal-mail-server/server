package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type AuthTokenClaims struct {
	TokenUse       string `json:"typ"`
	SessionVersion int    `json:"ver"`
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJWTIssuer(accessSecret, refreshSecret string) *JWTIssuer {
	return &JWTIssuer{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (j *JWTIssuer) IssuePair(now time.Time, loginID string, sessionVersion int, refreshTokenID string) (*IssuedTokenPair, error) {
	access, err := issueToken(j.accessSecret, now, now.Add(30*time.Minute), loginID, TokenUseAccess, sessionVersion)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	refreshExpiresAt := now.Add(7 * 24 * time.Hour)
	refresh, err := issueRefreshToken(j.refreshSecret, now, refreshExpiresAt, loginID, sessionVersion, refreshTokenID)
	if err != nil {
		return nil, fmt.Errorf("issue refresh token: %w", err)
	}

	return &IssuedTokenPair{
		AccessToken:           access,
		RefreshToken:          refresh,
		RefreshTokenID:        refreshTokenID,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func (j *JWTIssuer) VerifyAccessToken(rawToken string) (*AuthTokenClaims, error) {
	return verifyToken(rawToken, j.accessSecret, TokenUseAccess)
}

func (j *JWTIssuer) VerifyRefreshToken(rawToken string) (*AuthTokenClaims, error) {
	claims, err := verifyToken(rawToken, j.refreshSecret, TokenUseRefresh)
	if err != nil {
		return nil, err
	}
	if claims.ID == "" {
		return nil, fmt.Errorf("%w: malformed refresh token claims", ErrInvalidToken)
	}
	return claims, nil
}

func verifyToken(rawToken string, secret []byte, expectedUse string) (*AuthTokenClaims, error) {
	claims := &AuthTokenClaims{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("%w: unexpected signing method", ErrInvalidToken)
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: parse token: %v", ErrInvalidToken, err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("%w: token not valid", ErrInvalidToken)
	}
	if claims.TokenUse != expectedUse || claims.Subject == "" || claims.SessionVersion <= 0 {
		return nil, fmt.Errorf("%w: malformed token claims", ErrInvalidToken)
	}
	return claims, nil
}

func NewRefreshTokenID() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate refresh token id: %w", err)
	}
	return hex.EncodeToString(raw), nil
}

func issueToken(secret []byte, now, expiresAt time.Time, loginID, tokenUse string, sessionVersion int) (string, error) {
	claims := AuthTokenClaims{
		TokenUse:       tokenUse,
		SessionVersion: sessionVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   loginID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func issueRefreshToken(secret []byte, now, expiresAt time.Time, loginID string, sessionVersion int, tokenID string) (string, error) {
	claims := AuthTokenClaims{
		TokenUse:       TokenUseRefresh,
		SessionVersion: sessionVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   loginID,
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
