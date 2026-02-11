#!/bin/bash

# Build check script for Iran Proxy Unified
set -e

echo "🔍 Checking Go environment..."
go version
echo ""

echo "📦 Checking go.mod..."
cat src/go.mod | head -10
echo ""

echo "🏗️  Building Iran Proxy Unified..."
cd /workspaces/iran-proxy-unified/iran-proxy-unified-ultimate
export GOSUMDB=off

echo "📝 Source files count:"
find src -name "*.go" -type f | wc -l
echo ""

echo "🔧 Running build..."
go build -v -o iran-proxy ./src 2>&1 | head -50

echo ""
echo "✅ Build check complete!"
