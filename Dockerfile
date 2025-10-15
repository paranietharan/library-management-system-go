FROM golang:1.25.1-alpine AS builder

# Enable Go modules and disable CGO for static build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Create and set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for efficient caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the binary
RUN go build -o library-system ./cmd/run/main.go

# -----------------------------
# Stage 2: Create minimal runtime image
# -----------------------------
FROM alpine:3.20

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/library-system .

# Copy configuration files if any (optional)
COPY .env .env

# Expose the app port (adjust if your app runs on another port)
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Run the binary
ENTRYPOINT ["./library-system"]
