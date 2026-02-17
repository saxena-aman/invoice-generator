package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// googleUserInfo represents the response from Google's userinfo endpoint.
type googleUserInfo struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
}

// OAuthService wraps Google OAuth2 configuration.
type OAuthService struct {
	config *oauth2.Config
	store  *UserStore
}

// NewOAuthService creates a new OAuth service. Returns nil if Google credentials are not configured.
func NewOAuthService(clientID, clientSecret, redirectURL string, store *UserStore) *OAuthService {
	if clientID == "" || clientSecret == "" {
		return nil
	}

	return &OAuthService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		store: store,
	}
}

// GenerateStateToken creates a random state parameter for CSRF protection.
func (s *OAuthService) GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetAuthURL returns the Google consent page URL.
func (s *OAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleCallback exchanges the authorization code for a token, fetches user info,
// and upserts the user in the store.
func (s *OAuthService) HandleCallback(ctx context.Context, code string) (*User, error) {
	// Exchange code for token
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Fetch user info from Google
	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var userInfo googleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	if !userInfo.VerifiedEmail {
		return nil, fmt.Errorf("email not verified by Google")
	}

	// Upsert user
	user, err := s.store.UpsertOAuthUser(userInfo.Email, userInfo.Name, "google")
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return user, nil
}
