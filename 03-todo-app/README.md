# todo-app

ChameleonDB project initialized with `chameleon init`.

## Quick Start

### 1. Set up database connection

```bash
export DATABASE_URL="postgresql://user:password@localhost/dbname"
```

### 2. Validate your schema

```bash
chameleon migrate --check
```

### 3. Preview migration (dry-run)

```bash
chameleon migrate --dry-run
```

### 4. Apply migration to database

```bash
chameleon migrate --apply
```

## Project Structure

```
.
├── .chameleon.yml          Configuration (version controlled)
├── .chameleon/             Admin directory (local, not versioned)
│   ├── config.yml          Source config
│   ├── state/              Current DB state
│   ├── journal/            Audit logs
│   └── backups/            Migration backups
├── schemas/                Schema files
│   └── example.cham        Example schema
└── README.md               This file
```

## Configuration

Edit `.chameleon.yml` to:
- Change database driver (postgresql, mysql, sqlite)
- Set connection string or use `${DATABASE_URL}` env var
- Configure features (auto_migration, rollback, backup, audit_logging)
- Set safety options (validation, confirmation)

## Schema

Define your database schema in `schemas/*.cham`:

```
entity User {
    id: uuid primary,
    email: string unique,
    name: string,
    created_at: timestamp default now(),
}

entity Post {
    id: uuid primary,
    title: string,
    author_id: uuid,
    created_at: timestamp default now(),
}
```

Run `chameleon migrate --check` to validate.

## Migrations

Migrations are tracked in `.chameleon/state/migrations/manifest.json`.

Each migration:
- Has a unique version (timestamp-based)
- Includes schema hash for integrity
- Supports rollback (planned v0.2)
- Is backed up before applying

View history:

```bash
chameleon journal migrations
```

## Development

### Validate schema

```bash
chameleon migrate --check
```

### View audit log

```bash
chameleon journal last 10
```

### See migration history

```bash
chameleon journal migrations
```

## Learn More

- [ChameleonDB Documentation](https://chameleondb.dev/docs)
- [Schema Reference](https://chameleondb.dev/docs/schema)
- [Query API](https://chameleondb.dev/docs/query)
