<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# InstrLabs Backend Development Standards

## Microservice Architecture

The InstrLabs backend is organized as a Go microservice architecture using Fiber framework:

### Core Services
- **auth-service**: Authentication and session management with device binding
- **gateway-service**: API gateway with reverse proxy and service discovery
- **image-service**: Image processing pipeline with S3 storage integration
- **notification-service**: Real-time notifications via Server-Sent Events (SSE)

## Microservice Development Standards

### Go Microservices (InstrLabs Service Pattern)

#### Standard Project Structure
```
service-name/
├── main.go                 # Entry point
├── go.mod                  # Dependencies
├── go.sum                  # Dependency lock
├── Dockerfile              # Multi-stage build
├── .dockerignore           # Build optimization
├── internal/               # Private app code
│   ├── config.go           # Environment & configuration
│   ├── middleware.go       # HTTP middleware setup
│   ├── router.go           # Route definitions
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Data structures
│   ├── services/           # Business logic
│   ├── repositories/       # Data access layer
│   └── errors.go           # Error definitions
├── pkg/                    # Public reusable code
├── static/                 # Static assets (Swagger, etc.)
└── scripts/                # Build/deployment scripts
```

#### Efficient Initialization Pattern

```go
// main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/log"
)

func main() {
    // 1. Load configuration
    cfg, err := LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. Initialize app with optimizations
    app := fiber.New(fiber.Config{
        ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
        WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
        IdleTimeout:     time.Duration(cfg.IdleTimeout) * time.Second,
        CaseSensitive:   true,
        StrictRouting:   true,
        ServerHeader:    cfg.ServiceName,
        AppName:         cfg.ServiceName,
        DisableKeepalive: false,
        ReadBufferSize:  4096,
        WriteBufferSize: 4096,
    })

    // 3. Setup middleware in correct order
    SetupMiddleware(app, cfg)

    // 4. Setup routes
    SetupRoutes(app, cfg)

    // 5. Graceful shutdown
    go func() {
        if err := app.Listen(":" + cfg.Port); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    gracefulShutdown(app)
}

func gracefulShutdown(app *fiber.App) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := app.ShutdownWithContext(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Info("Server exited")
}
```

#### Configuration Management

```go
// internal/config.go
package internal

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    // Service
    ServiceName string `env:"SERVICE_NAME,required"`
    Port        string `env:"PORT,default=3000"`
    Environment string `env:"ENVIRONMENT,default=development"`

    // Database
    MongoURI   string `env:"MONGO_URI,required"`
    MongoDB    string `env:"MONGO_DB,required"`
    MongoTimeout int `env:"MONGO_TIMEOUT,default=10"`

    // Security
    JWTSecret        string `env:"JWT_SECRET,required"`
    TokenExpiryHours int    `env:"TOKEN_EXPIRY_HOURS,default=1"`

    // CORS
    Origins        string `env:"CORS_ORIGINS,default=http://localhost:3000"`
    CSRFEnabled    bool   `env:"CSRF_ENABLED,default=true"`

    // Rate limiting
    RateLimit    int           `env:"RATE_LIMIT,default=100"`
    RateWindow   time.Duration `env:"RATE_WINDOW,default=60s"`

    // Timeouts
    ReadTimeout  int `env:"READ_TIMEOUT,default=30"`
    WriteTimeout int `env:"WRITE_TIMEOUT,default=30"`
    IdleTimeout  int `env:"IDLE_TIMEOUT,default=60"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    // Use environment variable library like godotenv for local dev
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}
```

#### Middleware Stack Pattern

```go
// internal/middleware.go
package internal

import (
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/etag"
    "github.com/gofiber/fiber/v2/middleware/helmet"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func SetupMiddleware(app *fiber.App, cfg *Config) {
    // 1. Security first
    app.Use(helmet.New())

    // 2. Request tracking
    app.Use(requestid.New())

    // 3. Recovery
    app.Use(recover.New(recover.Config{
        EnableStackTrace: cfg.Environment == "development",
    }))

    // 4. Compression
    app.Use(compress.New(compress.Config{
        Level: compress.LevelBestSpeed,
    }))

    // 5. ETag for caching
    app.Use(etag.New())

    // 6. Rate limiting
    app.Use(limiter.New(limiter.Config{
        Max:        cfg.RateLimit,
        Expiration: cfg.RateWindow,
        SkipFailedRequests: false,
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "message": "Rate limit exceeded",
                "errors":  nil,
                "data":    nil,
            })
        },
    }))

    // 7. CORS
    app.Use(cors.New(cors.Config{
        AllowOrigins:     cfg.Origins,
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "content-type,cookie,authorization",
        AllowCredentials: true,
        MaxAge:           86400, // 24 hours
    }))

    // 8. Custom middleware
    setupCustomMiddleware(app, cfg)
}
```

#### Database Connection Pattern

```go
// internal/database.go
package internal

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
    Client   *mongo.Client
    Database *mongo.Database
}

func NewDatabase(cfg *Config) (*Database, error) {
    ctx, cancel := context.WithTimeout(context.Background(),
        time.Duration(cfg.MongoTimeout)*time.Second)
    defer cancel()

    clientOpts := options.Client().
        ApplyURI(cfg.MongoURI).
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(30 * time.Second)

    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        return nil, err
    }

    // Ping to verify connection
    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }

    return &Database{
        Client:   client,
        Database: client.Database(cfg.MongoDB),
    }, nil
}

func (d *Database) Close() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return d.Client.Disconnect(ctx)
}
```

### Fiber Framework Patterns

#### Route Registration Pattern

```go
// internal/router.go
package internal

import (
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *Config) {
    // API v1 routes
    v1 := app.Group("/api/v1")

    // Authentication routes
    auth := v1.Group("/auth")
    auth.Post("/login", authHandler.Login)
    auth.Post("/register", authHandler.Register)
    auth.Post("/logout", authHandler.Logout)

    // Protected routes
    protected := v1.Group("/")
    protected.Use(authMiddleware())

    // User routes
    user := protected.Group("/users")
    user.Get("/", userHandler.GetProfile)
    user.Put("/", userHandler.UpdateProfile)

    // Health check
    app.Get("/health", HealthCheck)
}
```

#### Proxy Configuration Pattern (Gateway Service)

```go
// internal/proxy.go
package internal

import (
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/proxy"
)

func SetupProxyRoutes(app *fiber.App, cfg *Config) {
    // Auth service proxy
    app.All("/auth/*", proxy.Forward(cfg.AuthService, proxy.Config{
        Timeout: 30 * time.Second,
        ModifyRequest: func(c *fiber.Ctx) error {
            c.Request().Header.Set("X-Forwarded-Host", c.Hostname())
            return nil
        },
    }))

    // Image service proxy
    app.All("/images/*", proxy.Forward(cfg.ImageService, proxy.Config{
        Timeout: 60 * time.Second, // Longer timeout for file uploads
        ModifyRequest: func(c *fiber.Ctx) error {
            c.Request().Header.Set("X-Forwarded-Host", c.Hostname())
            return nil
        },
    }))
}
```

## Common Patterns

### Environment Variables

Always use `.env.example` as a template:

```bash
# .env.example
# Service Configuration
SERVICE_NAME=auth-service
NODE_ENV=development
PORT=3001

# Database
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=auth_service

# Security
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=1h

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# Rate Limiting
RATE_LIMIT=100
RATE_WINDOW_MS=60000

# External Services
API_GATEWAY_URL=http://localhost:3000
```

### Docker Patterns

#### Go Dockerfile (Multi-stage)
```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Runtime stage
FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates dumb-init
WORKDIR /app
COPY --from=builder /go/src/app/app .
EXPOSE 3000
CMD ["dumb-init", "./app"]
```

#### Go Dockerfile (Production)

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Runtime stage
FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates tzdata dumb-init
WORKDIR /app
COPY --from=builder /go/src/app/app .
COPY --from=builder /go/src/app/static ./static
EXPOSE 3000
CMD ["dumb-init", "./app"]
```

### Health Checks

Implement consistent health check endpoints:

```go
// Go health check
func HealthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "ok",
        "timestamp": time.Now().UTC(),
        "service": cfg.ServiceName,
        "version": os.Getenv("VERSION"),
    })
}
```

### Standardized Response Format

All Go services should use consistent response format:

```json
{
  "message": "Operation successful",
  "errors": null,
  "data": { ... }
}
```

### Logging Standards

All services use Fiber's built-in logger for simple message-based logging:

```go
import "github.com/gofiber/fiber/v2/log"
```

#### Available Log Methods

```go
// Simple info messages
log.Info("Server started successfully")

// Formatted info with values
log.Infof("User logged in successfully: %s", email)

// Warnings
log.Warn("Invalid token provided")
log.Warnf("Rate limit exceeded for IP: %s", ipAddress)

// Errors
log.Error("Failed to connect to database")
log.Errorf("Failed to create session: %v", err)

// Fatal errors (exits the application)
log.Fatal("Unable to start server")
log.Fatalf("Failed to load config: %v", err)
```

#### Real Examples from Services

```go
// Auth Service - function context prefix pattern
log.Info("Login: Processing login request")
log.Infof("Login: Attempting to login user with email: %s", input.Email)
log.Errorf("Login: Failed to create session: %v", err)
log.Warnf("Login: Invalid request body: %v", err)

// PDF Service - operation logging
log.Errorf("Failed to compress PDF: %v", err)
log.Infof("PDF compressed successfully: %d bytes -> %d bytes", originalSize, compressedSize)

// Notification Service - state tracking
log.Infof("New SSE client connected for user %s. Total clients: %d", userId, totalClients)
log.Infof("SSE client disconnected for user %s", userId)
```

#### HTTP Request Logging Middleware

HTTP requests are automatically logged in JSON format via middleware:

```go
// Setup in main.go
middlewarex.SetupLogger(app)
```

This logs all HTTP requests with timestamp, method, path, status, latency, and user agent.

## API Documentation Standards

### Swagger/OpenAPI Documentation

All services must include comprehensive Swagger/OpenAPI 3.0.3 documentation in `static/swagger.json`:

**Required Documentation Elements:**
- API metadata (title, version, description)
- Tag organization for endpoints
- Complete path definitions with HTTP methods
- Request/response schemas for all endpoints
- Error responses with status codes
- Authentication requirements
- Parameter descriptions and examples
- Request body examples

**Response Format Standards:**
- Consistent JSON structure across all services
- Standardized error responses
- Success response format with message/errors/data pattern

### Service-Specific Patterns

#### Authentication Service
- PIN-based and OAuth authentication flows
- Session management with device binding
- JWT token generation and validation
- Multi-device session support

#### Gateway Service
- Reverse proxy with path-based routing
- Service health monitoring
- Request forwarding and load balancing
- CORS and security middleware

#### Image Service
- File upload and processing pipeline
- S3 storage integration
- Instruction management system
- Real-time status tracking

#### Notification Service
- Server-Sent Events (SSE) implementation
- NATS message bus integration
- User-specific notification streams
- Connection management and cleanup

## Repository Structure

### Backend Repository (Current)
This repository contains all microservices with:
- Shared configurations and standards
- Docker Compose for local development
- Common deployment patterns
- Centralized documentation

This documentation provides efficient, production-ready initialization patterns for Go microservices based on the observed patterns across all InstrLabs services.