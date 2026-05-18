# Start from the official Golang image
FROM golang:1.21-alpine as builder

WORKDIR /app

# Install git (required for go mod)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the backend binary
RUN go build -o server ./cmd/main.go

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the built binary from builder
COPY --from=builder /app/server .

# Expose the port your app runs on (change if needed)
EXPOSE 8080

# Set environment variables (override in Render dashboard as needed)
ENV PORT=8080

# Command to run the server
CMD ["./server"]
