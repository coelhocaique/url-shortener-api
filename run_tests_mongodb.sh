#!/bin/bash

# MongoDB Test Runner Script
# This script runs the test suite with MongoDB integration

set -e

echo "🧪 Running URL Shortener API Tests with MongoDB"
echo "=============================================="

# Check if MongoDB is running
echo "📋 Checking MongoDB connection..."
if ! docker ps | grep -q mongodb; then
    echo "⚠️  MongoDB container not found. Starting MongoDB..."
    docker run -d --name test-mongodb -p 27017:27017 mongo:7.0
    echo "⏳ Waiting for MongoDB to be ready..."
    sleep 5
    MONGO_STARTED=true
else
    echo "✅ MongoDB is already running"
    MONGO_STARTED=false
fi

# Set test environment variables
export MONGO_URI="mongodb://localhost:27017"
export DATABASE_NAME="url_shortener_test"
export PORT="8080"

echo ""
echo "🔧 Running Unit Tests..."
echo "========================"

# Run unit tests
echo "📦 Testing URL Storage..."
go test -v ./tests/unit/services/url_storage_test.go

echo ""
echo "📦 Testing URL Service..."
go test -v ./tests/unit/services/url_service_test.go

echo ""
echo "📦 Testing URL Handlers..."
go test -v ./tests/unit/handlers/url_handler_test.go

echo ""
echo "🔧 Running Integration Tests..."
echo "==============================="

# Run integration tests
go test -v ./tests/integration/api_integration_test.go

echo ""
echo "🔧 Running All Tests with Coverage..."
echo "===================================="

# Run all tests with coverage
go test -v -coverprofile=coverage.out ./tests/...

# Generate coverage report
if [ -f coverage.out ]; then
    echo ""
    echo "📊 Coverage Report:"
    go tool cover -func=coverage.out
    echo ""
    echo "📈 HTML coverage report generated: coverage.html"
    go tool cover -html=coverage.out -o coverage.html
fi

# Cleanup
if [ "$MONGO_STARTED" = true ]; then
    echo ""
    echo "🧹 Cleaning up test MongoDB container..."
    docker stop test-mongodb
    docker rm test-mongodb
fi

echo ""
echo "✅ All tests completed successfully!"
echo "🎉 MongoDB integration tests passed!"
