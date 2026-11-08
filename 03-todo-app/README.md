# 03 - Todo app

A small HTTP API for users and their todos, built on ChameleonDB. It shows a ChameleonDB engine inside a layered Go application, with a demo front end.

```
handler (HTTP)  ->  domain service (rules)  ->  repository (ChameleonDB queries and mutations)
```

- `schemas/`: the `User` and `Todo` entities.
- `internal/repository/`: all database access, through `engine.Query`, `Insert`, `Update` and `Delete`.
- `internal/domain/`: validation and business rules (password length, bcrypt, soft delete).
- `internal/handler/` and `internal/router/`: the REST API (chi).
- `seeds/seed.go`: a demo user and todos.
- `../03-todo-front`: a plain HTML/JS front end for the API.

## Prerequisites

- The `chameleon` CLI ([install](https://github.com/chameleon-db/chameleondb#install-linux-and-macos)) and `chameleon version` working.
- Go 1.25.6 or newer and a C compiler (the Go package links the native library the installer put in `/usr/local/lib`).
- PostgreSQL: Docker (for `make up`) or your own server.

## Quick start with Docker

```bash
make init
```

This starts PostgreSQL in Docker, applies the schema, loads the demo data and runs the API on `http://localhost:8080`. The same steps one by one:

```bash
make up        # PostgreSQL in Docker
make migrate   # chameleon migrate --apply
make seed      # demo user and todos
make api       # the API, in the foreground
```

## Quick start with your own PostgreSQL

```bash
createdb todo_chameleon
export DATABASE_URL="postgresql://postgres:postgres@localhost/todo_chameleon?sslmode=disable"

chameleon validate
chameleon migrate --dry-run
chameleon migrate --apply      # also creates the schema vault in .chameleon/vault/
go run ./seeds/seed.go
go run ./cmd/api/main.go
```

Adjust user, password and host to your setup. The API reads `DATABASE_URL` and `PORT` (default 8080).

## The schema

```
entity User {
    id: uuid primary default uuid(),
    email: string unique,
    name: string,
    password_hash: string,
    created_at: timestamp default now(),
    updated_at: timestamp default now(),
    is_active: bool default true,
}

entity Todo {
    id: uuid primary default uuid(),
    title: string,
    description: string nullable,
    completed: bool default false,
    user_id: uuid,
    created_at: timestamp default now(),
    updated_at: timestamp default now(),
    due_date: timestamp nullable,
}
```

Inserts only need to send the fields that have no default. A field that is not nullable and has no default is required: the engine rejects an insert that leaves it out. That is why `id`, `is_active` and `completed` have defaults here.

## API

Responses are `{"data": ...}` on success and `{"error": "..."}` on failure.

| Method and path | Body | What it does |
|---|---|---|
| `GET /health` | | Liveness check |
| `POST /users` | `email`, `name`, `password` (8 characters or more) | Create a user |
| `GET /users?limit=10&offset=0` | | List active users, newest first |
| `GET /users/{id}` | | Get a user |
| `PUT /users/{id}` | `name` | Rename a user |
| `DELETE /users/{id}` | | Deactivate a user (`is_active = false`) |
| `POST /login` | `email`, `password` | Check credentials and return the user. There is no session or token. |
| `POST /users/{userID}/todos` | `title`, `description` | Create a todo |
| `GET /users/{userID}/todos?completed=false&limit=10&offset=0` | | List a user's todos, newest first |
| `GET /users/{userID}/todos/overdue` | | Open todos whose `due_date` has passed |
| `GET /users/{userID}/todos/{id}` | | Get a todo |
| `PUT /users/{userID}/todos/{id}` | `title`, `description`, `completed` | Update a todo |
| `PATCH /users/{userID}/todos/{id}/toggle` | | Flip `completed` |
| `DELETE /users/{userID}/todos/{id}` | | Delete a todo |

Try it with the demo user (`test@example.com` / `password123`):

```bash
curl -s localhost:8080/health

curl -s -X POST localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"password123"}'

# use the "id" of the user from the response
curl -s localhost:8080/users/<id>/todos
curl -s -X POST localhost:8080/users/<id>/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"Read the docs","description":"docs.chameleondb.dev"}'
```

This example is a demo of ChameleonDB, not a model of authentication: there are no tokens and the routes are open.

## Front end

`../03-todo-front` is a static page that talks to this API on `http://localhost:8080`. With the API running:

```bash
cd ../03-todo-front
python3 -m http.server 3000
# open http://localhost:3000 and log in with the demo user
```

## Testing

The integration tests in `test/integration` start from empty `users` and `todos` tables, so they need a database of their own. Make one as a copy of the development database (stop the API first, because `createdb -T` needs no open connections to the template):

```bash
createdb -T todo_chameleon todo_chameleon_test
export TEST_DATABASE_URL="postgresql://postgres:postgres@localhost/todo_chameleon_test?sslmode=disable"
go test ./test/integration/...
```

Without `TEST_DATABASE_URL` the tests are skipped. They run from the project root's vault, so run `chameleon migrate --apply` at least once first.

## Start over

`migrate --apply` leaves the project in `readonly` integrity mode, so a second `migrate --apply` is refused. To start from scratch (with Docker):

```bash
make reset
```

or by hand: drop and recreate the database, then `rm -rf .chameleon` and migrate again.

To change the schema of a project you want to keep, set a mode password and switch to `standard` (`chameleon config auth set-password`, then `chameleon config set mode=standard`). Migrations are create-only: they create missing tables and never alter existing ones, so a new field in an existing table must be added by hand with `ALTER TABLE`.

## Limits worth knowing

- Filters do not accept `time.Time` values yet, so `GetOverdue` loads the open todos and compares `due_date` in Go.
- Updates do not touch `updated_at`; the column only has its default at insert time.
- `Todo.user_id` is a plain `uuid` field, not a relation, so the database has no foreign key to `users`.
