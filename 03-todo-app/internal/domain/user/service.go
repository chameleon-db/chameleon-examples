package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

// userService implements the Service interface
type userService struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &userService{repo: repo}
}

// Create creates a new user with password hashing
func (s *userService) Create(ctx context.Context, email, name, password string) (map[string]interface{}, error) {
	// Validation
	if email == "" || name == "" || password == "" {
		return nil, ErrInvalidInput
	}

	if len(password) < 8 {
		return nil, ErrWeakPassword
	}

	// Hash password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create via repository
	user, err := s.repo.Create(ctx, email, name, string(hash))
	if err != nil {
		// Check for duplicate email error
		// This is database-level, we bubble up
		return nil, err
	}

	// Remove sensitive fields
	delete(user, "password_hash")

	return user, nil
}

// GetByEmail retrieves user by email (no password returned)
func (s *userService) GetByEmail(ctx context.Context, email string) (map[string]interface{}, error) {
	if email == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	// Remove sensitive fields
	delete(user, "password_hash")

	return user, nil
}

// GetByID retrieves user by ID (no password returned)
func (s *userService) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	// Remove sensitive fields
	delete(user, "password_hash")

	return user, nil
}

// List returns all active users (paginated)
func (s *userService) List(ctx context.Context, limit, offset int) ([]map[string]interface{}, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Remove sensitive fields from all users
	for _, user := range users {
		delete(user, "password_hash")
	}

	return users, nil
}

// Update updates user profile (name only)
func (s *userService) Update(ctx context.Context, id, name string) error {
	if id == "" || name == "" {
		return ErrInvalidInput
	}

	return s.repo.Update(ctx, id, name)
}

// Delete soft-deletes a user
func (s *userService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}

// VerifyPassword verifies email + password combination
func (s *userService) VerifyPassword(ctx context.Context, email, password string) (map[string]interface{}, error) {
	if email == "" || password == "" {
		return nil, ErrInvalidInput
	}

	// Get user with password hash
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	if user == nil {
		return nil, ErrInvalidPassword
	}

	// Check if user is active
	isActive, ok := user["is_active"].(bool)
	if !ok || !isActive {
		return nil, ErrUserInactive
	}

	// Compare password with hash
	passwordHash, ok := user["password_hash"].(string)
	if !ok {
		return nil, ErrInvalidPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// Return user without password
	delete(user, "password_hash")

	return user, nil
}
