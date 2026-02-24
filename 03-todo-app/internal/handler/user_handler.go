package handler

import (
	"net/http"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/user"
	"github.com/go-chi/chi/v5"
)

// UserHandler handles user HTTP endpoints
type UserHandler struct {
	service user.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler(svc user.Service) *UserHandler {
	return &UserHandler{service: svc}
}

// CreateRequest is the request body for create user
type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// POST /users - Create user
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	u, err := h.service.Create(r.Context(), req.Email, req.Name, req.Password)
	if err != nil {
		switch err {
		case user.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid email, name, or password")
		case user.ErrWeakPassword:
			respondError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		case user.ErrDuplicateEmail:
			respondError(w, http.StatusConflict, "Email already exists")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to create user")
		}
		return
	}

	respondJSON(w, http.StatusCreated, u)
}

// GET /users/{id} - Get user by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case user.ErrNotFound:
			respondError(w, http.StatusNotFound, "User not found")
		case user.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid user ID")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	respondJSON(w, http.StatusOK, u)
}

// GET /users - List users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := queryIntParam(r, "limit", 10)
	offset := queryIntParam(r, "offset", 0)

	// Clamp limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondJSON(w, http.StatusOK, users)
}

// UpdateUserRequest is the request body for update user
type UpdateUserRequest struct {
	Name string `json:"name"`
}

// PUT /users/{id} - Update user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.service.Update(r.Context(), id, req.Name)
	if err != nil {
		switch err {
		case user.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid user ID or name")
		case user.ErrNotFound:
			respondError(w, http.StatusNotFound, "User not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// DELETE /users/{id} - Delete user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		switch err {
		case user.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid user ID")
		case user.ErrNotFound:
			respondError(w, http.StatusNotFound, "User not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /login - Login user
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	u, err := h.service.VerifyPassword(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case user.ErrInvalidPassword:
			respondError(w, http.StatusUnauthorized, "Invalid email or password")
		case user.ErrUserInactive:
			respondError(w, http.StatusForbidden, "User is inactive")
		case user.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "Invalid input")
		default:
			respondError(w, http.StatusInternalServerError, "Login failed")
		}
		return
	}

	respondJSON(w, http.StatusOK, u)
}
