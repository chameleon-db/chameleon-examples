package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/chameleon-db/chameleon-examples/todo-app/internal/config"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/todo"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/domain/user"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/handler"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/repository"
	"github.com/chameleon-db/chameleon-examples/todo-app/internal/router"

	"github.com/chameleon-db/chameleondb/chameleon/pkg/engine"
)

func main() {
	// Load configuration
	cfg := config.Load()

	log.Printf("Starting Todo App on port %d", cfg.Port)

	// Initialize database engine
	log.Println("Initializing database engine...")
	eng, err := engine.NewEngine()
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()

	// Connect to database
	ctx := context.Background()
	dbConfig, err := engine.ParseConnectionString(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	if err := eng.Connect(ctx, dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer eng.Close()
	log.Println("Database connected successfully")

	// Initialize repositories
	log.Println("Initializing repositories...")
	userRepo := repository.NewUserRepository(eng)
	todoRepo := repository.NewTodoRepository(eng)

	// Initialize domain services
	log.Println("Initializing domain services...")
	var userService user.Service = user.NewService(userRepo)
	var todoService todo.Service = todo.NewService(todoRepo)

	// Initialize handlers
	log.Println("Initializing handlers...")
	userHandler := handler.NewUserHandler(userService)
	todoHandler := handler.NewTodoHandler(todoService)

	// Create router
	log.Println("Creating router...")
	r := router.New(userHandler, todoHandler)

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
