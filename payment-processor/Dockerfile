# --------------------------
# Build Stage
# --------------------------
FROM golang:1.23-alpine AS builder

# Set working directory inside the builder
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/api

# --------------------------
# Runtime Stage
# --------------------------
FROM alpine:latest

# Set working directory inside the runtime container
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the .env file from builder stage
COPY --from=builder /app/cmd/api/.env /root/.env

# Expose port (adjust as needed)
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
