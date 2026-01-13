# Build Stage: Backend (Go)
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o scheduler ./apps/scheduler

# Final Stage: Run
FROM alpine:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/scheduler .

# Configuration
ENV ENABLED_JOBS=""
ENV PORT=8080

EXPOSE 8080

CMD ["./scheduler"]
