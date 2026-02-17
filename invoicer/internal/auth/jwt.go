package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles token generation and validation.
type JWTService struct {
	secret        []byte
	expiry        time.Duration
	refreshExpiry time.Duration
}

// NewJWTService creates a new JWT service.
func NewJWTService(secret string, expiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secret:        []byte(secret),
		expiry:        expiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateToken creates a signed access token for the given user.
func (s *JWTService) GenerateToken(user *User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID,
			Issuer:    "invoice-generator",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateRefreshToken creates a signed refresh token for the given user.
func (s *JWTService) GenerateRefreshToken(user *User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID,
			Issuer:    "invoice-generator",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken parses and validates the given token string.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GetExpiry returns the access token expiry duration.
func (s *JWTService) GetExpiry() time.Duration {
	return s.expiry
}
