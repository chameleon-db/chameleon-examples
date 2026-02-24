package repository

import (
	"context"
	"fmt"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/user"
	"github.com/chameleon-db/chameleondb/chameleon/pkg/engine"
	_ "github.com/chameleon-db/chameleondb/chameleon/pkg/engine/mutation"
)

// UserRepository implements user.Repository
type UserRepository struct {
	engine *engine.Engine
}

// NewUserRepository creates a new user repository
func NewUserRepository(eng *engine.Engine) user.Repository {
	return &UserRepository{engine: eng}
}

// Create inserts new user via ChameleonDB
func (r *UserRepository) Create(ctx context.Context, email, name, passwordHash string) (map[string]interface{}, error) {
	result, err := r.engine.Insert("User").
		Set("email", email).
		Set("name", name).
		Set("password_hash", passwordHash).
		Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("failed to create user: empty result")
	}

	if result.Record != nil {
		return result.Record, nil
	}

	if result.ID != nil {
		return map[string]interface{}{"id": result.ID}, nil
	}

	return nil, fmt.Errorf("failed to create user: missing record")
}

// GetByEmail retrieves user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (map[string]interface{}, error) {
	result, err := r.engine.Query("User").
		Filter("email", "eq", email).
		Filter("is_active", "eq", true).
		Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if result == nil || result.IsEmpty() {
		return nil, user.ErrNotFound
	}

	return rowToMap(result.Rows[0]), nil
}

// GetByID retrieves user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	result, err := r.engine.Query("User").
		Filter("id", "eq", id).
		Filter("is_active", "eq", true).
		Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if result == nil || result.IsEmpty() {
		return nil, user.ErrNotFound
	}

	return rowToMap(result.Rows[0]), nil
}

// List returns all active users (paginated)
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]map[string]interface{}, error) {
	query := r.engine.Query("User").
		Filter("is_active", "eq", true)

	if limit > 0 {
		query = query.Limit(uint64(limit))
	}
	if offset > 0 {
		query = query.Offset(uint64(offset))
	}

	result, err := query.Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("failed to list users: empty result")
	}

	return rowsToMaps(result.Rows), nil
}

// Update updates user (name only)
func (r *UserRepository) Update(ctx context.Context, id, name string) error {
	result, err := r.engine.Update("User").
		Filter("id", "eq", id).
		Set("name", name).
		Execute(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result != nil && result.Affected == 0 {
		return user.ErrNotFound
	}

	return nil
}

// Delete soft-deletes user (sets is_active = false)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.engine.Update("User").
		Filter("id", "eq", id).
		Set("is_active", false).
		Execute(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result != nil && result.Affected == 0 {
		return user.ErrNotFound
	}

	return nil
}
