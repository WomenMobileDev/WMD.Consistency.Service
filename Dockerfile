# Build stage
FROM golang:1.24.3-alpine AS builder

# Set necessary environment variables
ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  GO111MODULE=on

WORKDIR /app

# Install git and certificates for private repos if needed
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version info
RUN go build -ldflags="-s -w" -o main ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Add necessary runtime packages
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create a non-root user to run the application
RUN adduser -D -g '' appuser

# Copy binary and config from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

# Set proper permissions
RUN chown -R appuser:appuser /app

# Use the non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
