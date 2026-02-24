package integration

import (
	"context"
	"testing"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/todo"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/user"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/repository"
)

// TestTodoCreate tests todo creation
func TestTodoCreate(t *testing.T) {
	eng := setupTestEngine(t)
	userRepo := repository.NewUserRepository(eng)
	todoRepo := repository.NewTodoRepository(eng)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	ctx := context.Background()

	// Create user
	u, err := userSvc.Create(ctx, "test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID := u["id"].(string)

	// Create todo
	t1, err := todoSvc.Create(ctx, userID, "Test Todo", "This is a test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if t1 == nil {
		t.Fatal("Expected todo, got nil")
	}

	if title, ok := t1["title"].(string); !ok || title != "Test Todo" {
		t.Errorf("Expected title 'Test Todo', got %v", t1["title"])
	}

	if completed, ok := t1["completed"].(bool); !ok || completed {
		t.Errorf("Expected completed false, got %v", t1["completed"])
	}
}

// TestTodoCreateInvalidInput tests invalid input
func TestTodoCreateInvalidInput(t *testing.T) {
	eng := setupTestEngine(t)
	todoRepo := repository.NewTodoRepository(eng)
	todoSvc := todo.NewService(todoRepo)

	ctx := context.Background()

	// Empty title
	_, err := todoSvc.Create(ctx, "user-id", "", "Description")
	if err != todo.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}

	// Empty user ID
	_, err = todoSvc.Create(ctx, "", "Title", "Description")
	if err != todo.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
}

// TestTodoListByUser tests listing user's todos
func TestTodoListByUser(t *testing.T) {
	eng := setupTestEngine(t)
	userRepo := repository.NewUserRepository(eng)
	todoRepo := repository.NewTodoRepository(eng)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	ctx := context.Background()

	// Create user
	u, err := userSvc.Create(ctx, "test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID := u["id"].(string)

	// Create multiple todos
	for i := 1; i <= 3; i++ {
		title := "Todo " + string(rune(48+i))
		_, err := todoSvc.Create(ctx, userID, title, "Description")
		if err != nil {
			t.Fatalf("Failed to create todo: %v", err)
		}
	}

	// List todos
	todos, err := todoSvc.ListByUser(ctx, userID, 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(todos) < 3 {
		t.Errorf("Expected at least 3 todos, got %d", len(todos))
	}
}

// TestTodoListByUserFiltered tests filtering todos
func TestTodoListByUserFiltered(t *testing.T) {
	eng := setupTestEngine(t)
	userRepo := repository.NewUserRepository(eng)
	todoRepo := repository.NewTodoRepository(eng)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	ctx := context.Background()

	// Create user
	u, err := userSvc.Create(ctx, "test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID := u["id"].(string)

	// Create todos
	t1, _ := todoSvc.Create(ctx, userID, "Todo 1", "Description")
	todoID1 := t1["id"].(string)

	_, _ = todoSvc.Create(ctx, userID, "Todo 2", "Description")

	// Mark first todo as complete
	todoSvc.Update(ctx, todoID1, "Todo 1", "Description", true)

	// Filter by completed=true
	completed := true
	todos, err := todoSvc.ListByUserFiltered(ctx, userID, &completed, 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(todos) < 1 {
		t.Errorf("Expected at least 1 completed todo, got %d", len(todos))
	}

	// Filter by completed=false
	completed = false
	todos, err = todoSvc.ListByUserFiltered(ctx, userID, &completed, 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(todos) < 1 {
		t.Errorf("Expected at least 1 incomplete todo, got %d", len(todos))
	}
}

// TestTodoUpdate tests updating todo
func TestTodoUpdate(t *testing.T) {
	eng := setupTestEngine(t)
	userRepo := repository.NewUserRepository(eng)
	todoRepo := repository.NewTodoRepository(eng)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	ctx := context.Background()

	// Create user and todo
	u, _ := userSvc.Create(ctx, "test@example.com", "Test User", "password123")
	userID := u["id"].(string)

	t1, _ := todoSvc.Create(ctx, userID, "Original Title", "Original Desc")
	todoID := t1["id"].(string)

	// Update todo
	err := todoSvc.Update(ctx, todoID, "Updated Title", "Updated Desc", true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update
	updated, err := todoSvc.GetByID(ctx, todoID)
	if err != nil {
		t.Fatalf("Failed to get todo: %v", err)
	}

	if title, ok := updated["title"].(string); !ok || title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %v", updated["title"])
	}

	if completed, ok := updated["completed"].(bool); !ok || !completed {
		t.Errorf("Expected completed true, got %v", updated["completed"])
	}
}

// TestTodoDelete tests deleting todo
func TestTodoDelete(t *testing.T) {
	eng := setupTestEngine(t)
	userRepo := repository.NewUserRepository(eng)
	todoRepo := repository.NewTodoRepository(eng)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	ctx := context.Background()

	// Create user and todo
	u, _ := userSvc.Create(ctx, "test@example.com", "Test User", "password123")
	userID := u["id"].(string)

	t1, _ := todoSvc.Create(ctx, userID, "Test Todo", "Description")
	todoID := t1["id"].(string)

	// Delete todo
	err := todoSvc.Delete(ctx, todoID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = todoSvc.GetByID(ctx, todoID)
	if err != todo.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestTodoToggleCompletion tests toggling todo completion
func TestTodoToggleCompletion(t *testing.T) {
	eng := setupTestEngine(t)
	userRepo := repository.NewUserRepository(eng)
	todoRepo := repository.NewTodoRepository(eng)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	ctx := context.Background()

	// Create user and todo
	u, _ := userSvc.Create(ctx, "test@example.com", "Test User", "password123")
	userID := u["id"].(string)

	t1, _ := todoSvc.Create(ctx, userID, "Test Todo", "Description")
	todoID := t1["id"].(string)

	// Toggle completion
	err := todoSvc.ToggleCompletion(ctx, todoID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify toggle
	updated, _ := todoSvc.GetByID(ctx, todoID)
	if completed, ok := updated["completed"].(bool); !ok || !completed {
		t.Errorf("Expected completed true, got %v", updated["completed"])
	}

	// Toggle again
	todoSvc.ToggleCompletion(ctx, todoID)
	updated, _ = todoSvc.GetByID(ctx, todoID)
	if completed, ok := updated["completed"].(bool); !ok || completed {
		t.Errorf("Expected completed false, got %v", updated["completed"])
	}
}
