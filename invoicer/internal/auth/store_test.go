package auth

import (
	"testing"
)

func TestUserStore_CreateAndGetUser(t *testing.T) {
	store := NewUserStore()

	user, err := store.CreateUser("test@example.com", "password123", "Test User", "local")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", user.Email)
	}
	if user.PasswordHash == "" {
		t.Error("expected non-empty password hash")
	}
	if user.PasswordHash == "password123" {
		t.Error("password hash should not be the raw password")
	}

	// Get by email
	found, err := store.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("expected ID %q, got %q", user.ID, found.ID)
	}

	// Get by ID
	found2, err := store.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if found2.Email != user.Email {
		t.Errorf("expected email %q, got %q", user.Email, found2.Email)
	}
}

func TestUserStore_DuplicateEmail(t *testing.T) {
	store := NewUserStore()

	_, err := store.CreateUser("dup@example.com", "password123", "User 1", "local")
	if err != nil {
		t.Fatalf("first CreateUser failed: %v", err)
	}

	_, err = store.CreateUser("dup@example.com", "password123", "User 2", "local")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestUserStore_CheckPassword(t *testing.T) {
	store := NewUserStore()

	user, _ := store.CreateUser("pwd@example.com", "correctpassword", "User", "local")

	// Correct password
	if err := store.CheckPassword(user, "correctpassword"); err != nil {
		t.Errorf("expected correct password to pass: %v", err)
	}

	// Wrong password
	if err := store.CheckPassword(user, "wrongpassword"); err == nil {
		t.Error("expected wrong password to fail")
	}
}

func TestUserStore_UpsertOAuthUser(t *testing.T) {
	store := NewUserStore()

	// First upsert — creates user
	user1, err := store.UpsertOAuthUser("oauth@example.com", "OAuth User", "google")
	if err != nil {
		t.Fatalf("UpsertOAuthUser(create) failed: %v", err)
	}

	// Second upsert — returns same user
	user2, err := store.UpsertOAuthUser("oauth@example.com", "OAuth User", "google")
	if err != nil {
		t.Fatalf("UpsertOAuthUser(existing) failed: %v", err)
	}
	if user1.ID != user2.ID {
		t.Errorf("expected same user ID, got %q and %q", user1.ID, user2.ID)
	}
}

func TestUserStore_UserNotFound(t *testing.T) {
	store := NewUserStore()

	_, err := store.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("expected error for non-existent email")
	}

	_, err = store.GetUserByID("user_999")
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}
