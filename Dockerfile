FROM golang:1.21-alpine AS builder

# Install build dependencies for CGO (required for SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app/web-app

# Copy go mod files from web-app directory
COPY web-app/go.mod web-app/go.sum ./
RUN go mod download

# Copy source code from web-app
COPY web-app/ .

# Build the application with CGO enabled (required for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o ipo-pilot .

# Use alpine for smaller image
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

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
