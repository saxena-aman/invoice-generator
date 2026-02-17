package auth

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserStore is a thread-safe in-memory user store.
type UserStore struct {
	mu      sync.RWMutex
	users   map[string]*User  // keyed by user ID
	byEmail map[string]string // email -> user ID index
	nextID  int
}

// NewUserStore creates an empty user store.
func NewUserStore() *UserStore {
	return &UserStore{
		users:   make(map[string]*User),
		byEmail: make(map[string]string),
	}
}

// CreateUser hashes the password and stores a new user. Returns error if email is taken.
func (s *UserStore) CreateUser(email, password, name, provider string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byEmail[email]; exists {
		return nil, fmt.Errorf("email already registered")
	}

	var hash string
	if password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hash = string(h)
	}

	s.nextID++
	id := fmt.Sprintf("user_%d", s.nextID)

	user := &User{
		ID:           id,
		Email:        email,
		Name:         name,
		PasswordHash: hash,
		Provider:     provider,
		CreatedAt:    time.Now(),
	}

	s.users[id] = user
	s.byEmail[email] = id
	return user, nil
}

// GetUserByEmail looks up a user by email.
func (s *UserStore) GetUserByEmail(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, exists := s.byEmail[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return s.users[id], nil
}

// GetUserByID looks up a user by ID.
func (s *UserStore) GetUserByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// CheckPassword verifies a plaintext password against the stored hash.
func (s *UserStore) CheckPassword(user *User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

// UpsertOAuthUser creates or retrieves an existing OAuth user (no password).
func (s *UserStore) UpsertOAuthUser(email, name, provider string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id, exists := s.byEmail[email]; exists {
		return s.users[id], nil
	}

	s.nextID++
	id := fmt.Sprintf("user_%d", s.nextID)

	user := &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Provider:  provider,
		CreatedAt: time.Now(),
	}

	s.users[id] = user
	s.byEmail[email] = id
	return user, nil
}
