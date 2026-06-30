.PHONY: install run-backend run-frontend run test wire build

install:
	@echo "Installing backend dependencies..."
	cd apps/core && go mod download
	@echo "Installing frontend dependencies..."
	cd apps/portal && npm install

run-backend:
	@echo "Starting Go backend with Air hot reload..."
	cd apps/core && go run github.com/air-verse/air

run-frontend:
	@echo "Starting React frontend with Vite..."
	cd apps/portal && npm run dev

run:
	npx -y concurrently --kill-others -n "backend,frontend" -c "magenta,cyan" "make run-backend" "make run-frontend"

test:
	@echo "Running backend tests..."
	cd apps/core && go test ./...

wire:
	@echo "Generating wire code..."
	cd apps/core/cmd/api && go run github.com/google/wire/cmd/wire gen

build:
	@echo "Building backend..."
	cd apps/core && go build -o tmp/main ./cmd/api
	@echo "Building frontend..."
	cd apps/portal && npm run build
