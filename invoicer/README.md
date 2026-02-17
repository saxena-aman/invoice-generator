# Invoice Generator - Go Backend

A Go-based REST API for generating professional PDF invoices with JWT/OAuth2 authentication and rate limiting.

## Features

- ✅ RESTful API for PDF generation
- ✅ Support for item-level tax and discount
- ✅ Support for bill-level tax and discount  
- ✅ Professional PDF layout (minimal, corporate, modern templates)
- ✅ CORS enabled for React frontend
- ✅ JSON request/response
- ✅ Input validation
- ✅ **JWT authentication** (register, login, token refresh)
- ✅ **Google OAuth2** login
- ✅ **Rate limiting** (per-IP for anonymous, per-user for authenticated)

## Project Structure

```
invoicer/
├── main.go                          # HTTP server, routing, middleware wiring
├── internal/
│   ├── auth/
│   │   ├── config.go               # Auth configuration from env vars
│   │   ├── models.go               # User, Claims, request/response types
│   │   ├── jwt.go                  # JWT token generation & validation
│   │   ├── store.go                # In-memory user store with bcrypt
│   │   └── oauth.go                # Google OAuth2 service
│   ├── handlers/
│   │   ├── invoice.go              # Invoice PDF handler
│   │   └── auth_handler.go         # Auth endpoints (register, login, OAuth)
│   ├── middleware/
│   │   ├── auth_middleware.go      # JWT Bearer token validation
│   │   └── rate_limiter.go         # Per-IP / per-user rate limiting
│   ├── models/
│   │   └── invoice.go              # Invoice data models
│   └── pdf/
│       └── generator.go            # PDF generation logic
├── go.mod
└── go.sum
```

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `JWT_SECRET` | **Yes** | — | Secret key for signing JWTs |
| `JWT_EXPIRY_HOURS` | No | `24` | Access token expiry in hours |
| `GOOGLE_CLIENT_ID` | No | — | Google OAuth2 client ID |
| `GOOGLE_CLIENT_SECRET` | No | — | Google OAuth2 client secret |
| `GOOGLE_REDIRECT_URL` | No | — | Google OAuth2 redirect URL |
| `RATE_LIMIT_PER_MIN` | No | `30` | Requests/min for anonymous users |
| `RATE_LIMIT_AUTH_PER_MIN` | No | `60` | Requests/min for authenticated users |
| `ALLOWED_ORIGINS` | No | `localhost:5173,3000` | CORS allowed origins |

## API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/auth/register` | Register with email/password |
| `POST` | `/api/auth/login` | Login with email/password |
| `POST` | `/api/auth/refresh` | Refresh access token |
| `GET`  | `/api/auth/google` | Redirect to Google OAuth consent |
| `GET`  | `/api/auth/google/callback` | Google OAuth callback |

#### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Token Response:**
```json
{
  "accessToken": "eyJhbG...",
  "refreshToken": "eyJhbG...",
  "expiresIn": 86400,
  "tokenType": "Bearer"
}
```

### Generate PDF (🔒 Protected)
**POST** `/api/generate-pdf`

Requires `Authorization: Bearer <accessToken>` header.

```bash
curl -X POST http://localhost:8080/api/generate-pdf \
  -H "Authorization: Bearer <your-access-token>" \
  -H "Content-Type: application/json" \
  -d @test-invoice.json \
  --output invoice.pdf
```

### Health Check (Public)
**GET** `/health`

```bash
curl http://localhost:8080/health
```

### Rate Limiting

All endpoints are rate-limited:
- **Anonymous users**: 30 requests/minute per IP
- **Authenticated users**: 60 requests/minute per user

When rate limited, the API returns `429 Too Many Requests` with a `Retry-After` header.

## Running the Server

### Development
```bash
# Copy env file and set JWT_SECRET
cp .env.example .env
# Edit .env and set a strong JWT_SECRET

# Run
go run main.go
```

### Production Build
```bash
go build -o bin/invoicer
./bin/invoicer
```

## Testing

```bash
# Run all tests
go test ./... -v
```

## Dependencies

- `github.com/jung-kurt/gofpdf` - PDF generation
- `github.com/gorilla/mux` - HTTP routing
- `github.com/rs/cors` - CORS middleware
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `golang.org/x/crypto` - bcrypt password hashing
- `golang.org/x/oauth2` - Google OAuth2
- `golang.org/x/time` - Rate limiting

## License

MIT
