module invoice-generator/invoicer

go 1.24.5

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/rs/cors v1.11.1
	golang.org/x/crypto v0.48.0
	golang.org/x/oauth2 v0.35.0
	golang.org/x/time v0.14.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
