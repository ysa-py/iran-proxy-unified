#!/usr/bin/env bash

# ═══════════════════════════════════════════════════════════════════════════
# Iran Proxy Ultimate - Build Script
# ═══════════════════════════════════════════════════════════════════════════
# Description: Comprehensive build script for Iran Proxy Ultimate System
# Version: 3.2.0
# Usage: ./scripts/build.sh [options]
# ═══════════════════════════════════════════════════════════════════════════

set -euo pipefail

# ═══════════════════════════════════════════════════════════════════════════
# Configuration Variables
# ═══════════════════════════════════════════════════════════════════════════

PROJECT_NAME="iran-proxy-ultimate"
VERSION="3.2.0"
BUILD_DIR="bin"
SRC_DIR="src"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build information
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_ID="build-$(date +%Y%m%d-%H%M%S)"
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# ═══════════════════════════════════════════════════════════════════════════
# Helper Functions
# ═══════════════════════════════════════════════════════════════════════════

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.21 or higher."
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_success "Go version $GO_VERSION detected"
    
    if ! command -v git &> /dev/null; then
        log_warning "Git is not installed. Version information may be incomplete."
    fi
}

clean_build() {
    log_info "Cleaning previous build artifacts..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    log_success "Build directory cleaned"
}

download_dependencies() {
    log_info "Downloading Go dependencies..."
    cd "$SRC_DIR"
    
    if [ -f go.mod ]; then
        go mod download
        go mod verify
        log_success "Dependencies downloaded and verified"
    else
        log_error "go.mod not found in $SRC_DIR"
        exit 1
    fi
    
    cd ..
}

build_binary() {
    local os=$1
    local arch=$2
    local output_name=$3
    
    log_info "Building for $os/$arch..."
    
    cd "$SRC_DIR"
    
    GOOS=$os GOARCH=$arch go build \
        -v \
        -ldflags="-s -w \
            -X main.Version=$VERSION \
            -X main.BuildTime=$BUILD_TIME \
            -X main.BuildID=$BUILD_ID \
            -X main.GitCommit=$GIT_COMMIT" \
        -trimpath \
        -o "../$BUILD_DIR/$output_name" \
        main.go main_iran.go
    
    if [ $? -eq 0 ]; then
        log_success "Built $output_name"
    else
        log_error "Failed to build $output_name"
        exit 1
    fi
    
    cd ..
}

build_all_platforms() {
    log_info "Building for all supported platforms..."
    
    # Linux AMD64
    build_binary "linux" "amd64" "$PROJECT_NAME-linux-amd64"
    
    # Linux ARM64
    build_binary "linux" "arm64" "$PROJECT_NAME-linux-arm64"
    
    # macOS AMD64
    build_binary "darwin" "amd64" "$PROJECT_NAME-darwin-amd64"
    
    # macOS ARM64 (Apple Silicon)
    build_binary "darwin" "arm64" "$PROJECT_NAME-darwin-arm64"
    
    # Windows AMD64
    build_binary "windows" "amd64" "$PROJECT_NAME-windows-amd64.exe"
    
    log_success "All platform builds completed"
}

build_current_platform() {
    log_info "Building for current platform..."
    
    cd "$SRC_DIR"
    
    go build \
        -v \
        -ldflags="-s -w \
            -X main.Version=$VERSION \
            -X main.BuildTime=$BUILD_TIME \
            -X main.BuildID=$BUILD_ID \
            -X main.GitCommit=$GIT_COMMIT" \
        -trimpath \
        -o "../$BUILD_DIR/$PROJECT_NAME" \
        main.go main_iran.go
    
    if [ $? -eq 0 ]; then
        log_success "Built $PROJECT_NAME for current platform"
    else
        log_error "Failed to build $PROJECT_NAME"
        exit 1
    fi
    
    cd ..
}

run_tests() {
    log_info "Running tests..."
    
    cd "$SRC_DIR"
    
    go test -v -race -coverprofile=coverage.out ./...
    
    if [ $? -eq 0 ]; then
        log_success "All tests passed"
        
        # Generate coverage report
        go tool cover -func=coverage.out
        
        log_info "Coverage report generated: $SRC_DIR/coverage.out"
    else
        log_error "Tests failed"
        exit 1
    fi
    
    cd ..
}

create_checksums() {
    log_info "Creating checksums..."
    
    cd "$BUILD_DIR"
    
    for file in *; do
        if [ -f "$file" ]; then
            sha256sum "$file" > "$file.sha256"
            log_success "Created checksum for $file"
        fi
    done
    
    cd ..
}

show_build_info() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "  Build Information"
    echo "═══════════════════════════════════════════════════════════════"
    echo "  Project:     $PROJECT_NAME"
    echo "  Version:     $VERSION"
    echo "  Build ID:    $BUILD_ID"
    echo "  Build Time:  $BUILD_TIME"
    echo "  Git Commit:  $GIT_COMMIT"
    echo "  Output:      $BUILD_DIR/"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
}

show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Build script for Iran Proxy Ultimate System

OPTIONS:
    -h, --help          Show this help message
    -c, --clean         Clean build directory before building
    -a, --all           Build for all supported platforms
    -t, --test          Run tests before building
    -C, --checksums     Create SHA256 checksums for binaries
    -v, --verbose       Enable verbose output

EXAMPLES:
    $0                  Build for current platform
    $0 --all            Build for all platforms
    $0 --clean --test   Clean, test, and build
    $0 -a -C            Build all platforms and create checksums

EOF
}

# ═══════════════════════════════════════════════════════════════════════════
# Main Script Logic
# ═══════════════════════════════════════════════════════════════════════════

main() {
    local build_all=false
    local run_test=false
    local do_clean=false
    local create_sums=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -c|--clean)
                do_clean=true
                shift
                ;;
            -a|--all)
                build_all=true
                shift
                ;;
            -t|--test)
                run_test=true
                shift
                ;;
            -C|--checksums)
                create_sums=true
                shift
                ;;
            -v|--verbose)
                set -x
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Show build information
    show_build_info
    
    # Check prerequisites
    check_prerequisites
    
    # Clean if requested
    if [ "$do_clean" = true ]; then
        clean_build
    else
        mkdir -p "$BUILD_DIR"
    fi
    
    # Download dependencies
    download_dependencies
    
    # Run tests if requested
    if [ "$run_test" = true ]; then
        run_tests
    fi
    
    # Build
    if [ "$build_all" = true ]; then
        build_all_platforms
    else
        build_current_platform
    fi
    
    # Create checksums if requested
    if [ "$create_sums" = true ]; then
        create_checksums
    fi
    
    # Success message
    echo ""
    log_success "Build process completed successfully!"
    log_info "Binaries available in: $BUILD_DIR/"
    
    # List built binaries
    echo ""
    log_info "Built binaries:"
    ls -lh "$BUILD_DIR" | grep -v ".sha256$" | grep -v "^total"
    echo ""
}

# Run main function
main "$@"
