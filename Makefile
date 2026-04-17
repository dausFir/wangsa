.PHONY: help dev-backend dev-frontend install build clean vet fmt tidy db-create db-reset

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Download and tidy Go dependencies
	cd backend && go mod tidy && go mod download

install: tidy ## Install all dependencies (Go + npm)
	cd frontend && npm install

dev-backend: ## Run backend development server
	cd backend && go run ./cmd/server/main.go

dev-frontend: ## Run frontend development server
	cd frontend && npm run dev

build-backend: ## Build backend binary
	cd backend && go build -ldflags="-s -w" -o wangsa-server ./cmd/server/main.go

build-frontend: ## Build frontend for production
	cd frontend && npm run build

build: build-backend build-frontend ## Build everything

db-create: ## Create local PostgreSQL database for development
	createdb wangsa 2>/dev/null || echo "Database already exists"
	psql wangsa -c "CREATE USER wangsa WITH PASSWORD 'wangsa';" 2>/dev/null || true
	psql wangsa -c "GRANT ALL PRIVILEGES ON DATABASE wangsa TO wangsa;" 2>/dev/null || true
	@echo "✅ Database ready — run make dev-backend to apply schema"

db-reset: ## Drop and recreate local development database
	dropdb --if-exists wangsa
	$(MAKE) db-create

seed-admin: ## Create default superadmin account
	cd backend && go run ./cmd/seeder/main.go -create-admin

db-setup: db-create seed-admin ## Complete database setup (create + seed admin)

clean: ## Remove build artifacts
	rm -f backend/wangsa-server
	rm -rf frontend/dist
	@echo "✅ Cleaned"

vet: ## Run Go vet
	cd backend && go vet ./...

fmt: ## Format Go code
	cd backend && gofmt -w .
