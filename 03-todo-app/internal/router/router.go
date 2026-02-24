package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/handler"
	appMiddleware "github.com/chameleon-db/chameleon-examples/todo-app/internal/middleware"
)

// New creates and configures the HTTP router
func New(userHandler *handler.UserHandler, todoHandler *handler.TodoHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(appMiddleware.CORS)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(appMiddleware.RequestLogger)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)       // POST /users
		r.Get("/", userHandler.List)          // GET /users
		r.Get("/{id}", userHandler.GetByID)   // GET /users/{id}
		r.Put("/{id}", userHandler.Update)    // PUT /users/{id}
		r.Delete("/{id}", userHandler.Delete) // DELETE /users/{id}
	})

	// Auth route
	r.Post("/login", userHandler.Login) // POST /login

	// Todo routes
	r.Route("/users/{userID}/todos", func(r chi.Router) {
		r.Post("/", todoHandler.Create)                       // POST /users/{userID}/todos
		r.Get("/", todoHandler.ListByUser)                    // GET /users/{userID}/todos
		r.Get("/overdue", todoHandler.GetOverdue)             // GET /users/{userID}/todos/overdue
		r.Get("/{id}", todoHandler.GetByID)                   // GET /users/{userID}/todos/{id}
		r.Put("/{id}", todoHandler.Update)                    // PUT /users/{userID}/todos/{id}
		r.Delete("/{id}", todoHandler.Delete)                 // DELETE /users/{userID}/todos/{id}
		r.Patch("/{id}/toggle", todoHandler.ToggleCompletion) // PATCH /users/{userID}/todos/{id}/toggle
	})

	// Global todo routes (without userID in path)
	r.Get("/todos/{id}", todoHandler.GetByID)   // GET /todos/{id}
	r.Put("/todos/{id}", todoHandler.Update)    // PUT /todos/{id}
	r.Delete("/todos/{id}", todoHandler.Delete) // DELETE /todos/{id}

	return r
}
