# Use official Golang image
FROM golang:1.20

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire project
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 (or whichever your service runs on)
EXPOSE 8080

# Run the binary
CMD ["/app/main"]
