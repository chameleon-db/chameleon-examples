# Example Queries

This file shows common queries you can run with ChameleonDB.

## Fetch all users
```go
users, err := engine.Query("User").Execute(ctx)
if err != nil {
    log.Fatal(err)
}

for _, user := range users.Rows {
    fmt.Printf("User: %s (%s)\n", user.String("name"), user.String("email"))
}
```

**Generated SQL:**
```sql
SELECT id, name, email, created_at
FROM users;
```

---

## Find user by email
```go
user, err := engine.Query("User").
    Filter("email", "eq", "alice@example.com").
    Execute(ctx)
```

**Generated SQL:**
```sql
SELECT id, name, email, created_at
FROM users
WHERE email = 'alice@example.com';
```

---

## Order by name
```go
users, err := engine.Query("User").
    OrderBy("name", "asc").
    Execute(ctx)
```

**Generated SQL:**
```sql
SELECT id, name, email, created_at
FROM users
ORDER BY name ASC;
```

---

## Pagination
```go
users, err := engine.Query("User").
    OrderBy("created_at", "desc").
    Limit(10).
    Offset(0).
    Execute(ctx)
```

**Generated SQL:**
```sql
SELECT id, name, email, created_at
FROM users
ORDER BY created_at DESC
LIMIT 10 OFFSET 0;
```
