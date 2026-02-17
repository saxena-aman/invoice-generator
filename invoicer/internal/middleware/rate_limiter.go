package middleware

import (
	"encoding/json"
	"invoice-generator/invoicer/internal/auth"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// visitor holds the rate limiter and last-seen time for a single visitor.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter provides per-IP and per-user rate limiting.
type RateLimiter struct {
	mu              sync.RWMutex
	visitors        map[string]*visitor
	ratePerMin      int
	authRatePerMin  int
	cleanupInterval time.Duration
}

// NewRateLimiter creates a rate limiter with the specified limits.
// It starts a background goroutine that removes stale entries every 3 minutes.
func NewRateLimiter(ratePerMin, authRatePerMin int) *RateLimiter {
	rl := &RateLimiter{
		visitors:        make(map[string]*visitor),
		ratePerMin:      ratePerMin,
		authRatePerMin:  authRatePerMin,
		cleanupInterval: 3 * time.Minute,
	}

	go rl.cleanup()
	return rl
}

// cleanup periodically removes visitors that haven't been seen recently.
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.cleanupInterval)
		rl.mu.Lock()
		for key, v := range rl.visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

// getVisitor retrieves or creates a rate limiter for the given key and rate.
func (rl *RateLimiter) getVisitor(key string, perMin int) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[key]
	if !exists {
		// Token bucket: rate = perMin/60 tokens per second, burst = perMin
		limiter := rate.NewLimiter(rate.Limit(float64(perMin)/60.0), perMin)
		rl.visitors[key] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// Middleware returns an HTTP middleware that applies rate limiting.
// Authenticated users (with claims in context) get a higher limit keyed by user ID.
// Unauthenticated users are limited by IP address.
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var key string
			perMin := rl.ratePerMin

			// Check if user is authenticated (claims injected by auth middleware)
			if claims := GetClaims(r); claims != nil {
				key = "user:" + claims.UserID
				perMin = rl.authRatePerMin
			} else {
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					ip = r.RemoteAddr
				}
				key = "ip:" + ip
			}

			limiter := rl.getVisitor(key, perMin)
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(auth.ErrorResponse{
					Error:   "too_many_requests",
					Message: "Rate limit exceeded. Please try again later.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
