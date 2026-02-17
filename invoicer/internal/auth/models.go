package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents a registered user.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`        // never serialised to JSON
	Provider     string    `json:"provider"` // "local" or "google"
	CreatedAt    time.Time `json:"createdAt"`
}

// Claims are the JWT claims embedded in access and refresh tokens.
type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// RegisterRequest is the body for POST /api/auth/register.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest is the body for POST /api/auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest is the body for POST /api/auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// TokenResponse is returned after successful authentication.
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // seconds until access token expiry
	TokenType    string `json:"tokenType"` // always "Bearer"
}

// ErrorResponse is a standard JSON error envelope.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
