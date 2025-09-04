#!/bin/bash

echo "🚀 Running URL Shortener API Tests"
echo "=================================="
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy
echo ""

# Run unit tests
echo "🧪 Running Unit Tests..."
echo "------------------------"
go test -v ./tests/unit/...
echo ""

# Run integration tests
echo "🔗 Running Integration Tests..."
echo "-------------------------------"
go test -v ./tests/integration/...
echo ""

# Run all tests with coverage
echo "📊 Running All Tests with Coverage..."
echo "-------------------------------------"
go test -v -coverprofile=coverage.out ./...
echo ""

# Show coverage summary
if command -v go tool cover &> /dev/null; then
    echo "📈 Coverage Summary:"
    echo "-------------------"
    go tool cover -func=coverage.out
    echo ""
    
    echo "📊 Coverage Report:"
    echo "------------------"
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report saved to coverage.html"
else
    echo "⚠️  go tool cover not available, skipping coverage report"
fi

echo ""
echo "✅ All tests completed!"
