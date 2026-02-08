#!/bin/bash
set -e

cd web-app

echo "📦 Installing dependencies..."
go mod download

echo "🔨 Building ITail..."
go build -o ipo-pilot .

echo "✓ Build successful!"
echo "🚀 Starting IPO Pilot on :8080..."
exec ./ipo-pilot
