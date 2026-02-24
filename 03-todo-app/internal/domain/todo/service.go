package todo

import "context"

// todoService implements the Service interface
type todoService struct {
	repo Repository
}

// NewService creates a new todo service
func NewService(repo Repository) Service {
	return &todoService{repo: repo}
}

// Create creates a new todo for user
func (s *todoService) Create(ctx context.Context, userID, title, description string) (map[string]interface{}, error) {
	// Validation
	if userID == "" || title == "" {
		return nil, ErrInvalidInput
	}

	// Create via repository
	todo, err := s.repo.Create(ctx, userID, title, description)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// GetByID retrieves todo
func (s *todoService) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, ErrNotFound
	}

	return todo, nil
}

// ListByUser returns user's todos (paginated)
func (s *todoService) ListByUser(ctx context.Context, userID string, limit, offset int) ([]map[string]interface{}, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}

	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	todos, err := s.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

// ListByUserFiltered returns user's todos with filters
func (s *todoService) ListByUserFiltered(ctx context.Context, userID string, completed *bool, limit, offset int) ([]map[string]interface{}, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}

	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	todos, err := s.repo.ListByUserFiltered(ctx, userID, completed, limit, offset)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

// Update updates todo
func (s *todoService) Update(ctx context.Context, id, title, description string, completed bool) error {
	if id == "" || title == "" {
		return ErrInvalidInput
	}

	return s.repo.Update(ctx, id, title, description, completed)
}

// Delete deletes todo
func (s *todoService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}

// GetOverdue returns overdue todos for user
func (s *todoService) GetOverdue(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}

	todos, err := s.repo.GetOverdue(ctx, userID)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

// ToggleCompletion toggles todo completion status
func (s *todoService) ToggleCompletion(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	// Get current todo
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if todo == nil {
		return ErrNotFound
	}

	// Toggle completion
	completed, ok := todo["completed"].(bool)
	if !ok {
		return ErrInvalidInput
	}

	// Get other fields
	title, _ := todo["title"].(string)
	description, _ := todo["description"].(string)

	// Update with toggled value
	return s.repo.Update(ctx, id, title, description, !completed)
}
