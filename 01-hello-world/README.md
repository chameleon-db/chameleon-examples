# 01 - Hello World

The smallest ChameleonDB project: one entity, one table, a few queries from Go.

You will:

1. Declare an entity in a `.cham` file.
2. Preview and apply the migration on PostgreSQL.
3. Insert sample rows.
4. Query and insert from a small Go program.

Tested with ChameleonDB v1.1.3-alpha and PostgreSQL 16 (Linux x86-64).

## Prerequisites

- The `chameleon` CLI ([install](https://github.com/chameleon-db/chameleondb#install-linux-and-macos)):
  ```bash
  curl -sSL https://chameleondb.dev/install | sh
  chameleon version
  ```
- A running PostgreSQL server and the `psql` and `createdb` clients.
- For step 7 only: Go 1.25.6 or newer and a C compiler (the Go package links the native library that the installer put in `/usr/local/lib`).

## 1. Create the database and set the connection

```bash
createdb hello_chameleon
export DATABASE_URL="postgresql://postgres:postgres@localhost/hello_chameleon?sslmode=disable"
```

Adjust user, password and host to your setup. `.chameleon.yml` reads `DATABASE_URL`.

## 2. Validate the schema

```bash
chameleon validate
```

```
✓ Schema is valid
```

The schema is in [`schemas/user.cham`](schemas/user.cham):

```
entity User {
    id: uuid primary default uuid(),
    name: string,
    email: string unique,
    created_at: timestamp default now(),
}
```

- `uuid`, `string` and `timestamp` are the field types.
- `primary` marks the primary key. `default uuid()` generates the value in the database (`gen_random_uuid()`, PostgreSQL 13 or newer). Without a default you would have to provide the `id` yourself.
- `unique` adds a unique constraint. `default now()` fills in the creation time.

## 3. Preview the migration

```bash
chameleon migrate --dry-run
```

Among other lines, it prints the SQL it would run:

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    email VARCHAR NOT NULL UNIQUE,
    name VARCHAR NOT NULL
);
```

Columns are listed alphabetically, not in schema order. Migrations are create-only: they create missing tables and never alter or drop existing ones.

## 4. Apply it

```bash
chameleon migrate --apply
```

```
✓ Registered as v001 (hash: ...)
✓ Migration applied successfully
```

This also creates the schema vault in `.chameleon/vault/`, which records the schema version and its hash. Check it:

```bash
chameleon status      # Current version: v001, Status: Up to date
chameleon verify      # All checks passed
```

## 5. Insert sample data

```bash
psql "$DATABASE_URL" -f seed.sql
psql "$DATABASE_URL" -c "SELECT name, email FROM users ORDER BY name;"
```

You should see Alice, Bob and Charlie.

## 6. Query from the CLI

```bash
chameleon query User --debug
```

```
[SQL]
SELECT "name", "created_at", "email", "id"
FROM "users"

✓ Retrieved 3 row(s)
```

## 7. Query and insert from Go

[`main.go`](main.go) loads the schema from the vault, connects with `DATABASE_URL`, lists users, filters by email, inserts a user and reads a page of results.

```bash
go mod tidy
go run .
```

Expected output (the last section depends on how many times you ran it):

```
All users:
  Alice Smith <alice@example.com>
  Bob Johnson <bob@example.com>
  Charlie Brown <charlie@example.com>

Found by email: Alice Smith

Inserted Dana Scully <dana+...@example.com>

Newest 2 users:
  Dana Scully
  ...
```

To see the SQL behind each query:

```bash
CHAMELEON_DEBUG=sql go run .
```

## Start over

After the first `migrate --apply` the project is in `readonly` integrity mode, so running `migrate --apply` again fails with `readonly mode: schema locked`. To go back to a clean state in this example:

```bash
dropdb hello_chameleon && createdb hello_chameleon
rm -rf .chameleon
chameleon migrate --apply
```

To change a schema in a project you want to keep, see [integrity modes](https://docs.chameleondb.dev). Note that migrations do not alter existing tables yet.

## Troubleshooting

- `config file not found: .chameleon.yml`: run the commands from this folder.
- `connection refused` or authentication errors: check `DATABASE_URL` with `psql "$DATABASE_URL" -c "SELECT 1"`.
- `vault not initialized` when running the Go program: run `chameleon migrate --apply` first.
- `chameleon: error while loading shared libraries`: the library is not in the loader path. Reinstall with the default directories, or set `LD_LIBRARY_PATH` to the library directory.

## Next

- [02-blog](../02-blog): relations between entities.
- [Documentation](https://docs.chameleondb.dev)
