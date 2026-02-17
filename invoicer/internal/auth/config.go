package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// AuthConfig holds authentication configuration loaded from environment variables.
type AuthConfig struct {
	// JWT
	JWTSecret          string
	JWTExpiry          time.Duration
	JWTRefreshExpiry   time.Duration

	// OAuth2 (Google)
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// Rate limiting
	RateLimitPerMin     int
	RateLimitAuthPerMin int
}

// LoadAuthConfig reads auth configuration from environment variables.
func LoadAuthConfig() (*AuthConfig, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	expiry := 24 * time.Hour
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		hours, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
		}
		expiry = time.Duration(hours) * time.Hour
	}

	rateLimit := 30
	if v := os.Getenv("RATE_LIMIT_PER_MIN"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT_PER_MIN: %w", err)
		}
		rateLimit = n
	}

	rateLimitAuth := 60
	if v := os.Getenv("RATE_LIMIT_AUTH_PER_MIN"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT_AUTH_PER_MIN: %w", err)
		}
		rateLimitAuth = n
	}

	return &AuthConfig{
		JWTSecret:           secret,
		JWTExpiry:           expiry,
		JWTRefreshExpiry:    7 * 24 * time.Hour, // 7 days
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:   os.Getenv("GOOGLE_REDIRECT_URL"),
		RateLimitPerMin:     rateLimit,
		RateLimitAuthPerMin: rateLimitAuth,
	}, nil
}
