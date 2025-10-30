# Multi-stage build for Go backend
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/page-insight-tool ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/page-insight-tool .
COPY --from=builder /app/config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./page-insight-tool"]

