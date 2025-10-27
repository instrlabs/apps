# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based application with a Next.js frontend and four Go backend services communicating via NATS message broker. All services use Fiber web framework and share common infrastructure through `github.com/instrlabs/shared`.

### Services

**Gateway Service** (`gateway-service/`) - Port 3000
- API gateway that proxies requests to downstream services
- Handles CORS, rate limiting, CSRF protection, and JWT authentication
- Extracts JWT tokens from cookies and injects user metadata headers (`x-authenticated`, `x-user-id`, `x-user-roles`)
- Routes are dynamically configured via `Services` array in config

**Auth Service** (`auth-service/`) - Port 3001
- User authentication with email/PIN and Google OAuth
- JWT token generation and refresh
- MongoDB for user storage
- Unauthenticated endpoints: `/login`, `/refresh`, `/send-pin`, `/google`, `/google/callback`

**Image Service** (`image-service/`) - Port 3002
- Image processing instruction management
- S3 storage integration for input/output files
- NATS subscriber for async image processing (`image.requests` subject)
- Background cleanup job runs every 30 minutes
- MongoDB for instruction/file metadata

**Notification Service** (`notification-service/`) - Port 3003
- Server-Sent Events (SSE) for real-time notifications
- NATS subscriber for notification messages (`notifications.sse` subject)
- All endpoints require authentication

**Web Application** (`web/`) - Port 8000
- Next.js 15 with App Router and React 19
- See `web/CLAUDE.md` for detailed frontend guidance

## Development Commands

### Backend Services

**Run a service locally:**
```bash
cd <service-name>
cp .env.staging.example .env.staging  # Configure first
go run main.go
```

**Build a service:**
```bash
cd <service-name>
go build -o app .
```

**Run with Docker Compose:**
```bash
docker-compose up --build
# Gateway: http://localhost:3000
# Auth: http://localhost:3001
# Image: http://localhost:3002
# Notification: http://localhost:3003
```

**Note:** The web service is commented out in `docker-compose.yaml` and should be run separately.

### Frontend

**Development:**
```bash
cd web
npm install
npm run dev  # Starts on port 8000
```

**Production:**
```bash
cd web
npm run build
npm start
```

## Environment Configuration

Each service requires its own `.env` file. See `.env.example` in the root and each service directory.

**Key shared variables:**
- `JWT_SECRET` - Must be identical across gateway and auth services
- `MONGO_URI`, `MONGO_DB` - MongoDB connection
- `NATS_URI` - NATS message broker
- `S3_*` - S3/MinIO configuration (image-service only)
- Service URLs for inter-service communication

**Frontend-specific:**
- `GATEWAY_URL` - Backend API endpoint
- `NOTIFICATION_URL` - SSE service (server-side only, proxied via `/api/sse`)

## Shared Infrastructure Module

All services depend on `github.com/instrlabs/shared` (v0.0.10+) which provides:

**Initialization helpers** (`initx` package):
- `NewMongo()` - MongoDB client setup
- `NewS3()` - MinIO/S3 client setup
- `NewNats()` - NATS connection
- `SetupPrometheus()` - Metrics endpoint at `/metrics`
- `SetupServiceHealth()` - Health check at `/health`
- `SetupLogger()` - Request logging middleware
- `SetupServiceSwagger()` - Swagger UI at `/swagger/*`
- `SetupAuthenticated()` - Authentication middleware with exclusion list

**Common patterns:**
```go
cfg := internal.LoadConfig()
mongo := initx.NewMongo(&initx.MongoConfig{
    MongoURI: cfg.MongoURI,
    MongoDB: cfg.MongoDB,
})
defer mongo.Close()
```

## Authentication Flow

1. **Client → Gateway**: Request with cookies (`access_token`, `refresh_token`)
2. **Gateway middleware** (`gateway-service/internal/middleware.go:71-105`):
   - Extracts JWT from cookies
   - Validates and decodes token using `JWT_SECRET`
   - Strips `cookie` header and injects:
     - `x-authenticated: true`
     - `x-user-id: <id>`
     - `x-user-roles: <comma-separated>`
     - `x-gateway: true`
   - Returns 401 if token is invalid (but not if missing)
3. **Service middleware** (`SetupAuthenticated`):
   - Checks `x-authenticated` header
   - Returns 401 if false and endpoint not in exclusion list
4. **Service handler**: Accesses user info from headers

**Refresh flow:**
- If `access_token` missing but `refresh_token` present, frontend middleware calls `/auth/refresh`
- Gateway passes refresh token via `x-user-refresh` header on `/auth/refresh` endpoint

## Message Bus (NATS)

**Subjects:**
- `image.requests` - Image processing tasks (published by image-service HTTP handlers, consumed by image-service worker)
- `notifications.sse` - Real-time notifications (consumed by notification-service for SSE broadcasting)

**Subscription pattern:**
```go
nats := initx.NewNats(cfg.NatsURI)
defer nats.Close()

nats.Conn.Subscribe(cfg.NatsSubjectImageRequests, func(m *natsgo.Msg) {
    handler.ProcessMessage(m.Data)
})
```

## MongoDB Collections

Services use MongoDB for persistence with the following patterns:

**Auth Service:**
- `users` - User accounts with password hashing (bcrypt)

**Image Service:**
- `products` - Available image processing products
- `instructions` - Image processing jobs
- `files` - Input/output file metadata (actual files in S3)

Repository pattern with methods like `FindByID`, `Create`, `Update`, `Delete`.

## API Response Format

All API responses follow this structure:
```json
{
  "success": true,
  "message": "Success message",
  "data": { ... },
  "errors": null
}
```

**Frontend fetch utilities** (`web/utils/fetch.ts`) automatically handle this format.

## Docker Build Pattern

All Go services use multi-stage builds:
1. **Build stage**: `golang:1.24` with `CGO_ENABLED=0` for static binaries
2. **Release stage**: `alpine:latest` with `dumb-init` for proper signal handling
3. Services expose port `3000` internally (mapped to different ports in docker-compose)

## CI/CD

**GitHub Actions** (`.github/workflows/image-builder.yml`):
- Manual workflow dispatch for building and pushing service images
- Builds for `linux/amd64` platform
- Tags: `<dockerhub-user>/instrlabs-<service>:<sha>-<environment>`
- Supports build cache via GitHub Actions cache

**To trigger a build:**
Go to Actions → image-builder → Run workflow → Select service and environment

## Code Patterns

### Service Structure
```
<service>/
├── Dockerfile
├── go.mod
├── main.go          # Entry point with setup
├── internal/
│   ├── config.go    # Environment config loader
│   ├── *_handler.go # HTTP handlers (Fiber)
│   ├── *_repository.go # Database access
│   └── *.go         # Domain models
└── static/          # Static assets (if any)
```

### Handler Methods
```go
func (h *Handler) MethodName(c *fiber.Ctx) error {
    // Parse request
    var req RequestType
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Invalid request",
            "errors": err.Error(),
            "data": nil,
        })
    }

    // Business logic

    // Return response
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Success",
        "data": result,
        "errors": nil,
    })
}
```

### Adding New Routes
1. Define handler method in `internal/*_handler.go`
2. Register route in `main.go` (services) or `internal/router.go` (gateway)
3. Add to authenticated exclusion list if public endpoint

### Gateway Route Configuration
Routes are proxied based on config in `gateway-service/internal/router.go:35-63`:
```go
for _, srv := range config.Services {
    app.All(srv.Prefix+"/*", ...)
}
```

Service prefixes typically match: `/auth/*` → auth-service, `/images/*` → image-service, etc.

## Testing

No test infrastructure is currently present in the codebase. When adding tests:
- Use `testing` package for Go services
- Place test files alongside source: `*_test.go`
- Run with `go test ./...`
- For frontend, see `web/CLAUDE.md`

## Common Development Tasks

**Adding a new microservice:**
1. Copy structure from existing service (e.g., `auth-service`)
2. Update `go.mod` module name
3. Add to `docker-compose.yaml`
4. Register routes in `gateway-service/internal/config.go` `Services` array
5. Create `.env` file from `.env.example`

**Adding new NATS subjects:**
1. Define in `.env.example` and all service `.env` files
2. Add to config structs in `internal/config.go`
3. Publish/subscribe using `nats.Conn.Publish()` / `Subscribe()`

**Updating shared module:**
1. Make changes in `github.com/instrlabs/shared` repository
2. Update version in each service: `go get github.com/instrlabs/shared@v0.0.X`
3. Run `go mod tidy`