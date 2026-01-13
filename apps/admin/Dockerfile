# Build Stage 1: Frontend (Vite/React)
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm run build

# Build Stage 2: Backend (Go)
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o admin-service ./apps/admin

# Final Stage: Run
FROM alpine:latest
WORKDIR /app

# Copy binary from backend builder
COPY --from=backend-builder /app/admin-service .

# Copy static assets from frontend builder
COPY --from=frontend-builder /app/dist ./public

# Configuration
ENV PORT=8080
ENV STATIC_DIR=./public
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./admin-service"]
