# Use Debian-based Go image (full glibc support for SQLite)
FROM golang:1.21 AS builder

WORKDIR /app/web-app

# Copy go mod files from web-app directory
COPY web-app/go.mod web-app/go.sum ./
RUN go mod download

# Copy source code from web-app
COPY web-app/ .

# Build the application with CGO enabled (required for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -o ipo-pilot .

# Use Debian slim for smaller final image
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copy binary and assets from builder
COPY --from=builder /app/web-app/ipo-pilot .
COPY --from=builder /app/web-app/templates ./templates
COPY --from=builder /app/web-app/static ./static

# Expose port
EXPOSE 8080

# Set environment
ENV GIN_MODE=release
ENV PORT=8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run
CMD ["./ipo-pilot"]
