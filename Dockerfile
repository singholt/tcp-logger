# Build stage
FROM golang:latest AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY logger.go .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o logger .

# Run stage
FROM amazonlinux:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/logger .

# Command to run the executable
ENTRYPOINT ["./logger"]
