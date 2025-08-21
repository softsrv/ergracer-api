# Build stage
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with security hardening (pure Go, no CGO needed)
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o main .

# Runtime stage
FROM alpine:3.21

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Install CA certificates for HTTPS requests and update packages
RUN apk --no-cache add ca-certificates wget && \
    apk --no-cache upgrade

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]