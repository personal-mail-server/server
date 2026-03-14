package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func (j *JWTIssuer) IssuePair(now time.Time, loginID string) (string, string, error) {
	access, err := issueToken(j.accessSecret, now, now.Add(30*time.Minute), loginID, "access")
	if err != nil {
		return "", "", fmt.Errorf("issue access token: %w", err)
	}

	refresh, err := issueToken(j.refreshSecret, now, now.Add(7*24*time.Hour), loginID, "refresh")
	if err != nil {
		return "", "", fmt.Errorf("issue refresh token: %w", err)
	}

	return access, refresh, nil
}

func issueToken(secret []byte, now, expiresAt time.Time, loginID, tokenUse string) (string, error) {
	claims := jwt.MapClaims{
		"sub": loginID,
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": expiresAt.Unix(),
		"typ": tokenUse,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
