# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Install swag for Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Copy source code first
COPY . .

# Generate Swagger documentation and build
RUN go mod tidy && \
    go mod download && \
    swag init -g cmd/api/main.go -o cmd/api/docs && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates and wget for health checks
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migration files if they exist
COPY --from=builder /app/migrations ./migrations

# Create log directory
RUN mkdir -p /root/log

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 -O /dev/null http://localhost:8080/swagger/index.html || exit 1

# Command to run
CMD ["./main"]