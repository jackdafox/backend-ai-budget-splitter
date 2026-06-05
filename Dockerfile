# syntax=docker/dockerfile:1.6

# Build stage
FROM golang:1.22-alpine AS builder

# Install git for go mod operations
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go module files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# - CGO_ENABLED=0 produces a static binary
# - -trimpath removes file system paths from binary
# - -ldflags "-s -w" strips debug info (smaller binary)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /build/server \
    ./cmd/server

# Runtime stage
FROM alpine:3.19

# Install ca-certificates for HTTPS calls to AI providers
RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary and config from builder
COPY --from=builder /build/server /app/server
COPY config.yaml /app/config.yaml

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set default environment variables
ENV APP_HOST=0.0.0.0 \
    APP_PORT=8080 \
    APP_ENV=production

# Run the server
ENTRYPOINT ["/app/server"]
