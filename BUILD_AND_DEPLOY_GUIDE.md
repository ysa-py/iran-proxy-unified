# Iran Proxy Unified - Build & Deployment Guide

## 🇮🇷 Project Overview

**Iran Proxy Unified** is an enterprise-grade proxy management and configuration system with advanced DPI evasion capabilities specifically optimized for the Iranian network environment.

### Core Features

#### 1. 🛡️ Advanced Anti-DPI Technologies
- **uTLS Fingerprint Spoofing**: Randomizes TLS fingerprints to mimic legitimate browsers
- **SNI Fragmentation (Adaptive)**: Fragments Server Name Indication to evade detection
- **AI-Powered DPI Evasion**: Machine learning-based pattern adaptation
- **Timing Obfuscation**: Dynamic delay injection to defeat timing-based detection
- **Packet Padding (Dynamic)**: Intelligent packet size manipulation

#### 2. 🔧 Core Proxy Management
- **Multi-Endpoint Testing**: Tests against 4+ endpoints optimized for Iran
- **Intelligent Config Generation**: Automatically generates optimized configurations
- **Advanced Health Scoring**: Real-time proxy quality evaluation
- **Protocol Support**: VMess, VLESS, Trojan, ShadowSocks, and more
- **Multi-tier Fallback**: Automatic failover to alternative proxies

#### 3. 📊 System Intelligence
- **Self-Healing Mechanisms**: Automatically recovers from failures
- **Real-time Monitoring**: Comprehensive metrics and performance tracking
- **Performance Modes**: Speed, Balanced, and Quality optimization
- **Emergency Recovery**: Built-in disaster recovery procedures
- **Deep Analysis**: Advanced debugging and diagnostic capabilities

---

## 📦 Installation & Build

### Prerequisites

```bash
# Go 1.21 or later
go version

# Verify Go installation
which go
```

### Method 1: Using Make (Recommended)

```bash
cd /path/to/iran-proxy-unified-ultimate

# View available targets
make help

# Build the project
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Method 2: Using the Build Script

```bash
cd /path/to/iran-proxy-unified-ultimate

# Make script executable
chmod +x build-deploy.sh

# Run the build
./build-deploy.sh
```

### Method 3: Manual Build

```bash
cd /path/to/iran-proxy-unified-ultimate/src

# Set required environment variables
export GOSUMDB=off
export GO111MODULE=on

# Download dependencies
go mod download
go mod verify

# Build the binary
go build -v -o ../iran-proxy .

# Verify the build
../iran-proxy -version
```

### Method 4: Using the User Command (From Request)

```bash
cd /path/to/iran-proxy-unified-ultimate/src

# Execute exactly as requested
GOSUMDB=off go build -v -o ../iran-proxy .
```

---

## 🚀 Usage

### Basic Proxy Checking

```bash
./iran-proxy \
  -proxy-file edge/assets/p-list-february.txt \
  -output-file sub/ProxyIP-Daily.md \
  -iran-mode \
  -max-concurrent 100 \
  -timeout 10
```

### Advanced Mode with Anti-DPI

```bash
./iran-proxy \
  -iran-mode \
  -dpi-evasion-level maximum \
  -performance-mode balanced \
  -enable-self-healing \
  -enable-monitoring \
  -generate-configs \
  -test-configs \
  -max-concurrent 150
```

### Config-Only Generation

```bash
./iran-proxy \
  -configs-only \
  -generate-configs \
  -dpi-evasion-level aggressive \
  -config-output configs/iran-configs.txt
```

### Emergency Recovery Mode

```bash
./iran-proxy \
  -iran-mode \
  -emergency-mode \
  -deep-analysis \
  -enable-self-healing
```

---

## 📋 Command-Line Options

### File Paths
```
-proxy-file          Path to proxy list (CSV format)
-output-file         Path to output markdown file
-config-output       Path for generated configurations
```

### Performance Settings
```
-max-concurrent      Maximum concurrent connections (50-500, default: 100)
-timeout             Timeout in seconds (5-30, default: 10)
```

### Operation Modes
```
-iran-mode           Enable Iran-specific optimizations (default: true)
-generate-configs    Generate configs from active proxies (default: true)
-test-configs        Test generated configurations (default: true)
-configs-only        Only generate configs, skip proxy check
```

### Advanced Features
```
-performance-mode    speed, balanced, quality (default: balanced)
-dpi-evasion-level   standard, aggressive, maximum (default: aggressive)
-enable-fallback     Enable multi-tier fallback (default: true)
-enable-self-healing Enable self-healing (default: true)
-enable-monitoring   Enable monitoring & metrics (default: true)
```

### Recovery Options
```
-emergency-mode      Enable emergency recovery procedures
-deep-analysis       Run deep analysis mode
```

### Display Options
```
-help                Show help message
-version             Show version information
-verbose             Enable verbose output
```

---

## 🔌 Environment Variables

Configure behavior via environment variables:

```bash
# File Paths
export PROXY_FILE=edge/assets/p-list-february.txt
export OUTPUT_FILE=sub/ProxyIP-Daily.md
export CONFIG_OUTPUT=configs/iran-configs.txt

# Performance Settings
export MAX_CONCURRENT=100
export TIMEOUT=10

# Operation Modes
export IRAN_MODE=true
export GENERATE_CONFIGS=true
export TEST_CONFIGS=true
export CONFIGS_ONLY=false

# Advanced Features
export PERFORMANCE_MODE=balanced
export DPI_EVASION_LEVEL=aggressive
export ENABLE_FALLBACK=true
export ENABLE_SELF_HEALING=true
export ENABLE_MONITORING=true

# Recovery
export EMERGENCY_MODE=false
export DEEP_ANALYSIS=false

# Display
export VERBOSE=false

# Then run
./iran-proxy
```

---

## 🏗️ Project Structure

```
src/
├── main.go                          # Main entry point with CLI
├── main_iran.go                     # Iran-specific main variant
├── proxy_checker_iran.go            # Core proxy checking logic
├── enhanced_proxy_checker.go        # Advanced proxy checker with anti-DPI
├── enhanced_types.go                # Type definitions for enhanced features
│
├── config_generator.go              # Config generation logic
├── config_generator_ai.go           # AI-powered config generation
├── config_tester.go                 # Configuration testing
├── config_writer.go                 # Output formatting and writing
├── config_writer_enhanced.go        # Enhanced output capabilities
│
├── protocols.go                     # Protocol definitions and utilities
├── advanced_integration.go          # Advanced anti-DPI integration
├── advanced_health_scoring.go       # Health scoring algorithms
├── advanced_test.go                 # Advanced testing utilities
│
├── ai_anti_dpi.go                   # AI-powered DPI evasion engine
├── utls_fingerprint_spoofing.go     # uTLS fingerprint manipulation
├── sni_fragmentation.go             # SNI fragmentation implementation
│
├── monitoring.go                    # System monitoring and metrics
├── go.mod                           # Go module definition
└── go.sum                           # Go module checksums
```

---

## 🔍 Build Process

### Phase 1: Environment Setup
- Verify Go installation and version
- Check source directory structure
- Count Go source files

### Phase 2: Dependency Management
- Download Go dependencies
- Verify checksums
- Validate go.mod

### Phase 3: Code Quality
- Run static analysis (go vet)
- Check for obvious issues

### Phase 4: Compilation
- Compile with optimization flags
- Generate optimized binary
- Apply version information

### Phase 5: Verification
- Verify binary creation
- Check binary size
- Test basic functionality

---

## 📊 Build Statistics

| Metric | Value |
|--------|-------|
| Go Version | 1.21+ |
| Source Files | 20+ Go files |
| Total Lines | 10,000+ |
| External Dependencies | 15+ packages |
| Binary Size | ~15-20 MB (unstripped) |
| Build Time | 30-60 seconds |

---

## 🧪 Testing

### Run Unit Tests

```bash
cd src
go test -v ./...
```

### Generate Coverage Report

```bash
make coverage
# Opens coverage.html
```

### Run Full Test Suite

```bash
make test
```

---

## 🐳 Docker Build

### Build Docker Image

```bash
docker build -t iran-proxy:latest .
```

### Using Docker Compose

```bash
docker-compose up -d
```

### Run with Docker

```bash
docker run -it --rm \
  -v $(pwd)/configs:/app/configs \
  iran-proxy:latest \
  -iran-mode \
  -dpi-evasion-level maximum
```

---

## 🔧 Troubleshooting

### Build Issues

**Issue**: `go: no Go files in /...`
```bash
# Solution: Check working directory
cd src
go build -v -o ../iran-proxy .
```

**Issue**: `missing go.sum`
```bash
# Solution: Regenerate checksums
go mod tidy
go mod verify
```

**Issue**: `conflicting versions`
```bash
# Solution: Clean and reinstall
go clean -modcache
go mod download
```

### Runtime Issues

**Issue**: Proxy timeout
```bash
# Solution: Increase timeout
./iran-proxy -timeout 20 -max-concurrent 50
```

**Issue**: Connection refused
```bash
# Solution: Enable emergency mode
./iran-proxy -emergency-mode -enable-self-healing
```

---

## 📚 Documentation

- [CODE_FIXES_COMPLETE_FA.md](./CODE_FIXES_COMPLETE_FA.md) - فارسی: تمام اصلاحات کد
- [CODE_FIXES_SUMMARY_EN.md](./CODE_FIXES_SUMMARY_EN.md) - English: Code Fixes Summary
- [COMPLETE_FIXES_DOCUMENTATION.md](./COMPLETE_FIXES_DOCUMENTATION.md) - Comprehensive documentation
- [README.md](./README.md) - Project README
- [README-FA.md](./README-FA.md) - پروژه راهنما

---

## 🎯 Advanced Configuration

### For Maximum DPI Evasion

```bash
./iran-proxy \
  -dpi-evasion-level maximum \
  -enable-self-healing \
  -enable-fallback \
  -enable-monitoring \
  -deep-analysis
```

### For Speed Optimization

```bash
./iran-proxy \
  -performance-mode speed \
  -max-concurrent 300 \
  -timeout 5 \
  -dpi-evasion-level standard
```

### For Balanced Production

```bash
./iran-proxy \
  -performance-mode balanced \
  -dpi-evasion-level aggressive \
  -enable-self-healing \
  -enable-monitoring \
  -max-concurrent 100
```

---

## 📝 Version History

### v3.2.0 (Current)
- ✅ Unified proxy and config system
- ✅ AI-powered DPI evasion
- ✅ Advanced health scoring
- ✅ Self-healing mechanisms
- ✅ Real-time monitoring

---

## 📞 Support

For issues, questions, or contributions:
- Check the documentation files
- Review existing code comments
- Examine test files for examples

---

**🇮🇷 Built for Iran's Internet Freedom 🇮🇷**

*Enterprise-grade proxy management with advanced DPI evasion capabilities*
