package auth

import (
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

func (j *JWTIssuer) IssuePair(now time.Time, loginID string, sessionVersion int) (string, string, error) {
	access, err := issueToken(j.accessSecret, now, now.Add(30*time.Minute), loginID, TokenUseAccess, sessionVersion)
	if err != nil {
		return "", "", fmt.Errorf("issue access token: %w", err)
	}

	refresh, err := issueToken(j.refreshSecret, now, now.Add(7*24*time.Hour), loginID, TokenUseRefresh, sessionVersion)
	if err != nil {
		return "", "", fmt.Errorf("issue refresh token: %w", err)
	}

	return access, refresh, nil
}

func (j *JWTIssuer) VerifyAccessToken(rawToken string) (*AuthTokenClaims, error) {
	claims := &AuthTokenClaims{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("%w: unexpected signing method", ErrInvalidToken)
		}
		return j.accessSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: parse access token: %v", ErrInvalidToken, err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("%w: token not valid", ErrInvalidToken)
	}
	if claims.TokenUse != TokenUseAccess || claims.Subject == "" || claims.SessionVersion <= 0 {
		return nil, fmt.Errorf("%w: malformed access token claims", ErrInvalidToken)
	}
	return claims, nil
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
