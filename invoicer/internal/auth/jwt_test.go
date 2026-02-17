package auth

import (
	"testing"
	"time"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	svc := NewJWTService("test-secret-key", time.Hour, 7*24*time.Hour)

	user := &User{
		ID:    "user_1",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Generate access token
	token, err := svc.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// Validate token
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != user.ID {
		t.Errorf("expected UserID %q, got %q", user.ID, claims.UserID)
	}
	if claims.Email != user.Email {
		t.Errorf("expected Email %q, got %q", user.Email, claims.Email)
	}
	if claims.Type != "access" {
		t.Errorf("expected Type 'access', got %q", claims.Type)
	}
}

func TestJWTService_GenerateAndValidateRefreshToken(t *testing.T) {
	svc := NewJWTService("test-secret-key", time.Hour, 7*24*time.Hour)

	user := &User{ID: "user_1", Email: "test@example.com"}

	token, err := svc.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.Type != "refresh" {
		t.Errorf("expected Type 'refresh', got %q", claims.Type)
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	// Use a negative expiry so the token is immediately expired
	svc := NewJWTService("test-secret-key", -1*time.Hour, -1*time.Hour)

	user := &User{ID: "user_1", Email: "test@example.com"}
	token, err := svc.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Fatal("expected validation error for expired token")
	}
}

func TestJWTService_TamperedToken(t *testing.T) {
	svc := NewJWTService("test-secret-key", time.Hour, 7*24*time.Hour)

	user := &User{ID: "user_1", Email: "test@example.com"}
	token, err := svc.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Tamper with the token
	tampered := token + "tampered"
	_, err = svc.ValidateToken(tampered)
	if err == nil {
		t.Fatal("expected validation error for tampered token")
	}
}

func TestJWTService_WrongSecret(t *testing.T) {
	svc1 := NewJWTService("secret-one", time.Hour, 7*24*time.Hour)
	svc2 := NewJWTService("secret-two", time.Hour, 7*24*time.Hour)

	user := &User{ID: "user_1", Email: "test@example.com"}
	token, err := svc1.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = svc2.ValidateToken(token)
	if err == nil {
		t.Fatal("expected validation error for wrong secret")
	}
}
