package handlers

import (
	"encoding/json"
	"invoice-generator/invoicer/internal/auth"
	"net/http"
	"regexp"
	"strings"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	jwtService   *auth.JWTService
	userStore    *auth.UserStore
	oauthService *auth.OAuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(jwtService *auth.JWTService, userStore *auth.UserStore, oauthService *auth.OAuthService) *AuthHandler {
	return &AuthHandler{
		jwtService:   jwtService,
		userStore:    userStore,
		oauthService: oauthService,
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid JSON body",
		})
		return
	}
	defer r.Body.Close()

	// Validate
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if !emailRegex.MatchString(req.Email) {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid email address",
		})
		return
	}
	if len(req.Password) < 8 {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "validation_error",
			Message: "Password must be at least 8 characters",
		})
		return
	}
	if req.Name == "" {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "validation_error",
			Message: "Name is required",
		})
		return
	}

	// Create user
	user, err := h.userStore.CreateUser(req.Email, req.Password, req.Name, "local")
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			writeAuthJSON(w, http.StatusConflict, auth.ErrorResponse{
				Error:   "conflict",
				Message: "Email already registered",
			})
			return
		}
		writeAuthJSON(w, http.StatusInternalServerError, auth.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create user",
		})
		return
	}

	// Generate tokens
	tokenResponse, err := h.generateTokens(user)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, auth.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate tokens",
		})
		return
	}

	writeAuthJSON(w, http.StatusCreated, tokenResponse)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid JSON body",
		})
		return
	}
	defer r.Body.Close()

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Look up user
	user, err := h.userStore.GetUserByEmail(req.Email)
	if err != nil {
		writeAuthJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid email or password",
		})
		return
	}

	// Check password
	if err := h.userStore.CheckPassword(user, req.Password); err != nil {
		writeAuthJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid email or password",
		})
		return
	}

	// Generate tokens
	tokenResponse, err := h.generateTokens(user)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, auth.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate tokens",
		})
		return
	}

	writeAuthJSON(w, http.StatusOK, tokenResponse)
}

// Refresh handles POST /api/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid JSON body",
		})
		return
	}
	defer r.Body.Close()

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil || claims.Type != "refresh" {
		writeAuthJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid or expired refresh token",
		})
		return
	}

	// Look up user
	user, err := h.userStore.GetUserByID(claims.UserID)
	if err != nil {
		writeAuthJSON(w, http.StatusUnauthorized, auth.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not found",
		})
		return
	}

	// Generate new tokens
	tokenResponse, err := h.generateTokens(user)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, auth.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate tokens",
		})
		return
	}

	writeAuthJSON(w, http.StatusOK, tokenResponse)
}

// GoogleLogin handles GET /api/auth/google — redirects to Google consent screen.
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.oauthService == nil {
		writeAuthJSON(w, http.StatusServiceUnavailable, auth.ErrorResponse{
			Error:   "unavailable",
			Message: "Google OAuth is not configured",
		})
		return
	}

	state, err := h.oauthService.GenerateStateToken()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, auth.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate state token",
		})
		return
	}

	// Store state in a cookie for CSRF validation
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	url := h.oauthService.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles GET /api/auth/google/callback
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if h.oauthService == nil {
		writeAuthJSON(w, http.StatusServiceUnavailable, auth.ErrorResponse{
			Error:   "unavailable",
			Message: "Google OAuth is not configured",
		})
		return
	}

	// Validate state for CSRF protection
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid OAuth state",
		})
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		writeAuthJSON(w, http.StatusBadRequest, auth.ErrorResponse{
			Error:   "bad_request",
			Message: "Authorization code is required",
		})
		return
	}

	// Exchange code and get/create user
	user, err := h.oauthService.HandleCallback(r.Context(), code)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, auth.ErrorResponse{
			Error:   "internal_error",
			Message: "OAuth authentication failed",
		})
		return
	}

	// Generate tokens
	tokenResponse, err := h.generateTokens(user)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, auth.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate tokens",
		})
		return
	}

	writeAuthJSON(w, http.StatusOK, tokenResponse)
}

// generateTokens creates both access and refresh tokens for a user.
func (h *AuthHandler) generateTokens(user *auth.User) (*auth.TokenResponse, error) {
	accessToken, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &auth.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.jwtService.GetExpiry().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func writeAuthJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
