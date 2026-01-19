package main

import (
	"fmt"
	"invoice-generator/invoicer/internal/handlers"
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

	// Create router
	router := mux.NewRouter()

	// Initialize handlers
	invoiceHandler := handlers.NewInvoiceHandler()

	// Register routes
	router.HandleFunc("/api/generate-pdf", invoiceHandler.GeneratePDF).Methods("POST", "OPTIONS")
	router.HandleFunc("/health", invoiceHandler.HealthCheck).Methods("GET")

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
	fmt.Printf("📄 PDF generation endpoint: http://localhost:%s/api/generate-pdf\n", port)
	fmt.Printf("💚 Health check endpoint: http://localhost:%s/health\n", port)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
