# 🇮🇷 Iran Proxy Ultimate System - Project Summary

## Version 3.2.0 - Unified Edition

### Overview

This document provides a comprehensive summary of the Iran Proxy Ultimate System, a fully integrated and unified enterprise-grade proxy management platform specifically optimized for bypassing Iran's advanced Deep Packet Inspection (DPI) and internet filtering infrastructure. The system combines multiple previously separate workflows into a single, intelligent, and automated solution.

## Unified Architecture

### Integration Highlights

The project successfully unifies three distinct GitHub Actions workflows into a single comprehensive pipeline that handles all aspects of proxy management, testing, configuration generation, and deployment. The unified workflow includes preflight validation, code quality and security checks, automated building and testing, intelligent proxy checking with Iran-specific optimizations, configuration aggregation and optimization, comprehensive health monitoring and reporting, and optional Docker containerization and deployment.

### Core Components

The system is built on a modular architecture with clearly separated concerns. The source code directory contains all Go source files implementing core functionality including the main application entry point, Iran-specific optimizations, proxy checking logic, configuration generation with AI-powered optimization, DPI bypass techniques including SNI fragmentation and TLS fingerprint spoofing, advanced health scoring system, monitoring and metrics collection, and protocol handlers for multiple proxy types.

The GitHub Actions workflow provides comprehensive automation with scheduled execution at multiple intervals, manual triggering with extensive customization options, automatic execution on code changes and pull requests, intelligent caching for improved performance, multi-stage pipeline with proper dependency management, conditional execution based on configuration and results, comprehensive error handling and retry mechanisms, and detailed reporting and health checks.

### Key Features

#### Iran-Specific Optimizations

The system implements advanced DPI bypass capabilities through multiple layers of evasion techniques specifically designed for Iran's filtering infrastructure. This includes Reality Protocol support for next-generation TLS camouflage, xhttp Transport for advanced HTTP obfuscation, SNI Fragmentation to bypass Server Name Indication filtering, TLS Fingerprint Spoofing to mimic legitimate HTTPS traffic patterns, and multi-endpoint validation testing connectivity across multiple Iranian ISPs.

#### Performance and Reliability

The platform provides intelligent load balancing with automatic distribution across available proxies, adaptive timeout management that dynamically adjusts based on network conditions, a multi-tier fallback system for automatic failover when primary proxies are blocked, self-healing capabilities for automatic recovery from connection failures, and circuit breaker pattern implementation to prevent cascade failures.

#### Monitoring and Analytics

Comprehensive monitoring includes real-time health scoring with continuous evaluation of proxy quality, detailed statistics on success rates and performance, automatic identification of filtering patterns and network issues, historical analysis of proxy effectiveness, and quality assurance reporting with automated generation of quality metrics.

#### Security and Privacy

The system maintains a strong security posture with no logging or retention of user data or connection information, encrypted configurations for all proxy settings, complete transparency through open source implementation, and regular security audits through automated scanning in the CI/CD pipeline.

## Project Structure

The repository is organized into logical directories for maintainability and scalability. The .github/workflows directory contains the unified GitHub Actions workflow configuration. The src directory houses all Go source code organized by functionality. The configs directory stores generated proxy configurations in multiple formats including merged, by-protocol, by-region, and base64-encoded variants. The docs directory contains comprehensive documentation including quick start guides, deployment instructions, and API references. The scripts directory provides utility scripts for building, deployment, and quick setup. The metrics and stats directories store performance metrics and statistics collected during operation.

## Deployment Options

### Local Development

Developers can set up the project locally by cloning the repository, installing dependencies through either Make or manual Go commands, building the application using the provided Makefile or build scripts, and running the application with various configuration options for different use cases.

### Docker Deployment

The system includes comprehensive Docker support with a multi-stage Dockerfile optimized for security and size, a Docker Compose configuration for easy orchestration including optional monitoring stack with Prometheus and Grafana, health checks and automatic restart policies, resource limits and logging configuration, and volume mounts for persistent storage of configurations and metrics.

### Kubernetes Deployment

For production environments, the project provides Kubernetes deployment manifests with replica sets for high availability, resource requests and limits for proper scheduling, environment variable configuration for flexibility, and integration with existing Kubernetes infrastructure.

### GitHub Actions Automation

The automated workflow executes on multiple schedules including quick configuration updates every 15 minutes, comprehensive proxy checks every 4 hours, deep analysis daily at 3 AM Tehran time, and emergency recovery checks every hour. Manual triggering is supported through workflow dispatch with extensive customization options, and automatic execution occurs on code changes to main, develop, feature, hotfix, and release branches.

## Configuration and Customization

### Command Line Parameters

The application supports extensive configuration through command line flags. The --iran-mode flag enables Iran-specific optimizations and should be used for deployments targeting Iranian users. The --dpi-evasion-level parameter controls the aggressiveness of DPI bypass techniques with options for standard, aggressive, or maximum evasion. The --performance-mode setting selects between speed, balanced, or quality optimization profiles. Additional parameters control concurrent connection limits, timeout values, protocol selection, monitoring features, self-healing capabilities, and output directory configuration.

### Environment Variables

For containerized deployments, the system reads configuration from environment variables allowing for flexible deployment without code changes. Core settings include IRAN_MODE, DPI_EVASION_LEVEL, and PERFORMANCE_MODE. Performance tuning is controlled through MAX_CONCURRENT, TIMEOUT_SECONDS, and RETRY_ATTEMPTS. Feature flags enable or disable MONITORING, SELF_HEALING, FALLBACK, and CACHING capabilities.

### GitHub Actions Workflow Inputs

When triggering the workflow manually, users can specify check mode to control which components execute, performance mode for optimization strategy, DPI evasion level for filtering bypass aggressiveness, maximum concurrent connections, timeout values, protocol selection, and various feature flags for monitoring, self-healing, fallback systems, testing, security scanning, and artifact building.

## Monitoring and Metrics

### Quality Metrics

The system generates comprehensive quality metrics in JSON format including build identification and timestamp, version information, configuration parameters used during execution, success rates and performance indicators, retry attempt statistics, and feature enablement status. These metrics are stored in the metrics directory and can be used for analysis and optimization.

### Statistics Collection

Detailed statistics track proxy performance including total proxies tested, successful connection percentage, average response times, protocol distribution, DPI bypass success rates, and timestamp of last update. This information is available in the stats directory and is updated with each execution.

### Health Reporting

The GitHub Actions workflow generates comprehensive health reports in the workflow summary including build information and configuration parameters, execution status for each pipeline component, overall system health assessment, active features and capabilities, and alerts for critical failures with recommended actions.

## Advanced Features

### Self-Healing Capabilities

The system implements intelligent self-healing through automatic detection of proxy failures, dynamic adjustment of testing parameters based on network conditions, failover to backup proxies and alternative strategies, retry logic with exponential backoff, and recovery from transient network issues.

### Multi-Tier Fallback

When primary proxies fail, the system automatically falls back through multiple tiers including alternative proxy endpoints, different protocols when primary protocols are blocked, varied DPI evasion techniques, and emergency recovery mode for critical situations.

### Intelligent Caching

Performance is optimized through intelligent caching of Go modules and build artifacts with version-aware cache invalidation, dependency caching across workflow runs, and build artifact caching for faster subsequent builds.

### Circuit Breaker Pattern

The implementation includes circuit breaker functionality to prevent cascade failures by monitoring failure rates and automatically opening the circuit when thresholds are exceeded, allowing the system to recover before resuming normal operation, and preventing resource exhaustion during widespread failures.

## Security Considerations

### Code Security

The project maintains security through automated security scanning with gosec integrated into the CI/CD pipeline, upload of security scan results to GitHub Security tab for review, regular dependency updates to patch known vulnerabilities, and code quality checks including go vet and formatting verification.

### Container Security

Docker deployments follow security best practices including multi-stage builds to minimize attack surface, non-root user execution inside containers, minimal base images reducing potential vulnerabilities, health checks for monitoring container health, and proper secret management avoiding hardcoded credentials.

### Network Security

The system implements network security through encrypted proxy configurations, TLS certificate validation, secure communication protocols, and no retention of connection logs or user data.

## Performance Optimization

### Build Optimization

The build process is optimized through compilation with size and speed flags (-s -w), trimmed paths removing unnecessary debug information, cross-compilation support for multiple platforms, and parallel building when possible.

### Runtime Optimization

Performance during execution is enhanced through concurrent proxy checking with configurable limits, adaptive timeout adjustment based on network conditions, intelligent retry logic avoiding unnecessary attempts, and efficient memory usage through proper resource management.

### Caching Strategy

The caching implementation includes dependency caching to avoid redundant downloads, build artifact caching for faster rebuilds, version-aware cache invalidation ensuring freshness, and intelligent cache restoration with fallback keys.

## Troubleshooting

### Common Issues

When proxies are not functioning correctly, users should verify that Iran mode is enabled through the --iran-mode flag, consider increasing the DPI evasion level to maximum for more aggressive filtering bypass, check if specific ports are being blocked by their ISP, and enable fallback mode for additional redundancy.

For low success rates, adjusting the timeout parameter can help with slower connections, reducing concurrent connections may improve stability on resource-constrained systems, enabling self-healing allows automatic recovery from failures, and switching to quality performance mode prioritizes reliability over speed.

Build errors typically indicate issues with Go version compatibility (requiring 1.21 or higher), missing dependencies that need to be downloaded with go mod download, or missing system tools required for compilation.

### Debugging and Logging

Verbose logging can be enabled with the --verbose flag to provide detailed execution information. Metrics are available in JSON format in the metrics/quality-metrics.json file providing insights into system performance. Statistics can be reviewed in stats/proxy-stats.json showing historical trends and patterns. Workflow execution logs in GitHub Actions provide complete visibility into automated runs.

## Future Enhancements

The project roadmap includes plans for additional protocol support expanding beyond the current vmess, vless, trojan, and shadowsocks implementations. Enhanced machine learning capabilities will improve proxy quality prediction and selection. A web dashboard for real-time monitoring will provide better visibility into system operation. Mobile application support will extend accessibility to smartphones and tablets. Advanced analytics and reporting will offer deeper insights into performance trends. Community feedback integration will ensure the system evolves based on user needs and experiences.

## Contributing

Contributors are welcome and encouraged to participate in the project development. The process involves forking the repository, creating a feature branch for changes, implementing modifications with appropriate tests, running quality checks including linting and testing, committing changes with descriptive messages, pushing to the forked repository, and opening a pull request for review. All contributions should follow Go best practices, include tests for new functionality, update documentation as needed, and ensure all CI checks pass successfully.

## Support and Resources

For assistance with the project, users can open issues on GitHub for bugs or feature requests, review existing documentation in the docs directory for detailed information, examine closed issues for solutions to similar problems, and consult the README file for quick reference information.

## License

The project is licensed under the MIT License, allowing free use, modification, and distribution while maintaining attribution to the original authors. The complete license text is available in the LICENSE file in the repository root.

## Acknowledgments

The development of this system benefited from the work of the V2Ray and Xray communities in protocol implementation, contributors to Iran proxy filtering research who provided valuable insights, and the open source security tools and libraries that form the foundation of the security features.

## Conclusion

The Iran Proxy Ultimate System represents a comprehensive solution for proxy management and testing with specific optimizations for bypassing Iranian internet filtering. Through the unification of multiple workflows, modular architecture, extensive automation, and focus on reliability and security, the system provides a robust platform for maintaining access to unrestricted internet connectivity. The combination of advanced DPI bypass techniques, intelligent monitoring, self-healing capabilities, and flexible deployment options makes this an enterprise-grade solution suitable for both individual users and large-scale deployments.

---

**Project:** Iran Proxy Ultimate System  
**Version:** 3.2.0 - Unified Edition  
**Last Updated:** February 2024  
**Maintained by:** Iran Proxy Ultimate Team
