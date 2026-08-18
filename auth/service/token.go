package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenClaims represents JWT token claims
type TokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// TokenManager handles JWT token generation and validation
type TokenManager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) *TokenManager {
	return &TokenManager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// GenerateAccessToken generates a new access token
func (tm *TokenManager) GenerateAccessToken(email string) (string, error) {
	claims := TokenClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.secretKey))
}

// GenerateRefreshToken generates a new refresh token
func (tm *TokenManager) GenerateRefreshToken(email string) (string, error) {
	claims := TokenClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.refreshTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.secretKey))
}

// ValidateAccessToken validates an access token
func (tm *TokenManager) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return tm.validateToken(tokenString)
}

// ValidateRefreshToken validates a refresh token
func (tm *TokenManager) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return tm.validateToken(tokenString)
}

// validateToken is a helper function to validate tokens
func (tm *TokenManager) validateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
