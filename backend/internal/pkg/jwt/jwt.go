package jwtutil

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(secret string, accessTTL time.Duration) *Manager {
	return &Manager{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: 30 * 24 * time.Hour, // 30 days
	}
}

// Generate creates a short-lived access token (JWT).
func (m *Manager) Generate(userID int64, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// Validate parses and validates an access token.
func (m *Manager) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// GenerateRefreshToken creates a cryptographically random opaque token (32 bytes).
// Returns the raw token (stored in cookie) and its SHA-256 hash (stored in DB).
func (m *Manager) GenerateRefreshToken() (raw string, hash string, expiresAt time.Time, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", time.Time{}, fmt.Errorf("generate random bytes: %w", err)
	}
	raw      = hex.EncodeToString(b)
	sum      := sha256.Sum256([]byte(raw))
	hash     = hex.EncodeToString(sum[:])
	expiresAt = time.Now().Add(m.refreshTokenTTL)
	return raw, hash, expiresAt, nil
}

// RefreshTokenTTL returns the duration used for refresh tokens.
func (m *Manager) RefreshTokenTTL() time.Duration {
	return m.refreshTokenTTL
}

// AccessTokenTTL returns the duration used for access tokens.
func (m *Manager) AccessTokenTTL() time.Duration {
	return m.accessTokenTTL
}
