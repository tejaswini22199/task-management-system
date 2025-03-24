# Use official Golang image
FROM golang:1.24 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build multiple services separately
RUN go build -o main ./cmd
RUN go build -o auth ./cmd/authservice
RUN go build -o tasks ./cmd/taskservice

# Use a minimal image for production
FROM ubuntu:22.04

WORKDIR /root/

# Copy the compiled binaries from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/auth .
COPY --from=builder /app/tasks .

# Expose the ports
EXPOSE 8080 8000 8001

# Default command (can be overridden in docker-compose)
CMD ["./main"]
