#!/bin/bash

echo "🚀 Starting project with PostgreSQL and Redis..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start databases first
echo "🐘 Starting PostgreSQL and Redis..."
make db-up

# Install dependencies if needed
if [ ! -f "go.sum" ]; then
    echo "📦 Installing dependencies..."
    make setup
fi

# Wait a bit more for databases to be fully ready
echo "⏳ Waiting for databases to be fully ready..."
sleep 5

# Run the application
echo "🌟 Starting server on port 8080..."
echo "📊 PostgreSQL: localhost:5432 (app_db)"
echo "🔴 Redis: localhost:6379"
echo "🌐 API: http://localhost:8080"
echo ""
echo "🧪 Test endpoints: ./test_api.sh"
echo ""
make run
