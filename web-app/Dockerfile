FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o ipo-pilot-web .

# Use alpine for smaller image
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy binary and assets
COPY --from=builder /app/ipo-pilot-web .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose port
EXPOSE 8080

# Set environment
ENV GIN_MODE=release
ENV PORT=8080

# Run
CMD ["./ipo-pilot-web"]
