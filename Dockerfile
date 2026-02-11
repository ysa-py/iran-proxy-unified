# ═══════════════════════════════════════════════════════════════════════════
# Multi-Stage Dockerfile for Iran Proxy Ultimate System
# Optimized for: Security, Size, Performance
# ═══════════════════════════════════════════════════════════════════════════

# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Build arguments
ARG VERSION=3.2.0
ARG BUILD_ID=local
ARG BUILD_TIME

# Set working directory
WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Copy go mod files
COPY src/go.mod src/go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY src/ ./

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildID=${BUILD_ID} -X main.BuildTime=${BUILD_TIME}" \
    -trimpath \
    -o iran-proxy-ultimate \
    main.go main_iran.go

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 iran-proxy && \
    adduser -D -u 1000 -G iran-proxy iran-proxy

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/iran-proxy-ultimate /app/

# Copy configuration files
COPY configs/ /app/configs/

# Set permissions
RUN chown -R iran-proxy:iran-proxy /app && \
    chmod +x /app/iran-proxy-ultimate

# Switch to non-root user
USER iran-proxy

# Set timezone to Tehran
ENV TZ=Asia/Tehran

# Expose health check port (if applicable)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/iran-proxy-ultimate", "--health-check"] || exit 1

# Default command
ENTRYPOINT ["/app/iran-proxy-ultimate"]
CMD ["--iran-mode", "--performance-mode", "balanced", "--dpi-evasion-level", "aggressive"]

# Labels
LABEL org.opencontainers.image.title="Iran Proxy Ultimate"
LABEL org.opencontainers.image.description="Advanced Iran Proxy System with DPI Bypass"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.vendor="Iran Proxy Team"
LABEL org.opencontainers.image.licenses="MIT"
