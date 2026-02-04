# Example Queries

This example demonstrates how ChameleonDB handles relations.

## Basic Queries

### Get all published posts
```go
posts, err := engine.Query("Post").
    Filter("published", "eq", true).
    OrderBy("created_at", "desc").
    Execute(ctx)
```

**Generated SQL:**
```sql
SELECT id, title, content, published, created_at, updated_at, author_id
FROM posts
WHERE published = true
ORDER BY created_at DESC;
```

---

## Relations (Eager Loading)

### Get post with author
```go
post, err := engine.Query("Post").
    Filter("id", "eq", postID).
    Include("author").
    Execute(ctx)
```

**Generated SQL:**
```sql
-- Main query
SELECT id, title, content, published, created_at, updated_at, author_id
FROM posts
WHERE id = '...';

-- Eager load author
SELECT id, username, email, bio, created_at
FROM users
WHERE id IN ('...');
```

**Result structure:**
```
post.Rows[0].String("title")          → "Getting Started..."
post.Relations["author"][0].String("username") → "alice"
```

---

### Get user with all their posts
```go
user, err := engine.Query("User").
    Filter("username", "eq", "alice").
    Include("posts").
    Execute(ctx)
```

**Generated SQL:**
```sql
-- Main query
SELECT id, username, email, bio, created_at
FROM users
WHERE username = 'alice';

-- Eager load posts
SELECT id, title, content, published, created_at, updated_at, author_id
FROM posts
WHERE author_id IN ('...');
```

---

### Get post with author AND comments
```go
post, err := engine.Query("Post").
    Filter("id", "eq", postID).
    Include("author").
    Include("comments").
    Execute(ctx)
```

**Generated SQL:**
```sql
-- Main query
SELECT ... FROM posts WHERE id = '...';

-- Load author
SELECT ... FROM users WHERE id IN ('...');

-- Load comments
SELECT ... FROM comments WHERE post_id IN ('...');
```

---

## Nested Relations

### Get post with comments AND comment authors
```go
post, err := engine.Query("Post").
    Filter("id", "eq", postID).
    Include("comments").
    Include("comments.author").
    Execute(ctx)
```

**Generated SQL:**
```sql
-- Main query
SELECT ... FROM posts WHERE id = '...';

-- Load comments
SELECT ... FROM comments WHERE post_id IN ('...');

-- Load comment authors
SELECT ... FROM users WHERE id IN ('...');  -- IDs from comments
```

**Result structure:**
```
post.Rows[0].String("title")                    → "Getting Started..."
post.Relations["comments"][0].String("content") → "Great introduction!"
post.Relations["author"][0].String("username")  → "bob"
```

---

## Complex Queries

### Get users who have commented on published posts
```go
users, err := engine.Query("User").
    Filter("comments.post.published", "eq", true).
    Execute(ctx)
```

**Generated SQL:**
```sql
SELECT DISTINCT users.id, users.username, users.email, users.bio, users.created_at
FROM users
INNER JOIN comments ON comments.author_id = users.id
INNER JOIN posts ON posts.id = comments.post_id
WHERE posts.published = true;
```

---

### Get all posts by users with more than 2 posts
```go
// This shows the power of graph navigation
// In traditional SQL you'd need a subquery or GROUP BY + HAVING

users, err := engine.Query("User").
    Include("posts").
    Execute(ctx)

// Filter in application (for now)
for _, user := range users.Rows {
    posts := users.Relations["posts"]
    if len(posts) > 2 {
        // This user has more than 2 posts
    }
}
```

**Note:** Aggregations in queries coming in future versions.

---

## Pagination

### Get latest 5 posts with authors
```go
posts, err := engine.Query("Post").
    Filter("published", "eq", true).
    Include("author").
    OrderBy("created_at", "desc").
    Limit(5).
    Execute(ctx)
```

---

## Key Concepts

### HasMany
A `User` has many `Post`s:
```rust
posts: [Post] via author_id,
```

The `via` specifies which field in `Post` points back to `User`.

### BelongsTo
A `Post` belongs to a `User`:
```rust
author_id: uuid,
author: User,
```

The foreign key field (`author_id`) and the relation (`author`) work together.

### Eager Loading
Without `Include()`:
- You get the main entity only
- Relations are not loaded

With `Include()`:
- Relations are loaded in separate queries
- No N+1 problem
- Results are merged automatically

### Nested Includes
```go
Include("comments").Include("comments.author")
```

Loads comments, then loads the author of each comment.
