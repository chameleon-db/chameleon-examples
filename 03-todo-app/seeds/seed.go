package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()

	fmt.Println("🌱 Seeding database...")
	fmt.Println()

	// Get connection string from env or use default
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgresql://chameleon:chameleon@localhost:5432/chameleon?sslmode=disable"
	}

	sqlPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("❌ Failed to create SQL pool: %v", err)
	}
	defer sqlPool.Close()

	if err := sqlPool.Ping(ctx); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Connected to database")
	fmt.Println()

	// ============================================
	// 2. Seed Demo User
	// ============================================

	const demoEmail = "test@example.com"
	const demoPassword = "password123"

	var demoUserID string

	err = sqlPool.QueryRow(
		ctx,
		`SELECT id FROM users WHERE email = $1 LIMIT 1`,
		demoEmail,
	).Scan(&demoUserID)

	if err != nil && err.Error() != "no rows in result set" {
		log.Fatalf("❌ Failed to query demo user: %v", err)
	}

	if demoUserID != "" {
		fmt.Printf("⏭️  Demo user already exists: %s (id: %s)\n", demoEmail, demoUserID)
	} else {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("❌ Failed to hash demo password: %v", err)
		}

		demoUserID = uuid.New().String()
		_, err = sqlPool.Exec(
			ctx,
			`INSERT INTO users (id, email, name, password_hash, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
			demoUserID,
			demoEmail,
			"Demo User",
			string(passwordHash),
			true,
		)
		if err != nil {
			log.Fatalf("❌ Failed to create demo user: %v", err)
		}

		fmt.Printf("✅ Created demo user: %s (id: %s)\n", demoEmail, demoUserID)
	}

	fmt.Println()

	// ============================================
	// 3. Seed Todos
	// ============================================

	todos := []struct {
		title       string
		description string
		completed   bool
	}{
		{
			title:       "Learn ChameleonDB",
			description: "Understand Schema Vault, Integrity Modes, and query building",
			completed:   false,
		},
		{
			title:       "Test basic CRUD",
			description: "Create, read, update, and delete todos using the API",
			completed:   false,
		},
		{
			title:       "Break GORM forever",
			description: "Never look back at magic ORM behavior again",
			completed:   true, // Already done! 🎉
		},
		{
			title:       "Drink mate peacefully",
			description: "Enjoy the serenity of type-safe queries",
			completed:   false,
		},
		{
			title:       "Explore Debug mode",
			description: "Use .Debug() to see generated SQL",
			completed:   false,
		},
		{
			title:       "Test Schema Vault",
			description: "Run 'chameleon journal schema' to see version history",
			completed:   false,
		},
	}

	fmt.Println("📝 Creating todos...")

	successCount := 0
	skipCount := 0

	for _, todo := range todos {
		var existingID string
		err := sqlPool.QueryRow(
			ctx,
			`SELECT id FROM todos WHERE user_id = $1 AND title = $2 LIMIT 1`,
			demoUserID,
			todo.title,
		).Scan(&existingID)

		if err != nil && err.Error() != "no rows in result set" {
			log.Printf("⚠️  Warning: Failed to check existing todo '%s': %v", todo.title, err)
			continue
		}

		if existingID != "" {
			fmt.Printf("⏭️  Skipped: '%s' (already exists)\n", todo.title)
			skipCount++
			continue
		}

		// Insert new todo
		todoID := uuid.New().String()
		_, err = sqlPool.Exec(
			ctx,
			`INSERT INTO todos (id, title, description, completed, user_id, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
			todoID,
			todo.title,
			todo.description,
			todo.completed,
			demoUserID,
		)

		if err != nil {
			log.Printf("❌ Failed to create todo '%s': %v", todo.title, err)
			continue
		}

		fmt.Printf("✅ Created: '%s' (id: %s)\n", todo.title, todoID)
		successCount++
	}

	fmt.Println()
	fmt.Printf("📊 Summary: %d created, %d skipped, %d total\n",
		successCount, skipCount, len(todos))
	fmt.Println()

	// ============================================
	// 4. Verify seeded data
	// ============================================

	fmt.Println("🔍 Verifying seeded data...")

	var total int
	if err := sqlPool.QueryRow(ctx, `SELECT COUNT(*) FROM todos WHERE user_id = $1`, demoUserID).Scan(&total); err != nil {
		log.Fatalf("❌ Failed to count total todos: %v", err)
	}

	var completed int
	if err := sqlPool.QueryRow(ctx, `SELECT COUNT(*) FROM todos WHERE user_id = $1 AND completed = true`, demoUserID).Scan(&completed); err != nil {
		log.Fatalf("❌ Failed to count completed todos: %v", err)
	}

	pending := total - completed

	fmt.Printf("   Total todos in database: %d\n", total)
	fmt.Printf("   Completed todos: %d\n", completed)
	fmt.Printf("   Pending todos: %d\n", pending)
	fmt.Println()

	fmt.Println()
	fmt.Println("✅ Seeding complete!")
	fmt.Println()
	fmt.Printf("🔐 Demo login: %s / %s\n", demoEmail, demoPassword)
	fmt.Println()
	fmt.Println("💡 Next steps:")
	fmt.Println("   • Start the API:  go run ./cmd/api")
	fmt.Println("   • Test endpoint:  curl http://localhost:8080/todos")
	fmt.Println("   • View vault:     chameleon journal schema")
	fmt.Println()
}
