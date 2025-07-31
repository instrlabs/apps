FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev pkgconfig imagemagick-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o app .

# Create final image
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache imagemagick

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Run the application
CMD ["./app"]