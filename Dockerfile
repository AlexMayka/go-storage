# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go mod download && \
    swag init -g cmd/api/main.go -o cmd/api/docs && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/migrations ./migrations

RUN mkdir -p /root/log

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 -O /dev/null http://localhost:8080/swagger/index.html || exit 1

CMD ["./main"]