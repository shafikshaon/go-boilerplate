.PHONY: help setup db-up db-down run test clean

# Default target
help:
	@echo "🚀 go-boilerplate - Go API with PostgreSQL & Redis"
	@echo "=============================================="
	@echo "Available commands:"
	@echo "  make setup     - Install dependencies"
	@echo "  make db-up     - Start PostgreSQL and Redis"
	@echo "  make db-down   - Stop databases"
	@echo "  make run       - Run the application"
	@echo "  make test      - Test the API"
	@echo "  make clean     - Clean up containers and data"
	@echo "  make dev       - Start databases and run app"
	@echo ""
	@echo "🏃‍♂️ Quick start: make setup && make dev"

# Install dependencies
setup:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed!"

# Start databases
db-up:
	@echo "🐘 Starting PostgreSQL and Redis..."
	docker-compose up -d postgres redis
	@echo "⏳ Waiting for databases to be ready..."
	@sleep 5
	@echo "✅ Databases are running!"
	@echo "📊 PostgreSQL: localhost:5432 (user: postgres, pass: postgres, db: app_db)"
	@echo "🔴 Redis: localhost:6379"

# Stop databases
db-down:
	@echo "🛑 Stopping databases..."
	docker-compose down
	@echo "✅ Databases stopped!"

# Run the application
run:
	@echo "🚀 Starting Go application..."
	go run main.go

# Test the API
test:
	@echo "🧪 Testing API endpoints..."
	@chmod +x test_api.sh
	@./test_api.sh

# Clean up everything
clean:
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	docker system prune -f
	@echo "✅ Cleanup completed!"

# Development workflow
dev: db-up
	@echo "⏳ Waiting a bit more for databases..."
	@sleep 3
	@echo "🚀 Starting application..."
	@go run main.go

# Connect to PostgreSQL
psql:
	docker exec -it go-boilerplate_postgres psql -U postgres -d app_db

# Connect to Redis CLI
redis-cli:
	docker exec -it go-boilerplate_redis redis-cli

# Build binary
build:
	go build -o bin/app main.go

# Docker operations
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
