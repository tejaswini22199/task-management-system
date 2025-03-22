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

# Build the Go binary
RUN go build -o main ./cmd

# Use a minimal image for production
FROM ubuntu:22.04

WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Start the application
CMD ["./main"]
