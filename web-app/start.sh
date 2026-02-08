#!/bin/bash
set -e

echo "📦 Installing dependencies..."
go mod download

echo "🔨 Building IPO Pilot..."
go build -o ipo-pilot .

echo "✓ Build successful!"
echo "🚀 Starting IPO Pilot on :8080..."
exec ./ipo-pilot
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
