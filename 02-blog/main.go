// Blog: relations and eager loading with ChameleonDB.
//
// Run `chameleon migrate --apply` and load seed.sql first.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chameleon-db/chameleondb/chameleon/pkg/engine"
	_ "github.com/chameleon-db/chameleondb/chameleon/pkg/engine/mutation"
)

// childrenOf returns the rows whose foreign key points at parentID.
// Eager-loaded relations come back as flat lists, so the program groups them.
func childrenOf(rows []engine.Row, fk, parentID string) []engine.Row {
	var out []engine.Row
	for _, r := range rows {
		if r.String(fk) == parentID {
			out = append(out, r)
		}
	}
	return out
}

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

	// 1. Published posts, newest first.
	posts, err := eng.Query("Post").
		Filter("published", "eq", true).
		OrderBy("created_at", "desc").
		Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Published posts: %d\n", posts.Count())
	for _, p := range posts.Rows {
		fmt.Printf("  %s\n", p.String("title"))
	}

	// 2. One user with their posts, and each post with its comments.
	//    Include("posts") and Include("posts.comments") run one extra query each.
	res, err := eng.Query("User").
		Filter("username", "eq", "alice").
		Include("posts").
		Include("posts.comments").
		Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.IsEmpty() {
		log.Fatal("alice not found: did you load seed.sql?")
	}
	alice := res.Rows[0]
	fmt.Printf("\n%s and her posts:\n", alice.String("username"))
	for _, p := range childrenOf(res.Relations["posts"], "author_id", alice.String("id")) {
		cs := childrenOf(res.Relations["posts.comments"], "post_id", p.String("id"))
		fmt.Printf("  %s (%d comments)\n", p.String("title"), len(cs))
	}

	// 3. Filter through a HasMany relation: users with at least one published post.
	authors, err := eng.Query("User").
		Filter("posts.published", "eq", true).
		OrderBy("username", "asc").
		Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nUsers with a published post:")
	for _, u := range authors.Rows {
		fmt.Printf("  %s\n", u.String("username"))
	}

	// 4. Latest 2 users, each with their posts.
	latest, err := eng.Query("User").
		Include("posts").
		OrderBy("created_at", "desc").
		Limit(2).
		Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nLatest 2 users:")
	for _, u := range latest.Rows {
		n := len(childrenOf(latest.Relations["posts"], "author_id", u.String("id")))
		fmt.Printf("  %s (%d posts)\n", u.String("username"), n)
	}
}
