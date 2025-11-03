# Gateway Service

API Gateway for Instrlabs platform providing reverse proxy functionality with service discovery, health monitoring, and request routing.

## Features

- Reverse proxy with path-based routing
- Service health monitoring and status checking
- Request forwarding with timeout handling
- CORS and security middleware
- Swagger API documentation
- Prometheus metrics integration
- Centralized logging

## Quick Start

```bash
# Local Setup
cp .env.example .env
go mod download
go run main.go

# Docker Setup
docker-compose up -d --build gateway-service
```

## Configuration

Refer to `.env.example` file for the list of required environment variables and their default values.

## Service Routing

The gateway service routes incoming requests to appropriate microservices based on path prefixes:

### Route Configuration

```
/auth/*    → auth-service (authentication & user management)
/images/*  → image-service (image processing & storage)
```

### Request Flow

```
Client Request → Gateway Service → Target Service → Response
```

**Example Request:**
```
GET /auth/login → Gateway → auth-service:3000/auth/login
POST /images/upload → Gateway → image-service:3000/images/upload
```

### Health Monitoring

**Gateway Health Check**
```
GET /health
Returns:
{
  "status": "ok",
  "services": {
    "auth-service": "ok",
    "image-service": "ok"
  }
}
```

**Service Health Status**
- Gateway periodically checks health of all registered services
- Returns service status in the health endpoint
- Routes requests only to healthy services
- Provides 502 Bad Gateway for unhealthy services

### Proxy Features

**Request Forwarding**
- Preserves HTTP method and headers
- Maintains query parameters
- 30-second timeout for all requests
- Automatic error handling for service unavailability

**Error Responses**
```json
{
  "error": "Bad Gateway",
  "message": "The service is currently unavailable"
}
```

```json
{
  "error": "Not Found",
  "message": "The requested resource does not exist"
}
```

## Project Structure

```
gateway-service/
├── main.go                    # Fiber app entry point
├── internal/
│   ├── config.go              # Config & service definitions
│   ├── router.go              # Route setup & proxy logic
│   ├── middleware.go          # CORS, security, logging
│   ├── swagger.go             # API documentation setup
│   ├── token.go               # JWT validation utilities
│   └── errors.go              # Error handling
├── static/                    # Static assets
└── Dockerfile
```

## Service Configuration

**Environment Variables**
```bash
# Gateway Configuration
PORT=:3000
ENVIRONMENT=development
ORIGINS_ALLOWED=http://localhost:8000
JWT_SECRET=your-jwt-secret
CSRF_ENABLED=true

# Service URLs
AUTH_SERVICE=http://auth-service:3000
IMAGE_SERVICE=http://image-service:3000
```

**Service Registration**
Services are automatically registered at startup and configured via environment variables. Each service requires:
- Name: Service identifier
- URL: Service endpoint
- Prefix: URL path prefix for routing

## Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [Fiber Proxy](https://github.com/gofiber/fiber/tree/master/middleware/proxy) - Reverse proxy middleware
- [Shared Init](github.com/instrlabs/shared/init) - Common initialization utilities