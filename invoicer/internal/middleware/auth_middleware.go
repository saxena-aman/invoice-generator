package middleware

import (
	"context"
	"encoding/json"
	"invoice-generator/invoicer/internal/auth"
	"net/http"
	"strings"
)

// contextKey is an unexported type for context keys to avoid collisions.
type contextKey string

const (
	// UserClaimsKey is the context key for the authenticated user's claims.
	UserClaimsKey contextKey = "userClaims"
)

// AuthMiddleware returns an HTTP middleware that validates JWT Bearer tokens.
// Requests without a valid token receive a 401 Unauthorized response.
func AuthMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
					Error:   "unauthorized",
					Message: "Authorization header is required",
				})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
					Error:   "unauthorized",
					Message: "Authorization header must be in the format: Bearer <token>",
				})
				return
			}

			claims, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
					Error:   "unauthorized",
					Message: "Invalid or expired token",
				})
				return
			}

			if claims.Type != "access" {
				writeJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
					Error:   "unauthorized",
					Message: "Invalid token type",
				})
				return
			}

			// Inject claims into request context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims extracts user claims from the request context.
func GetClaims(r *http.Request) *auth.Claims {
	claims, _ := r.Context().Value(UserClaimsKey).(*auth.Claims)
	return claims
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
