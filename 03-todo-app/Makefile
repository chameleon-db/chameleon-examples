.PHONY: help setup build up down wait-db migrate seed api init reset clean

DB_URL ?= postgresql://chameleon:chameleon@localhost:5432/chameleon?sslmode=disable
DB_SERVICE ?= db
DB_USER ?= chameleon
DB_NAME ?= chameleon
DB_WAIT_RETRIES ?= 40
DB_WAIT_SECONDS ?= 2

help:
	@echo "Targets disponibles:"
	@echo "  make setup    - descarga artefactos necesarios"
	@echo "  make build    - construye imágenes Docker"
	@echo "  make up       - levanta contenedores en background"
	@echo "  make wait-db  - espera hasta que Postgres acepte conexiones"
	@echo "  make migrate  - aplica migraciones con ChameleonDB"
	@echo "  make seed     - siembra datos demo"
	@echo "  make api      - ejecuta la API en local (foreground)"
	@echo "  make init     - inicializa todo y deja la API funcionando"
	@echo "  make down     - baja contenedores"
	@echo "  make reset    - reinicia DB y vuelve a migrar/seed"
	@echo "  make clean    - limpia TODO (contenedores, volúmenes, imágenes, artefactos)"

setup:
	@echo "🔧 Preparing artifacts..."
	./setup.sh

build: setup
	@echo "🐳 Building images..."
	docker compose build

up:
	@echo "🚀 Starting stack in background..."
	@docker compose stop api >/dev/null 2>&1 || true
	docker compose up -d db

down:
	@echo "🛑 Stopping stack..."
	docker compose down --remove-orphans

wait-db:
	@echo "⏳ Waiting for PostgreSQL to be ready..."
	@retries=0; \
	until docker compose exec -T $(DB_SERVICE) pg_isready -U $(DB_USER) -d $(DB_NAME) >/dev/null 2>&1; do \
		retries=$$((retries+1)); \
		if [ $$retries -ge $(DB_WAIT_RETRIES) ]; then \
			echo "❌ PostgreSQL not ready after $$retries attempts"; \
			exit 1; \
		fi; \
		echo "   retry $$retries/$(DB_WAIT_RETRIES)..."; \
		sleep $(DB_WAIT_SECONDS); \
	done; \
	echo "✅ PostgreSQL is ready"

migrate: wait-db
	@echo "🧬 Applying migrations..."
	DATABASE_URL="$(DB_URL)" chameleon migrate --apply

seed: wait-db
	@echo "🌱 Seeding database..."
	DATABASE_URL="$(DB_URL)" go run ./seeds/seed.go

api:
	@echo "🌐 Starting API on localhost:8080..."
	go run ./cmd/api/main.go

init: setup build up wait-db migrate seed
	@echo "✅ Infra lista. Iniciando API..."
	@$(MAKE) api

reset:
	docker compose down -v --remove-orphans
	$(MAKE) up
	$(MAKE) wait-db
	$(MAKE) migrate
	$(MAKE) seed

clean:
	@echo "🧹 Cleaning Docker + generated artifacts..."
	docker compose down --volumes --rmi all --remove-orphans
	rm -rf .artifacts/
	rm -rf bin/
	rm -rf .chameleon/state/ .chameleon/journal/ .chameleon/backups/
	go clean -testcache
