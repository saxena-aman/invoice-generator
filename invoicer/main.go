package main

import (
	"fmt"
	"invoice-generator/invoicer/internal/auth"
	"invoice-generator/invoicer/internal/handlers"
	"invoice-generator/invoicer/internal/middleware"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Load auth configuration
	authConfig, err := auth.LoadAuthConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load auth config: %v", err)
	}

	// Initialize auth services
	jwtService := auth.NewJWTService(authConfig.JWTSecret, authConfig.JWTExpiry, authConfig.JWTRefreshExpiry)
	userStore := auth.NewUserStore()
	oauthService := auth.NewOAuthService(
		authConfig.GoogleClientID,
		authConfig.GoogleClientSecret,
		authConfig.GoogleRedirectURL,
		userStore,
	)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(authConfig.RateLimitPerMin, authConfig.RateLimitAuthPerMin)

	// Create router
	router := mux.NewRouter()

	// Apply rate limiting globally
	router.Use(rateLimiter.Middleware())

	// Initialize handlers
	invoiceHandler := handlers.NewInvoiceHandler()
	authHandler := handlers.NewAuthHandler(jwtService, userStore, oauthService)

	// ── Public routes (no auth required) ─────────────────────────────
	router.HandleFunc("/health", invoiceHandler.HealthCheck).Methods("GET")

	// Auth routes
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")
	authRouter.HandleFunc("/google", authHandler.GoogleLogin).Methods("GET")
	authRouter.HandleFunc("/google/callback", authHandler.GoogleCallback).Methods("GET")

	// ── Protected routes (JWT auth required) ─────────────────────────
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(jwtService))
	protectedRouter.HandleFunc("/generate-pdf", invoiceHandler.GeneratePDF).Methods("POST", "OPTIONS")

	// Get allowed origins from environment
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
	} else {
		allowedOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap router with CORS middleware
	handler := corsHandler.Handler(router)

	// Start server
	port := "8080"
	fmt.Printf("🚀 Invoice Generator API server starting on port %s\n", port)
	fmt.Printf("📄 PDF generation endpoint: http://localhost:%s/api/generate-pdf (🔒 protected)\n", port)
	fmt.Printf("🔑 Auth endpoints:          http://localhost:%s/api/auth/*\n", port)
	fmt.Printf("💚 Health check endpoint:    http://localhost:%s/health\n", port)
	if oauthService != nil {
		fmt.Printf("🌐 Google OAuth:             http://localhost:%s/api/auth/google\n", port)
	} else {
		fmt.Println("⚠️  Google OAuth is not configured (set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET)")
	}
	fmt.Printf("🛡️  Rate limiting:            %d req/min (anonymous), %d req/min (authenticated)\n",
		authConfig.RateLimitPerMin, authConfig.RateLimitAuthPerMin)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
