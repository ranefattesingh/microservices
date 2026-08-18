package service

import (
	"errors"
	"testing"
	"time"
)

// TestGenerateAccessToken tests access token generation
func TestGenerateAccessToken(t *testing.T) {
	tm := NewTokenManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	token, err := tm.GenerateAccessToken("test@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected token, got empty string")
	}

	// Verify token is a valid JWT string (has 3 parts separated by dots)
	dotCount := 0
	for _, ch := range token {
		if ch == '.' {
			dotCount++
		}
	}
	if dotCount != 2 {
		t.Fatalf("expected JWT with 2 dots (3 parts), got %d dots", dotCount)
	}
}

// TestGenerateRefreshToken tests refresh token generation
func TestGenerateRefreshToken(t *testing.T) {
	tm := NewTokenManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	token, err := tm.GenerateRefreshToken("test@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected token, got empty string")
	}
}

// TestValidateAccessToken tests access token validation
func TestValidateAccessToken(t *testing.T) {
	tm := NewTokenManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	email := "test@example.com"
	token, err := tm.GenerateAccessToken(email)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := tm.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	if claims.Email != email {
		t.Fatalf("expected email %s, got %s", email, claims.Email)
	}

	if claims.ExpiresAt == nil {
		t.Fatalf("expected expiration time in token")
	}
}

// TestValidateRefreshToken tests refresh token validation
func TestValidateRefreshToken(t *testing.T) {
	tm := NewTokenManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	email := "test@example.com"
	token, err := tm.GenerateRefreshToken(email)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := tm.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	if claims.Email != email {
		t.Fatalf("expected email %s, got %s", email, claims.Email)
	}
}

// TestValidateInvalidToken tests validation of invalid token
func TestValidateInvalidToken(t *testing.T) {
	tm := NewTokenManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	claims, err := tm.ValidateAccessToken("invalid.token.here")
	if err == nil {
		t.Fatalf("expected error for invalid token, got nil")
	}

	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	if claims != nil {
		t.Fatalf("expected nil claims for invalid token")
	}
}

// TestValidateTokenWithWrongSecret tests validation with wrong secret
func TestValidateTokenWithWrongSecret(t *testing.T) {
	tm1 := NewTokenManager("secret1", 15*time.Minute, 7*24*time.Hour)
	tm2 := NewTokenManager("secret2", 15*time.Minute, 7*24*time.Hour)

	token, err := tm1.GenerateAccessToken("test@example.com")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Try to validate with different secret
	claims, err := tm2.ValidateAccessToken(token)
	if err == nil {
		t.Fatalf("expected error when validating with wrong secret, got nil")
	}

	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	if claims != nil {
		t.Fatalf("expected nil claims")
	}
}

// TestTokenExpiration tests token expiration
func TestTokenExpiration(t *testing.T) {
	// Create token manager with short TTL
	tm := NewTokenManager("test-secret", 100*time.Millisecond, 7*24*time.Hour)

	token, err := tm.GenerateAccessToken("test@example.com")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Immediately validate - should be valid
	claims, err := tm.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("expected valid token immediately after generation, got error: %v", err)
	}

	if claims == nil {
		t.Fatalf("expected claims for valid token")
	}

	// Wait for token to expire
	time.Sleep(150 * time.Millisecond)

	// Try to validate expired token
	claims, err = tm.ValidateAccessToken(token)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}

	if claims != nil {
		t.Fatalf("expected nil claims for expired token")
	}
}

// TestDifferentTokenTypes tests that tokens have different expirations
func TestDifferentTokenTypes(t *testing.T) {
	tm := NewTokenManager("test-secret", 1*time.Hour, 7*24*time.Hour)

	accessToken, err := tm.GenerateAccessToken("test@example.com")
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	refreshToken, err := tm.GenerateRefreshToken("test@example.com")
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	accessClaims, err := tm.ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}

	refreshClaims, err := tm.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}

	// Access token should expire sooner than refresh token
	if accessClaims.ExpiresAt.After(refreshClaims.ExpiresAt.Time) {
		t.Fatalf("expected access token to expire before refresh token")
	}
}
