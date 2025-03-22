FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Create proto directory
RUN mkdir -p proto

# Copy proto files
COPY proto/generate.sh proto/
COPY ../proto/ecommerce.proto ../proto/

# Install protobuf compiler
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy source code
COPY . .

# Generate proto files
RUN mkdir -p proto/ecommerce && chmod +x proto/generate.sh && proto/generate.sh

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cart-service .

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cart-service .

# Create secrets directory
RUN mkdir -p /app/secrets

# Set non-root user for security
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser && \
    chown -R appuser:appgroup /app

USER appuser

# Command to run
CMD ["./cart-service"] 