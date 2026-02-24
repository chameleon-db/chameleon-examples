package todo

import "context"

// Service defines todo business logic contracts
type Service interface {
	// Create creates a new todo for user
	Create(ctx context.Context, userID, title, description string) (map[string]interface{}, error)

	// GetByID retrieves todo
	GetByID(ctx context.Context, id string) (map[string]interface{}, error)

	// ListByUser returns user's todos (paginated)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]map[string]interface{}, error)

	// ListByUserFiltered returns user's todos with filters
	ListByUserFiltered(ctx context.Context, userID string, completed *bool, limit, offset int) ([]map[string]interface{}, error)

	// Update updates todo
	Update(ctx context.Context, id, title, description string, completed bool) error

	// Delete deletes todo
	Delete(ctx context.Context, id string) error

	// GetOverdue returns overdue todos for user
	GetOverdue(ctx context.Context, userID string) ([]map[string]interface{}, error)

	// ToggleCompletion toggles todo completion status
	ToggleCompletion(ctx context.Context, id string) error
}

// Repository defines data access contracts
type Repository interface {
	Create(ctx context.Context, userID, title, description string) (map[string]interface{}, error)
	GetByID(ctx context.Context, id string) (map[string]interface{}, error)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]map[string]interface{}, error)
	ListByUserFiltered(ctx context.Context, userID string, completed *bool, limit, offset int) ([]map[string]interface{}, error)
	Update(ctx context.Context, id, title, description string, completed bool) error
	Delete(ctx context.Context, id string) error
	GetOverdue(ctx context.Context, userID string) ([]map[string]interface{}, error)
	GetByIDForUser(ctx context.Context, id, userID string) (map[string]interface{}, error)
}
