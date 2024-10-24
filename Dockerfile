# Step 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files and download the necessary dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Step 2: Create a smaller image for running the application
FROM alpine:latest

# Set the working directory inside the new smaller image
WORKDIR /app

# Copy the binary from the builder stage to the new image
COPY --from=builder /app/main .

# Expose the necessary port for the REST API
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
