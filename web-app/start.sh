#!/bin/bash

echo "🚀 IPO Pilot Web Platform - Startup Script"
echo "=========================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "   Download from: https://golang.org/dl/"
    exit 1
fi

echo "✓ Go is installed: $(go version)"
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
echo "✓ Dependencies installed"
echo ""

# Build the application
echo "🔨 Building application..."
go build -o ipo-pilot-web .
if [ $? -eq 0 ]; then
    echo "✓ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi
echo ""

# Set environment variables
export GIN_MODE=debug
export PORT=8080

echo "🌐 Starting IPO Pilot Web Platform..."
echo "   URL: http://localhost:8080"
echo "   Admin: admin@ipopilot.com / admin123"
echo ""
echo "📝 Press Ctrl+C to stop the server"
echo "=========================================="
echo ""

# Run the application
./ipo-pilot-web
