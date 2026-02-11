# 🇮🇷 Iran Proxy Ultimate System - Enterprise v3.2.0

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Version](https://img.shields.io/badge/version-3.2.0-blue)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)
![Iran Optimized](https://img.shields.io/badge/Iran-Optimized-red)

**Enterprise-grade proxy system specifically optimized for bypassing Iran's advanced Deep Packet Inspection (DPI) and internet filtering infrastructure.**

## 🌟 Key Features

### 🇮🇷 Iran-Specific Optimizations
- **Advanced DPI Bypass**: Multiple layers of evasion techniques specifically designed for Iran's filtering infrastructure
- **Reality Protocol Support**: Next-generation TLS camouflage technology
- **xhttp Transport**: Advanced HTTP obfuscation for maximum stealth
- **SNI Fragmentation**: Bypass Server Name Indication filtering
- **TLS Fingerprint Spoofing**: Mimic legitimate HTTPS traffic patterns
- **Multi-Endpoint Validation**: Test connectivity across multiple Iranian ISPs

### 🚀 Performance & Reliability
- **Intelligent Load Balancing**: Automatic distribution across available proxies
- **Adaptive Timeout Management**: Dynamic adjustment based on network conditions
- **Multi-Tier Fallback System**: Automatic failover when primary proxies are blocked
- **Self-Healing Capabilities**: Automatic recovery from connection failures
- **Circuit Breaker Pattern**: Prevent cascade failures

### 📊 Monitoring & Analytics
- **Real-Time Health Scoring**: Continuous evaluation of proxy quality
- **Comprehensive Metrics**: Detailed statistics on success rates and performance
- **Anomaly Detection**: Automatic identification of filtering patterns
- **Performance Tracking**: Historical analysis of proxy effectiveness

### 🔒 Security & Privacy
- **No Logging**: Zero retention of user data or connection information
- **Encrypted Configurations**: All proxy configurations are encrypted
- **Open Source**: Complete transparency in implementation
- **Regular Security Audits**: Automated security scanning in CI/CD pipeline

## 📋 Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
- [GitHub Actions Integration](#github-actions-integration)
- [Docker Deployment](#docker-deployment)
- [Development](#development)
- [Architecture](#architecture)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## 🚀 Installation

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, for using Makefile)

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/iran-proxy-ultimate.git
cd iran-proxy-ultimate

# Build using Make
make build

# Or build manually
cd src
go build -o ../bin/iran-proxy-ultimate main.go main_iran.go
```

### Using Docker

```bash
# Build Docker image
docker build -t iran-proxy-ultimate:latest .

# Run container
docker run --rm iran-proxy-ultimate:latest
```

### Pre-built Binaries

Download the latest release from the [Releases](https://github.com/yourusername/iran-proxy-ultimate/releases) page.

## ⚡ Quick Start

### Basic Usage

```bash
# Run with default Iran optimizations
./bin/iran-proxy-ultimate --iran-mode

# Run with maximum DPI evasion
./bin/iran-proxy-ultimate \
  --iran-mode \
  --dpi-evasion-level maximum \
  --performance-mode balanced

# Run with custom settings
./bin/iran-proxy-ultimate \
  --iran-mode \
  --max-concurrent 200 \
  --timeout 15 \
  --protocols "vmess,vless,trojan,shadowsocks" \
  --enable-monitoring \
  --enable-self-healing
```

### Using Makefile

```bash
# Build and run with Iran optimizations
make run-iran

# Build for all platforms
make build-all

# Run tests
make test

# Generate coverage report
make coverage

# Show all available commands
make help
```

## ⚙️ Configuration

### Command Line Options

| Flag | Description | Default | Values |
|------|-------------|---------|--------|
| `--iran-mode` | Enable Iran-specific optimizations | false | boolean |
| `--dpi-evasion-level` | Level of DPI evasion techniques | aggressive | standard, aggressive, maximum |
| `--performance-mode` | Performance optimization mode | balanced | speed, balanced, quality |
| `--max-concurrent` | Maximum concurrent connections | 200 | 50-500 |
| `--timeout` | Connection timeout in seconds | 15 | 5-30 |
| `--protocols` | Protocols to test (comma-separated) | all | vmess,vless,trojan,shadowsocks |
| `--enable-monitoring` | Enable monitoring and metrics | true | boolean |
| `--enable-self-healing` | Enable self-healing capabilities | true | boolean |
| `--enable-fallback` | Enable multi-tier fallback | true | boolean |
| `--output` | Output directory for configs | ./configs | path |
| `--verbose` | Enable verbose logging | false | boolean |

### Environment Variables

```bash
# Core settings
export IRAN_MODE=true
export DPI_EVASION_LEVEL=aggressive
export PERFORMANCE_MODE=balanced

# Performance tuning
export MAX_CONCURRENT=200
export TIMEOUT_SECONDS=15
export RETRY_ATTEMPTS=3

# Feature flags
export ENABLE_MONITORING=true
export ENABLE_SELF_HEALING=true
export ENABLE_FALLBACK=true
export CACHE_ENABLED=true
```

## 🎯 Usage

### Generating Proxy Configurations

```bash
# Generate configs with Iran optimizations
./bin/iran-proxy-ultimate \
  --iran-mode \
  --output ./configs \
  --enable-monitoring

# Generate configs for specific protocols
./bin/iran-proxy-ultimate \
  --iran-mode \
  --protocols "vmess,vless" \
  --output ./configs
```

### Testing Proxy Quality

```bash
# Test all proxies with health scoring
./bin/iran-proxy-ultimate \
  --iran-mode \
  --enable-health-scoring \
  --quality-threshold 70

# Deep analysis mode
./bin/iran-proxy-ultimate \
  --iran-mode \
  --deep-analysis \
  --verbose
```

### Emergency Recovery

```bash
# Run emergency recovery check
./bin/iran-proxy-ultimate \
  --iran-mode \
  --emergency-recovery \
  --enable-fallback
```

## 🔄 GitHub Actions Integration

The project includes a comprehensive GitHub Actions workflow that automatically:

- Runs code quality and security checks
- Builds and tests the application
- Generates and optimizes proxy configurations
- Updates configurations on schedule (every 15 minutes, 4 hours, and daily)
- Provides comprehensive health reporting
- Builds and pushes Docker images

### Workflow Configuration

The workflow is located at `.github/workflows/iran-proxy-ultimate.yml` and supports:

**Scheduled Runs:**
- Quick config updates every 15 minutes
- Comprehensive proxy checks every 4 hours
- Deep analysis daily at 3 AM Tehran time
- Emergency recovery checks every hour

**Manual Triggers:**
Use GitHub Actions workflow dispatch with customizable parameters including check mode, performance settings, DPI evasion level, and more.

**Automatic Triggers:**
- On push to main, develop, feature, hotfix, and release branches
- On pull requests to main and develop
- On version tags

### Required Secrets

No additional secrets are required beyond the default `GITHUB_TOKEN` which is automatically provided.

### Outputs

The workflow generates:
- Optimized proxy configurations in multiple formats
- Quality metrics and statistics
- Comprehensive health reports
- Docker images (optional)
- Build artifacts

## 🐳 Docker Deployment

### Using Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  iran-proxy-ultimate:
    build: .
    container_name: iran-proxy-ultimate
    restart: unless-stopped
    environment:
      - IRAN_MODE=true
      - DPI_EVASION_LEVEL=aggressive
      - PERFORMANCE_MODE=balanced
      - ENABLE_MONITORING=true
    volumes:
      - ./configs:/app/configs
      - ./metrics:/app/metrics
    command:
      - "--iran-mode"
      - "--performance-mode"
      - "balanced"
      - "--dpi-evasion-level"
      - "aggressive"
      - "--enable-monitoring"
```

Run with:

```bash
docker-compose up -d
```

### Kubernetes Deployment

Example deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: iran-proxy-ultimate
spec:
  replicas: 3
  selector:
    matchLabels:
      app: iran-proxy-ultimate
  template:
    metadata:
      labels:
        app: iran-proxy-ultimate
    spec:
      containers:
      - name: iran-proxy-ultimate
        image: ghcr.io/iran-proxy-ultimate:latest
        args:
        - "--iran-mode"
        - "--performance-mode"
        - "balanced"
        - "--dpi-evasion-level"
        - "aggressive"
        env:
        - name: ENABLE_MONITORING
          value: "true"
        - name: ENABLE_SELF_HEALING
          value: "true"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

## 💻 Development

### Project Structure

```
iran-proxy-unified-ultimate/
├── .github/
│   └── workflows/
│       └── iran-proxy-ultimate.yml    # Main CI/CD workflow
├── src/
│   ├── main.go                        # Main application entry
│   ├── main_iran.go                   # Iran-specific optimizations
│   ├── proxy_checker_iran.go          # Proxy checking logic
│   ├── config_generator.go            # Configuration generation
│   ├── config_generator_ai.go         # AI-powered config optimization
│   ├── ai_anti_dpi.go                 # DPI bypass techniques
│   ├── sni_fragmentation.go           # SNI fragmentation
│   ├── utls_fingerprint_spoofing.go   # TLS fingerprinting
│   ├── advanced_health_scoring.go     # Health scoring system
│   ├── monitoring.go                  # Monitoring and metrics
│   └── protocols.go                   # Protocol handlers
├── configs/                           # Generated configurations
├── metrics/                           # Performance metrics
├── stats/                             # Statistics
├── docs/                              # Documentation
├── scripts/                           # Utility scripts
├── tests/                             # Test files
├── Dockerfile                         # Docker configuration
├── Makefile                          # Build automation
├── go.mod                            # Go module definition
└── README.md                         # This file
```

### Building

```bash
# Standard build
make build

# Build for all platforms
make build-all

# Development build with debug symbols
cd src
go build -gcflags="all=-N -l" -o ../bin/iran-proxy-ultimate-debug main.go main_iran.go
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific test
cd src
go test -v -run TestSpecificFunction
```

### Code Quality

```bash
# Format code
make fmt

# Run linters
make lint

# Run all checks
make check
```

## 🏗️ Architecture

### Core Components

**Proxy Checker:** Validates proxy connectivity and measures performance across multiple endpoints including Iranian ISPs.

**Config Generator:** Creates optimized configurations for various protocols with Iran-specific enhancements.

**DPI Bypass Engine:** Implements multiple evasion techniques including SNI fragmentation, TLS fingerprint spoofing, and protocol obfuscation.

**Health Scoring System:** Continuously evaluates proxy quality based on latency, success rate, and stability.

**Monitoring System:** Collects and analyzes metrics for performance optimization and anomaly detection.

**Self-Healing Module:** Automatically detects and recovers from failures, adjusting strategies as needed.

### Data Flow

The system follows this general flow: Fetch proxy sources, validate connectivity, apply Iran optimizations, generate health scores, create optimized configurations, monitor performance, and trigger self-healing when needed.

## 🔧 Troubleshooting

### Common Issues

**Proxies not working:**
- Ensure Iran mode is enabled with `--iran-mode`
- Try increasing DPI evasion level to `maximum`
- Check if your ISP is blocking specific ports
- Enable fallback mode for additional redundancy

**Low success rate:**
- Adjust timeout with `--timeout` parameter
- Reduce concurrent connections with `--max-concurrent`
- Enable self-healing with `--enable-self-healing`
- Switch performance mode to `quality` for better stability

**Build errors:**
- Verify Go version is 1.21 or higher
- Run `go mod download` to ensure dependencies are available
- Check for missing build tools on your system

### Logs and Debugging

```bash
# Enable verbose logging
./bin/iran-proxy-ultimate --iran-mode --verbose

# Check metrics
cat metrics/quality-metrics.json

# Review statistics
cat stats/proxy-stats.json
```

## 🤝 Contributing

Contributions are welcome and appreciated. Please follow these guidelines:

### Development Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linters (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Standards

- Follow Go best practices and idioms
- Write tests for new functionality
- Update documentation as needed
- Ensure all CI checks pass
- Add comments for complex logic

## 📄 License

This project is licensed under the MIT License. See the LICENSE file for details.

## 🙏 Acknowledgments

- The V2Ray and Xray communities for protocol implementations
- Contributors to Iran proxy filtering research
- Open source security tools and libraries

## 📞 Support

For issues, questions, or suggestions:
- Open an issue on GitHub
- Check existing documentation in the `docs/` directory
- Review closed issues for similar problems

---

**Developed with ❤️ for the Iranian community**

**Version:** 3.2.0 Ultimate Edition  
**Last Updated:** 2024

---

## 📊 Statistics

The system automatically generates and updates statistics including total proxies tested, success rate percentage, average response time, protocols supported, DPI bypass success rate, and last update timestamp.

Check the `stats/` and `metrics/` directories for detailed analytics.
