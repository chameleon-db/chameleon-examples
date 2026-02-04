# Blog Example

A blog platform with users, posts, and comments.

**Complexity:** Intermediate (15 minutes)

## What you'll learn

- ✓ HasMany relations (User has many Posts)
- ✓ BelongsTo relations (Post belongs to User)
- ✓ Eager loading with `Include()`
- ✓ Nested eager loading (posts with comments with authors)
- ✓ Filtering on related entities
- ✓ Real-world data patterns

## Prerequisites

- Completed [01-hello-world](../01-hello-world)
- PostgreSQL running
- 15 minutes

## Quick Start

### 1. Create database
```bash
psql -U postgres -c "CREATE DATABASE blog_chameleon;"
```

### 2. Validate schema
```bash
chameleon validate
```

### 3. Apply migration
```bash
chameleon migrate --apply
```

### 4. Seed data
```bash
psql -U postgres -d blog_chameleon < seed.sql
```

You should see:
```
 username | post_count | comment_count 
----------+------------+---------------
 alice    |          3 |             2
 bob      |          1 |             2
 charlie  |          0 |             1
```

### 5. Explore the data
```bash
# All published posts
psql -U postgres -d blog_chameleon -c "SELECT title, username FROM posts p JOIN users u ON p.author_id = u.id WHERE published = true;"

# Posts with comment counts
psql -U postgres -d blog_chameleon -c "
SELECT 
    p.title, 
    u.username as author,
    COUNT(c.id) as comments
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN comments c ON c.post_id = p.id
GROUP BY p.id, p.title, u.username
ORDER BY comments DESC;
"
```

## Schema Overview
```
User (3 rows)
  ├─ posts: [Post]          → HasMany
  └─ comments: [Comment]    → HasMany

Post (4 rows)
  ├─ author: User           → BelongsTo
  └─ comments: [Comment]    → HasMany

Comment (5 rows)
  ├─ author: User           → BelongsTo
  └─ post: Post             → BelongsTo
```

## Key Patterns

### 1. Bidirectional Relations

User ↔ Post is defined on both sides:

**User side:**
```rust
posts: [Post] via author_id,
```

**Post side:**
```rust
author_id: uuid,
author: User,
```

This allows navigation in both directions:
- "Give me all posts by Alice"
- "Who wrote this post?"

### 2. Eager Loading Avoids N+1

**Bad (N+1 queries):**
```go
// Get all posts
posts := db.Query("Post").Execute(ctx)

// For each post, get author (N queries!)
for _, post := range posts.Rows {
    author := db.Query("User").Filter("id", "eq", post["author_id"]).Execute(ctx)
}
```

**Good (2 queries total):**
```go
// Get posts with authors in one go
posts := db.Query("Post").Include("author").Execute(ctx)
// Main query + 1 eager query = 2 total
```

### 3. Nested Includes

Load entire graph in one call:
```go
user := db.Query("User").
    Include("posts").
    Include("posts.comments").
    Include("posts.comments.author").
    Execute(ctx)
```

This loads:
1. User (main query)
2. Their posts
3. Comments on those posts
4. Authors of those comments

All with optimal SQL (4 queries instead of potentially hundreds).

## Common Queries

See [queries.md](./queries.md) for detailed examples with generated SQL.

## Next Steps

- Try [03-ecommerce](../03-ecommerce) for more complex domains
- Read about [multi-tenancy](../04-saas-multitenant)
- Explore [social networks](../05-social-network) for many-to-many relations

## Troubleshooting

### "Relation X references unknown entity Y"

Make sure the target entity exists in the schema and is spelled correctly.

### "Missing foreign key"

HasMany relations need `via`:
```rust
posts: [Post] via author_id,  // ← author_id must exist in Post
```

### Migration fails

Drop and recreate:
```bash
psql -U postgres -c "DROP DATABASE blog_chameleon;"
psql -U postgres -c "CREATE DATABASE blog_chameleon;"
chameleon migrate --apply
```
