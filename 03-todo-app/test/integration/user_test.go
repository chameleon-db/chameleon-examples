package integration

import (
	"context"
	"testing"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/user"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/repository"
)

// TestUserCreate tests user creation
func TestUserCreate(t *testing.T) {
	// Setup
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	// Test
	ctx := context.Background()
	u, err := svc.Create(ctx, "test@example.com", "Test User", "password123")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if u == nil {
		t.Fatal("Expected user, got nil")
	}

	if email, ok := u["email"].(string); !ok || email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", u["email"])
	}

	if _, ok := u["password_hash"]; ok {
		t.Error("Password hash should not be returned")
	}
}

// TestUserCreateInvalidInput tests invalid input
func TestUserCreateInvalidInput(t *testing.T) {
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	ctx := context.Background()

	// Empty email
	_, err := svc.Create(ctx, "", "Test", "password123")
	if err != user.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}

	// Weak password
	_, err = svc.Create(ctx, "test@example.com", "Test", "short")
	if err != user.ErrWeakPassword {
		t.Errorf("Expected ErrWeakPassword, got %v", err)
	}
}

// TestUserGetByEmail tests getting user by email
func TestUserGetByEmail(t *testing.T) {
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	ctx := context.Background()

	// Create user
	_, err := svc.Create(ctx, "test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Retrieve user
	u, err := svc.GetByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if u == nil {
		t.Fatal("Expected user, got nil")
	}

	if email, ok := u["email"].(string); !ok || email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", u["email"])
	}
}

// TestUserGetByEmailNotFound tests getting non-existent user
func TestUserGetByEmailNotFound(t *testing.T) {
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	ctx := context.Background()

	u, err := svc.GetByEmail(ctx, "nonexistent@example.com")
	if err != user.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	if u != nil {
		t.Errorf("Expected nil, got %v", u)
	}
}

// TestUserVerifyPassword tests password verification
func TestUserVerifyPassword(t *testing.T) {
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	ctx := context.Background()

	// Create user
	_, err := svc.Create(ctx, "test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify correct password
	u, err := svc.VerifyPassword(ctx, "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if u == nil {
		t.Fatal("Expected user, got nil")
	}

	// Verify incorrect password
	_, err = svc.VerifyPassword(ctx, "test@example.com", "wrongpassword")
	if err != user.ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got %v", err)
	}
}

// TestUserList tests listing users
func TestUserList(t *testing.T) {
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	ctx := context.Background()

	// Create multiple users
	for i := 1; i <= 3; i++ {
		email := "user" + string(rune(48+i)) + "@example.com"
		_, err := svc.Create(ctx, email, "User "+string(rune(48+i)), "password123")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	// List users
	users, err := svc.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(users) < 3 {
		t.Errorf("Expected at least 3 users, got %d", len(users))
	}
}

// TestUserUpdate tests updating user
func TestUserUpdate(t *testing.T) {
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	ctx := context.Background()

	// Create user
	u, err := svc.Create(ctx, "test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID := u["id"].(string)

	// Update user
	err = svc.Update(ctx, userID, "Updated Name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update
	updated, err := svc.GetByID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if name, ok := updated["name"].(string); !ok || name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %v", updated["name"])
	}
}

// TestUserDelete tests soft-deleting user
func TestUserDelete(t *testing.T) {
	eng := setupTestEngine(t)
	repo := repository.NewUserRepository(eng)
	svc := user.NewService(repo)

	ctx := context.Background()

	// Create user
	u, err := svc.Create(ctx, "test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID := u["id"].(string)

	// Delete user
	err = svc.Delete(ctx, userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify user is not found
	_, err = svc.GetByID(ctx, userID)
	if err != user.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
