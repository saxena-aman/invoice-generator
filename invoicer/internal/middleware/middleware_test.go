package middleware

import (
	"invoice-generator/invoicer/internal/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", time.Hour, 7*24*time.Hour)
	mw := AuthMiddleware(jwtService)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", time.Hour, 7*24*time.Hour)
	mw := AuthMiddleware(jwtService)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", time.Hour, 7*24*time.Hour)
	user := &auth.User{ID: "user_1", Email: "test@example.com"}

	token, err := jwtService.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	var handlerCalled bool
	mw := AuthMiddleware(jwtService)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		claims := GetClaims(r)
		if claims == nil {
			t.Error("expected claims in context")
			return
		}
		if claims.UserID != "user_1" {
			t.Errorf("expected UserID 'user_1', got %q", claims.UserID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Error("handler should have been called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestAuthMiddleware_RefreshTokenRejected(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", time.Hour, 7*24*time.Hour)
	user := &auth.User{ID: "user_1", Email: "test@example.com"}

	// Generate a refresh token (should be rejected by auth middleware)
	refreshToken, err := jwtService.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	mw := AuthMiddleware(jwtService)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called with refresh token")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestRateLimiter_AllowsUnderLimit(t *testing.T) {
	rl := NewRateLimiter(100, 200) // high limit for test

	handler := rl.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRateLimiter_RejectsOverLimit(t *testing.T) {
	rl := NewRateLimiter(1, 1) // Extremely low: 1 req/min

	handler := rl.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request should pass (burst allows 1)
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "10.0.0.1:12345"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("first request: expected 200, got %d", rr1.Code)
	}

	// Second request should be rate limited
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "10.0.0.1:12345"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: expected 429, got %d", rr2.Code)
	}
}
