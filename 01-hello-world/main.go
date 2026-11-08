// Hello World: query and insert with ChameleonDB.
//
// Run `chameleon migrate --apply` first: the engine loads the schema from the
// vault that the first migration creates.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chameleon-db/chameleondb/chameleon/pkg/engine"
	_ "github.com/chameleon-db/chameleondb/chameleon/pkg/engine/mutation"
)

func main() {
	ctx := context.Background()

	eng, err := engine.NewEngine()
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()

	cfg, err := engine.ParseConnectionString(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DATABASE_URL: %v", err)
	}
	if err := eng.Connect(ctx, cfg); err != nil {
		log.Fatal(err)
	}

	// 1. All users, ordered by name.
	users, err := eng.Query("User").OrderBy("name", "asc").Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All users:")
	for _, u := range users.Rows {
		fmt.Printf("  %s <%s>\n", u.String("name"), u.String("email"))
	}

	// 2. One user, by email.
	found, err := eng.Query("User").Filter("email", "eq", "alice@example.com").Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !found.IsEmpty() {
		fmt.Printf("\nFound by email: %s\n", found.Rows[0].String("name"))
	}

	// 3. Insert a user. `id` and `created_at` come from the schema defaults.
	email := fmt.Sprintf("dana+%d@example.com", time.Now().Unix())
	if _, err := eng.Insert("User").Set("name", "Dana Scully").Set("email", email).Execute(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInserted Dana Scully <%s>\n", email)

	// 4. First page of 2 users, newest first.
	page, err := eng.Query("User").OrderBy("created_at", "desc").Limit(2).Offset(0).Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nNewest 2 users:")
	for _, u := range page.Rows {
		fmt.Printf("  %s\n", u.String("name"))
	}
}
