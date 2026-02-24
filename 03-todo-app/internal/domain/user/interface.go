package user

import "context"

// Service defines user business logic contracts
type Service interface {
	// Create creates a new user with password hashing
	Create(ctx context.Context, email, name, password string) (map[string]interface{}, error)

	// GetByEmail retrieves user by email (password_hash excluded)
	GetByEmail(ctx context.Context, email string) (map[string]interface{}, error)

	// GetByID retrieves user by ID (password_hash excluded)
	GetByID(ctx context.Context, id string) (map[string]interface{}, error)

	// List returns all active users (paginated)
	List(ctx context.Context, limit, offset int) ([]map[string]interface{}, error)

	// Update updates user profile (name only)
	Update(ctx context.Context, id, name string) error

	// Delete soft-deletes a user (sets is_active = false)
	Delete(ctx context.Context, id string) error

	// VerifyPassword verifies email + password combination
	VerifyPassword(ctx context.Context, email, password string) (map[string]interface{}, error)
}

// Repository defines data access contracts
type Repository interface {
	Create(ctx context.Context, email, name, passwordHash string) (map[string]interface{}, error)
	GetByEmail(ctx context.Context, email string) (map[string]interface{}, error)
	GetByID(ctx context.Context, id string) (map[string]interface{}, error)
	List(ctx context.Context, limit, offset int) ([]map[string]interface{}, error)
	Update(ctx context.Context, id, name string) error
	Delete(ctx context.Context, id string) error
}
