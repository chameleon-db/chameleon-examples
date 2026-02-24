package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/todo"
	"github.com/chameleon-db/chameleondb/chameleon/pkg/engine"
	_ "github.com/chameleon-db/chameleondb/chameleon/pkg/engine/mutation"
	"github.com/google/uuid"
)

// TodoRepository implements todo.Repository
type TodoRepository struct {
	engine *engine.Engine
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(eng *engine.Engine) todo.Repository {
	return &TodoRepository{engine: eng}
}

// Create inserts new todo via ChameleonDB
func (r *TodoRepository) Create(ctx context.Context, userID, title, description string) (map[string]interface{}, error) {
	result, err := r.engine.Insert("Todo").
		Set("id", uuid.New().String()).
		Set("user_id", userID).
		Set("title", title).
		Set("description", description).
		Set("completed", false).
		Debug().
		Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("failed to create todo: empty result")
	}

	if result.Record != nil {
		return result.Record, nil
	}

	if result.ID != nil {
		return map[string]interface{}{"id": result.ID}, nil
	}

	return nil, fmt.Errorf("failed to create todo: missing record")
}

// GetByID retrieves todo by ID
func (r *TodoRepository) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	result, err := r.engine.Query("Todo").
		Filter("id", "eq", id).
		Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query todo: %w", err)
	}

	if result == nil || result.IsEmpty() {
		return nil, todo.ErrNotFound
	}

	return rowToMap(result.Rows[0]), nil
}

// ListByUser returns user's todos (paginated)
func (r *TodoRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]map[string]interface{}, error) {
	query := r.engine.Query("Todo").
		Filter("user_id", "eq", userID)

	if limit > 0 {
		query = query.Limit(uint64(limit))
	}
	if offset > 0 {
		query = query.Offset(uint64(offset))
	}

	result, err := query.Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("failed to list todos: empty result")
	}

	return rowsToMaps(result.Rows), nil
}

// ListByUserFiltered returns user's todos with completion filter
func (r *TodoRepository) ListByUserFiltered(ctx context.Context, userID string, completed *bool, limit, offset int) ([]map[string]interface{}, error) {
	query := r.engine.Query("Todo").
		Filter("user_id", "eq", userID)

	if completed != nil {
		query = query.Filter("completed", "eq", *completed)
	}
	if limit > 0 {
		query = query.Limit(uint64(limit))
	}
	if offset > 0 {
		query = query.Offset(uint64(offset))
	}

	result, err := query.Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("failed to list todos: empty result")
	}

	return rowsToMaps(result.Rows), nil
}

// Update updates todo
func (r *TodoRepository) Update(ctx context.Context, id, title, description string, completed bool) error {
	result, err := r.engine.Update("Todo").
		Filter("id", "eq", id).
		Set("title", title).
		Set("description", description).
		Set("completed", completed).
		Debug().
		Execute(ctx)

	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	if result != nil && result.Affected == 0 {
		return todo.ErrNotFound
	}

	return nil
}

// Delete deletes todo
func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	result, err := r.engine.Delete("Todo").
		Filter("id", "eq", id).
		Debug().
		Execute(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if result != nil && result.Affected == 0 {
		return todo.ErrNotFound
	}

	return nil
}

// GetOverdue returns overdue todos for user
func (r *TodoRepository) GetOverdue(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	now := time.Now()

	result, err := r.engine.Query("Todo").
		Filter("user_id", "eq", userID).
		Filter("completed", "eq", false).
		Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query overdue todos: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("failed to query overdue todos: empty result")
	}

	// Filter overdue todos (due_date < now)
	var overdue []map[string]interface{}
	for _, t := range result.Rows {
		dueDate, ok := t["due_date"].(time.Time)
		if ok && dueDate.Before(now) {
			overdue = append(overdue, rowToMap(t))
		}
	}

	return overdue, nil
}

// GetByIDForUser retrieves todo and validates it belongs to user
func (r *TodoRepository) GetByIDForUser(ctx context.Context, id, userID string) (map[string]interface{}, error) {
	result, err := r.engine.Query("Todo").
		Filter("id", "eq", id).
		Filter("user_id", "eq", userID).
		Execute(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query todo: %w", err)
	}

	if result == nil || result.IsEmpty() {
		return nil, todo.ErrNotFound
	}

	return rowToMap(result.Rows[0]), nil
}
