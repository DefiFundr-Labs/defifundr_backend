# Multi-stage build for DefiFundr backend
FROM golang:1.23-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o defifundr_backend ./cmd/api

# Development with hot reload
FROM builder AS development

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@v1.45.0

# Expose port
EXPOSE 8080

# Start with hot reload
CMD ["air"]

# Production image
FROM alpine:latest AS production

# Install necessary certificates
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/defifundr_backend .

# Copy config directory if it exists
COPY --from=builder /app/config/ /app/config/

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./defifundr_backend"]