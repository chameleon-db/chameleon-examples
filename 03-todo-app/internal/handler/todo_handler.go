package handler

import (
	"net/http"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/todo"
	"github.com/go-chi/chi/v5"
)

// TodoHandler handles todo HTTP endpoints
type TodoHandler struct {
	service todo.Service
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(svc todo.Service) *TodoHandler {
	return &TodoHandler{service: svc}
}

// CreateTodoRequest is the request body for create todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// POST /users/{userID}/todos - Create todo
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var req CreateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	t, err := h.service.Create(r.Context(), userID, req.Title, req.Description)
	if err != nil {
		switch err {
		case todo.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid user ID or title")
		case todo.ErrInvalidUserID:
			respondError(w, http.StatusBadRequest, "Invalid user ID")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to create todo")
		}
		return
	}

	respondJSON(w, http.StatusCreated, t)
}

// GET /todos/{id} - Get todo by ID
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	t, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case todo.ErrNotFound:
			respondError(w, http.StatusNotFound, "Todo not found")
		case todo.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid todo ID")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to fetch todo")
		}
		return
	}

	respondJSON(w, http.StatusOK, t)
}

// GET /users/{userID}/todos - List user's todos
func (h *TodoHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	limit := queryIntParam(r, "limit", 10)
	offset := queryIntParam(r, "offset", 0)
	completed := queryBoolParam(r, "completed")

	// Clamp limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var todos []map[string]interface{}
	var err error

	if completed != nil {
		todos, err = h.service.ListByUserFiltered(r.Context(), userID, completed, limit, offset)
	} else {
		todos, err = h.service.ListByUser(r.Context(), userID, limit, offset)
	}

	if err != nil {
		switch err {
		case todo.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid user ID")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to fetch todos")
		}
		return
	}

	respondJSON(w, http.StatusOK, todos)
}

// UpdateTodoRequest is the request body for update todo
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// PUT /todos/{id} - Update todo
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.service.Update(r.Context(), id, req.Title, req.Description, req.Completed)
	if err != nil {
		switch err {
		case todo.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid todo ID or title")
		case todo.ErrNotFound:
			respondError(w, http.StatusNotFound, "Todo not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to update todo")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Todo updated successfully"})
}

// DELETE /todos/{id} - Delete todo
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		switch err {
		case todo.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid todo ID")
		case todo.ErrNotFound:
			respondError(w, http.StatusNotFound, "Todo not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to delete todo")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Todo deleted successfully"})
}

// GET /users/{userID}/todos/overdue - Get overdue todos
func (h *TodoHandler) GetOverdue(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	todos, err := h.service.GetOverdue(r.Context(), userID)
	if err != nil {
		switch err {
		case todo.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid user ID")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to fetch overdue todos")
		}
		return
	}

	respondJSON(w, http.StatusOK, todos)
}

// ToggleCompletionRequest is the request body for toggle completion
type ToggleCompletionRequest struct {
	Completed bool `json:"completed"`
}

// PATCH /todos/{id}/toggle - Toggle todo completion
func (h *TodoHandler) ToggleCompletion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.ToggleCompletion(r.Context(), id)
	if err != nil {
		switch err {
		case todo.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid todo ID")
		case todo.ErrNotFound:
			respondError(w, http.StatusNotFound, "Todo not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to toggle todo completion")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Todo toggled successfully"})
}
