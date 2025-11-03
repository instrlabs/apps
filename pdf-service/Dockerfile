# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pdf-service .

# Runtime stage
FROM alpine:latest AS runtime

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata dumb-init

# Create non-root user
RUN addgroup -g 1001 -S instrlabs && \
    adduser -u 1001 -S instrlabs -G instrlabs

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/pdf-service .

# Copy static files
COPY --from=builder /app/static ./static

# Change ownership to non-root user
RUN chown -R instrlabs:instrlabs /app

# Switch to non-root user
USER instrlabs

# Expose port
EXPOSE 3004

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3004/health || exit 1

# Run the application
ENTRYPOINT ["dumb-init", "./pdf-service"]