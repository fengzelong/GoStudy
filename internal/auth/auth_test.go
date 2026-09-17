package auth

import (
	"errors"
	"strings"
	"testing"
	"time"

	"GoStudy/internal/domain"
)

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("secret1")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	ok, err := VerifyPassword(hash, "secret1")
	if err != nil {
		t.Fatalf("verify password: %v", err)
	}
	if !ok {
		t.Fatal("expected password match")
	}

	ok, err = VerifyPassword(hash, "wrong")
	if err != nil {
		t.Fatalf("verify wrong password: %v", err)
	}
	if ok {
		t.Fatal("expected password mismatch")
	}
}

func TestTokenGenerateAndParse(t *testing.T) {
	manager := NewManager("test-secret", time.Hour)

	token, err := manager.Generate(7, domain.UserRoleAdmin)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Subject != 7 {
		t.Fatalf("expected subject 7, got %d", claims.Subject)
	}
	if claims.Role != domain.UserRoleAdmin {
		t.Fatalf("expected admin role, got %s", claims.Role)
	}
	if claims.TokenID == "" {
		t.Fatal("expected JWT token ID")
	}
	if len(strings.Split(token, ".")) != 3 {
		t.Fatalf("expected JWT with three parts, got %s", token)
	}
}

func TestTokenRefreshAndRevoke(t *testing.T) {
	manager := NewManager("test-secret", time.Hour)
	token, err := manager.Generate(7)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	refreshed, err := manager.Refresh(claims)
	if err != nil {
		t.Fatalf("refresh token: %v", err)
	}
	if refreshed == token {
		t.Fatal("expected refreshed token to differ")
	}
	if _, err := manager.Parse(token); !errors.Is(err, ErrRevokedToken) {
		t.Fatalf("expected revoked old token, got %v", err)
	}
	refreshedClaims, err := manager.Parse(refreshed)
	if err != nil {
		t.Fatalf("parse refreshed token: %v", err)
	}
	manager.Revoke(refreshedClaims)
	if _, err := manager.Parse(refreshed); !errors.Is(err, ErrRevokedToken) {
		t.Fatalf("expected revoked refreshed token, got %v", err)
	}
}

func TestTokenRejectsTampering(t *testing.T) {
	manager := NewManager("test-secret", time.Hour)
	token, err := manager.Generate(7)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	parts := strings.Split(token, ".")
	parts[1] = parts[1] + "x"
	if _, err := manager.Parse(strings.Join(parts, ".")); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected invalid token, got %v", err)
	}
}

func TestTokenExpired(t *testing.T) {
	manager := &Manager{secret: []byte("test-secret"), ttl: -time.Hour}

	token, err := manager.Generate(7)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	_, err = manager.Parse(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("expected expired token, got %v", err)
	}
}
