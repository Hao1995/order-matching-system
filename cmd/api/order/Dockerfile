# Use the official Golang 1.22.1 image as a base
FROM golang:1.23.5-alpine

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules manifest and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY ./pkg ./pkg
COPY ./internal/common ./internal/common
COPY ./cmd/api/order ./cmd/api/order
COPY ./internal/api/order ./internal/api/order

# Build the Go application
RUN go build -o main ./cmd/api/order

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
