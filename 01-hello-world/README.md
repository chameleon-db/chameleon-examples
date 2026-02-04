# Hello World

The simplest possible ChameleonDB example. One entity, basic queries.

## What you'll learn

- ✓ Define an entity with fields and constraints
- ✓ Generate a migration
- ✓ Apply migration to PostgreSQL
- ✓ Insert sample data
- ✓ Run basic queries

## Prerequisites

- ChameleonDB CLI installed
- PostgreSQL running locally
- 5 minutes

## Quick Start

### 1. Create the database
```bash
# Create database (only needed once)
createdb hello_chameleon
```

### 2. Validate the schema
```bash
chameleon validate
```

You should see:
```
✓ Schema is valid
```

### 3. Generate migration
```bash
chameleon migrate
```

This shows the SQL that will be executed:
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 4. Apply migration
```bash
chameleon migrate --apply
```

This creates the `users` table in PostgreSQL.

### 5. Insert sample data
```bash
psql -U postgres hello_chameleon < seed.sql
```

### 6. Verify
```bash
psql -U postgres hello_chameleon -c "SELECT * FROM users;"
```

You should see 3 users (Alice, Bob, Charlie).

## Next Steps

- See [queries.md](./queries.md) for example queries
- Try the [blog example](../02-blog) to learn about relations
- Read the [ChameleonDB documentation](https://github.com/chameleon-db/chameleondb)

## Schema Breakdown
```rust
entity User {
    id: uuid primary,           // Primary key, auto-generated
    name: string,               // Required field
    email: string unique,       // Must be unique across all users
    created_at: timestamp default now(),  // Auto-set to current time
}
```

### Field Types

- `uuid` — Universally unique identifier
- `string` — Variable-length text (VARCHAR in PostgreSQL)
- `timestamp` — Date and time

### Constraints

- `primary` — Primary key (unique, not null, indexed)
- `unique` — Must be unique across all rows
- `default now()` — Automatically set to current timestamp

## Troubleshooting

### Database doesn't exist
```bash
createdb hello_chameleon
```

### Permission denied

Check your PostgreSQL user has CREATE permissions:
```bash
psql -U postgres -c "CREATE DATABASE hello_chameleon;"
```

### Migration already applied

Drop and recreate:
```bash
dropdb hello_chameleon
createdb hello_chameleon
chameleon migrate --apply
```
