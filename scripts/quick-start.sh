#!/usr/bin/env bash

# ═══════════════════════════════════════════════════════════════════════════
# Iran Proxy Ultimate - Quick Start Script
# ═══════════════════════════════════════════════════════════════════════════
# Description: Automated setup and quick start for Iran Proxy Ultimate System
# Version: 3.2.0
# Usage: ./scripts/quick-start.sh
# ═══════════════════════════════════════════════════════════════════════════

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="Iran Proxy Ultimate"
VERSION="3.2.0"

# ═══════════════════════════════════════════════════════════════════════════
# Helper Functions
# ═══════════════════════════════════════════════════════════════════════════

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_banner() {
    clear
    cat << "EOF"
╔═══════════════════════════════════════════════════════════════════════════╗
║                                                                           ║
║            🇮🇷 Iran Proxy Ultimate System - Quick Start                    ║
║                       Version 3.2.0 - Enterprise                          ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
EOF
    echo ""
}

check_requirements() {
    log_info "Checking system requirements..."
    
    local missing_deps=()
    
    # Check Go
    if ! command -v go &> /dev/null; then
        missing_deps+=("Go (1.21+)")
    else
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        log_success "Go $GO_VERSION installed"
    fi
    
    # Check Git
    if ! command -v git &> /dev/null; then
        missing_deps+=("Git")
    else
        log_success "Git installed"
    fi
    
    # Check Make (optional)
    if command -v make &> /dev/null; then
        log_success "Make installed"
    else
        log_warning "Make not found (optional)"
    fi
    
    # Check Docker (optional)
    if command -v docker &> /dev/null; then
        log_success "Docker installed"
    else
        log_warning "Docker not found (optional)"
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing required dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        exit 1
    fi
    
    echo ""
}

setup_directories() {
    log_info "Setting up project directories..."
    
    mkdir -p configs/{merged,by-protocol,by-region,base64}
    mkdir -p stats
    mkdir -p metrics
    mkdir -p logs
    mkdir -p bin
    
    log_success "Directories created"
    echo ""
}

install_dependencies() {
    log_info "Installing Go dependencies..."
    
    cd src
    
    if [ -f go.mod ]; then
        go mod download
        go mod verify
        log_success "Dependencies installed and verified"
    else
        log_error "go.mod not found"
        exit 1
    fi
    
    cd ..
    echo ""
}

build_project() {
    log_info "Building Iran Proxy Ultimate..."
    
    if [ -f Makefile ] && command -v make &> /dev/null; then
        make build
    else
        ./scripts/build.sh
    fi
    
    log_success "Build completed"
    echo ""
}

show_menu() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}  Select Startup Mode${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo "  1) 🇮🇷 Iran Mode (Recommended)"
    echo "     - Maximum DPI bypass"
    echo "     - Aggressive evasion techniques"
    echo "     - Optimized for Iranian filtering"
    echo ""
    echo "  2) ⚡ Speed Mode"
    echo "     - Maximum concurrent connections"
    echo "     - Faster processing"
    echo "     - Less filtering bypass"
    echo ""
    echo "  3) 🎯 Quality Mode"
    echo "     - Maximum reliability"
    echo "     - Comprehensive checks"
    echo "     - Best proxy quality"
    echo ""
    echo "  4) ⚙️  Custom Configuration"
    echo "     - Specify your own parameters"
    echo ""
    echo "  5) 🐳 Docker Mode"
    echo "     - Run using Docker Compose"
    echo ""
    echo "  6) ❌ Exit"
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
}

run_iran_mode() {
    log_info "Starting in Iran Mode..."
    echo ""
    
    ./bin/iran-proxy-ultimate \
        --iran-mode \
        --dpi-evasion-level maximum \
        --performance-mode balanced \
        --max-concurrent 200 \
        --timeout 15 \
        --protocols "vmess,vless,trojan,shadowsocks" \
        --enable-monitoring \
        --enable-self-healing \
        --enable-fallback \
        --output ./configs
}

run_speed_mode() {
    log_info "Starting in Speed Mode..."
    echo ""
    
    ./bin/iran-proxy-ultimate \
        --performance-mode speed \
        --max-concurrent 500 \
        --timeout 10 \
        --output ./configs
}

run_quality_mode() {
    log_info "Starting in Quality Mode..."
    echo ""
    
    ./bin/iran-proxy-ultimate \
        --performance-mode quality \
        --max-concurrent 100 \
        --timeout 30 \
        --enable-monitoring \
        --enable-self-healing \
        --output ./configs
}

run_custom_mode() {
    log_info "Custom Configuration Mode"
    echo ""
    
    read -p "Enable Iran mode? (y/n): " iran_mode
    read -p "DPI evasion level (standard/aggressive/maximum): " dpi_level
    read -p "Performance mode (speed/balanced/quality): " perf_mode
    read -p "Max concurrent connections (50-500): " max_concurrent
    read -p "Timeout in seconds (5-30): " timeout
    
    local iran_flag=""
    if [[ "$iran_mode" =~ ^[Yy]$ ]]; then
        iran_flag="--iran-mode"
    fi
    
    echo ""
    log_info "Starting with custom configuration..."
    
    ./bin/iran-proxy-ultimate \
        $iran_flag \
        --dpi-evasion-level "$dpi_level" \
        --performance-mode "$perf_mode" \
        --max-concurrent "$max_concurrent" \
        --timeout "$timeout" \
        --enable-monitoring \
        --output ./configs
}

run_docker_mode() {
    log_info "Starting Docker Mode..."
    echo ""
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi
    
    log_info "Building and starting containers..."
    docker-compose up --build -d
    
    log_success "Docker containers started"
    log_info "View logs with: docker-compose logs -f"
}

show_results() {
    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}  Execution Complete!${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo "Generated configurations are available in:"
    echo "  - configs/merged/           (All configs combined)"
    echo "  - configs/by-protocol/      (Separated by protocol)"
    echo "  - configs/by-region/        (Separated by region)"
    echo "  - configs/base64/           (Base64 encoded)"
    echo ""
    echo "Metrics and statistics:"
    echo "  - metrics/quality-metrics.json"
    echo "  - stats/proxy-stats.json"
    echo ""
    log_success "Iran Proxy Ultimate is ready to use!"
    echo ""
}

# ═══════════════════════════════════════════════════════════════════════════
# Main Script Logic
# ═══════════════════════════════════════════════════════════════════════════

main() {
    print_banner
    check_requirements
    setup_directories
    install_dependencies
    build_project
    
    while true; do
        show_menu
        read -p "Enter your choice (1-6): " choice
        
        case $choice in
            1)
                run_iran_mode
                show_results
                break
                ;;
            2)
                run_speed_mode
                show_results
                break
                ;;
            3)
                run_quality_mode
                show_results
                break
                ;;
            4)
                run_custom_mode
                show_results
                break
                ;;
            5)
                run_docker_mode
                break
                ;;
            6)
                echo ""
                log_info "Exiting..."
                exit 0
                ;;
            *)
                log_error "Invalid choice. Please select 1-6."
                sleep 2
                ;;
        esac
    done
}

# Run main function
main "$@"
