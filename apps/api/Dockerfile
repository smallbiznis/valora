# Build Stage: Backend (Go)
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api-service ./apps/api

# Final Stage: Run
FROM alpine:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api-service .

# Configuration
ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./api-service"]
