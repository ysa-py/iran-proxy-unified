#!/bin/bash
set -e

echo "🇮🇷 Building Iran Proxy Unified..."

# Navigate to project root
cd /workspaces/iran-proxy-unified/iran-proxy-unified-ultimate

# Clean and prepare
export GOSUMDB=off
export GOPROXY=direct
export GO111MODULE=on

echo "📦 Cleaning cache..."
go clean -modcache || true

echo "📋 Running go mod tidy..."
go mod tidy

echo "🏗️  Building project..."
cd src
go build -v -o ../iran-proxy . 2>&1

# Check if build succeeded
if [ -f ../iran-proxy ]; then
    echo ""
    echo "✅  Build successful!"
    ../iran-proxy -version
    
    ls -lh ../iran-proxy
else
    echo "❌ Build failed!"
    exit 1
fi
