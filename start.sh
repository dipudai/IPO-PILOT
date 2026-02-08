#!/bin/bash

echo "🚀 IPO Pilot Web Platform - Railway Startup"
echo "==========================================="
echo ""

# Navigate to web-app directory
cd web-app || exit 1

echo "📁 Current directory: $(pwd)"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi

echo "✓ Go version: $(go version)"
echo ""

# Build the application
echo "🔨 Building application..."
go build -o ipo-pilot .

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "🚀 Starting IPO Pilot..."
    ./ipo-pilot
else
    echo "❌ Build failed"
    exit 1
fi
