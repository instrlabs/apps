# Design: Product Service Architecture Refactoring

## Context

The product-service currently has a basic structure that mixes concerns and doesn't follow the clean architecture patterns established in the auth-service. The existing code has handlers that contain business logic, no clear service layer, and inconsistent organization patterns.

**Current Issues:**
- Business logic mixed with HTTP handling in `product_handler.go`
- No service layer for business rule separation
- Inconsistent directory structure compared to auth-service
- Limited testing support due to tight coupling
- Missing input validation layer
- Configuration management is minimal

**Goals:**
- Align with auth-service architecture patterns
- Enable clean testing with proper separation
- Improve code organization and maintainability
- Maintain full backward compatibility
- Follow established microservice conventions

## Decisions

### Decision 1: Feature-Based Organization (Auth-Service Pattern)

**What:** Reorganize code to match auth-service structure:

```
internal/
├── config/           # Configuration with validation
├── models/           # Domain models with business logic
├── handlers/         # Thin HTTP handlers
├── services/         # Business logic services
├── repositories/     # Data access layer
├── validators/       # Input validation
└── helpers/          # Utility functions
```

**Why:**
- Consistency across microservices
- Easier for developers to switch between services
- Clear separation of concerns
- Follows proven patterns from auth-service

### Decision 2: Service Layer for Business Logic

**What:** Extract business logic from handlers into dedicated service:

```go
type ProductService interface {
    ListProducts(productType string, page, limit int) (*ProductListResult, error)
    GetProductByID(id string, productType string) (*Product, error)
}

type productService struct {
    repo repositories.ProductRepositoryInterface
}

func (s *productService) ListProducts(productType string, page, limit int) (*ProductListResult, error) {
    // Business logic for listing products with pagination
    products, err := s.repo.List(productType)
    if err != nil {
        return nil, err
    }

    // Apply pagination logic
    return &ProductListResult{
        Products: paginateProducts(products, page, limit),
        Pagination: calculatePagination(len(products), page, limit),
    }, nil
}
```

**Why:**
- Separates HTTP concerns from business rules
- Makes business logic reusable and testable
- Enables clean handler implementations
- Follows single responsibility principle

### Decision 3: Enhanced Configuration Management

**What:** Implement comprehensive configuration similar to auth-service:

```go
type Config struct {
    ServiceName     string        `env:"SERVICE_NAME,required"`
    Port            string        `env:"PORT,default=3002"`
    Environment     string        `env:"ENVIRONMENT,default=development"`
    MongoURI        string        `env:"MONGO_URI,required"`
    MongoDB         string        `env:"MONGO_DB,required"`
    MongoTimeout    int           `env:"MONGO_TIMEOUT,default=10"`
    ReadTimeout     int           `env:"READ_TIMEOUT,default=30"`
    WriteTimeout    int           `env:"WRITE_TIMEOUT,default=30"`
    IdleTimeout     int           `env:"IDLE_TIMEOUT,default=60"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return cfg, nil
}
```

**Why:**
- Consistent configuration across services
- Proper validation at startup
- Environment-specific configurations
- Clear error messages for misconfiguration

### Decision 4: Centralized Input Validation

**What:** Add dedicated validation layer:

```go
type RequestValidator struct {
    validator *validator.Validate
}

func (v *RequestValidator) ValidateProductListRequest(c *fiber.Ctx) (*ProductListRequest, error) {
    req := &ProductListRequest{
        Type:  c.Query("type", ""),
        Page:  c.QueryInt("page", 1),
        Limit: c.QueryInt("limit", 50),
    }

    if err := v.validatePagination(req.Page, req.Limit); err != nil {
        return nil, err
    }

    return req, nil
}
```

**Why:**
- Clean separation of validation logic
- Reusable validation rules
- Consistent error messages
- Easier testing of validation rules

### Decision 5: Improved Main.go Structure

**What:** Refactor main.go to follow auth-service patterns with graceful shutdown:

```go
func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize database
    client, db := initx.NewMongo()
    defer initx.CloseMongo(client)

    // Create Fiber app
    app := setupFiberApp(cfg)

    // Setup routes
    setupRoutes(app, db, cfg)

    // Start server
    startServer(app, cfg)
}
```

**Why:**
- Consistent application startup across services
- Proper graceful shutdown handling
- Clean separation of concerns
- Better error handling and logging

## Target Directory Structure

```
product-service/
├── main.go                          # Clean setup with constructor injection
├── go.mod                           # Dependencies
├── go.sum                           # Dependency lock
├── Dockerfile                       # Multi-stage build
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration with validation
│   ├── models/
│   │   └── product.go              # Product domain model with business logic
│   ├── handlers/
│   │   └── product_handler.go      # Thin HTTP handlers focused on HTTP concerns
│   ├── services/
│   │   └── product_service.go      # Business logic layer
│   ├── repositories/
│   │   └── product_repository.go   # Data access layer with clean interface
│   ├── validators/
│   │   └── request_validator.go    # Input validation logic
│   └── helpers/
│       └── response_helper.go      # HTTP response utilities
├── static/
│   └── swagger.json               # API documentation
└── tests/                         # Integration tests
```

## Migration Strategy

### Phase 1: Directory Structure (Day 1)
1. Create new directory structure matching auth-service
2. Move existing files to appropriate directories
3. Update imports throughout the codebase
4. Ensure service builds successfully

### Phase 2: Service Layer Extraction (Day 2)
1. Create service interfaces and implementations
2. Extract business logic from handlers into services
3. Add proper error handling and validation
4. Test service layer independently

### Phase 3: Handler Refactoring (Day 2-3)
1. Refactor handlers to be thin HTTP layers
2. Add input validation layer
3. Implement consistent response formatting
4. Add comprehensive error handling

### Phase 4: Configuration and Main (Day 3)
1. Enhance configuration management
2. Refactor main.go for consistent patterns
3. Add graceful shutdown
4. Update Docker and deployment configurations

## Benefits Over Current Structure

**Consistency:**
- Matches auth-service architecture patterns
- Easier team mobility between services
- Shared conventions and practices

**Maintainability:**
- Clear separation of concerns
- Smaller, focused files
- Single responsibility principle

**Testability:**
- Service layer enables unit testing
- Clean interfaces for mocking
- Separated validation logic

**Scalability:**
- Easy to add new product-related features
- Clear patterns for extending functionality
- Consistent with microservice architecture

## When to Evolve

This refactoring provides a solid foundation that matches the current service complexity. Consider evolution only when:

- Service grows significantly (>20 endpoints)
- Complex business rules emerge (pricing, inventory, etc.)
- Multiple data sources are integrated
- Advanced caching or performance optimization is needed

For now, focus on clean, maintainable code that follows established patterns.