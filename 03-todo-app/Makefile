.PHONY: help up down wait-db migrate seed api init reset clean

DB_URL ?= postgresql://chameleon:chameleon@localhost:5432/chameleon?sslmode=disable
DB_SERVICE ?= db
DB_USER ?= chameleon
DB_NAME ?= chameleon
DB_WAIT_RETRIES ?= 40
DB_WAIT_SECONDS ?= 2

help:
	@echo "Targets (they use PostgreSQL in Docker; with your own PostgreSQL see the README):"
	@echo "  make up       - start PostgreSQL in Docker"
	@echo "  make wait-db  - wait until PostgreSQL accepts connections"
	@echo "  make migrate  - apply the schema with the chameleon CLI"
	@echo "  make seed     - load the demo user and todos"
	@echo "  make api      - run the API on localhost:8080 (foreground)"
	@echo "  make init     - up, wait-db, migrate, seed, then api"
	@echo "  make down     - stop the containers"
	@echo "  make reset    - recreate the database and the local vault, migrate and seed"
	@echo "  make clean    - remove containers, volumes and the local vault"

up:
	@echo "Starting PostgreSQL..."
	docker compose up -d db

down:
	docker compose down --remove-orphans

wait-db:
	@echo "Waiting for PostgreSQL..."
	@retries=0; \
	until docker compose exec -T $(DB_SERVICE) pg_isready -U $(DB_USER) -d $(DB_NAME) >/dev/null 2>&1; do \
		retries=$$((retries+1)); \
		if [ $$retries -ge $(DB_WAIT_RETRIES) ]; then \
			echo "PostgreSQL not ready after $$retries attempts"; \
			exit 1; \
		fi; \
		echo "   retry $$retries/$(DB_WAIT_RETRIES)..."; \
		sleep $(DB_WAIT_SECONDS); \
	done; \
	echo "PostgreSQL is ready"

migrate: wait-db
	@echo "Applying the schema..."
	DATABASE_URL="$(DB_URL)" chameleon migrate --apply

seed: wait-db
	@echo "Seeding the database..."
	DATABASE_URL="$(DB_URL)" go run ./seeds/seed.go

api:
	@echo "API on http://localhost:8080"
	DATABASE_URL="$(DB_URL)" go run ./cmd/api/main.go

init: up wait-db migrate seed
	@$(MAKE) api

# The vault records the schema version of this project. A new database needs a
# new vault, otherwise `migrate --apply` stops: the project is in readonly mode.
reset:
	docker compose down -v --remove-orphans
	rm -rf .chameleon
	$(MAKE) up
	$(MAKE) wait-db
	$(MAKE) migrate
	$(MAKE) seed

clean:
	docker compose down --volumes --remove-orphans
	rm -rf .chameleon bin/
	go clean -testcache
