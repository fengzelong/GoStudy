package auth

import (
	"errors"
	"testing"
	"time"
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

	token, err := manager.Generate(7)
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
