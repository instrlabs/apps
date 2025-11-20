# Design: Image Service Architecture Refactoring

## Architecture Overview

This design document outlines the technical approach for refactoring the image-service to follow the clean architecture patterns established in the auth-service.

## Current Architecture Analysis

### Problems with Current Structure

1. **Monolithic Handler**: `instruction_handler.go` contains 13,618 lines mixing:
   - HTTP request/response handling
   - Business logic for image processing
   - Database operations
   - File validation and processing
   - External service calls

2. **Flat Directory Structure**: All files in `internal/` without clear separation:
   ```
   internal/
   ├── config.go
   ├── image_service.go
   ├── instruction_detail_repository.go
   ├── instruction_handler.go (13,618 lines)
   ├── instruction_handler_test.go
   ├── instruction_processor.go
   ├── instruction_repository.go
   └── product_client.go
   ```

3. **Tight Coupling**: Direct dependencies make testing difficult:
   - Handlers directly access repositories
   - Business logic embedded in HTTP layer
   - Configuration scattered across files

## Target Architecture

### Layered Architecture Pattern

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (Handlers)                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ Instruction     │  │ Health          │  │ Auth         │ │
│  │ Handler         │  │ Handler         │  │ Middleware   │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   Business Layer (Services)                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ Image           │  │ Instruction     │  │ Processing   │ │
│  │ Service         │  │ Service         │  │ Service      │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                  Data Layer (Repositories)                  │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │ Instruction     │  │ Instruction     │                   │
│  │ Repository      │  │ Detail Repo     │                   │
│  └─────────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

### Directory Structure Design

```
image-service/
├── main.go                           # Application bootstrap
├── internal/
│   ├── config/
│   │   └── config.go                 # Configuration management
│   ├── models/                       # Domain models
│   │   ├── instruction.go
│   │   ├── image.go
│   │   └── processing_status.go
│   ├── handlers/                     # HTTP layer
│   │   ├── instruction_handler.go
│   │   └── health_handler.go
│   ├── services/                     # Business logic
│   │   ├── image_service.go
│   │   ├── instruction_service.go
│   │   └── processing_service.go
│   ├── repositories/                 # Data access
│   │   ├── instruction_repository.go
│   │   └── instruction_detail_repository.go
│   ├── validators/                   # Input validation
│   │   ├── request_validator.go
│   │   └── image_validator.go
│   ├── helpers/                      # Utilities
│   │   ├── response_helper.go
│   │   ├── s3_helper.go
│   │   └── nats_helper.go
│   ├── middleware/                   # HTTP middleware
│   │   └── auth_middleware.go
│   └── errors.go                     # Error definitions
├── pkg/
│   └── utils/                        # Reusable utilities
│       ├── mime.go
│   └── slice.go
└── static/
    └── swagger.json                  # API documentation
```

## Component Design

### 1. Configuration Layer

**Purpose**: Centralized configuration management with validation

**Key Features**:
- Environment variable loading with validation
- Default values for development
- Type-safe configuration structure
- Validation on startup

```go
type Config struct {
    Service         ServiceConfig
    Database        DatabaseConfig
    Storage         S3Config
    MessageBus      NatsConfig
    ExternalAPIs    ExternalAPIConfig
    Security        SecurityConfig
}
```

### 2. Domain Models

**Purpose**: Clear domain entities with validation

**Key Models**:
- `Instruction`: Core instruction entity
- `Image`: Image processing entity
- `ProcessingStatus`: Status enumeration
- `InstructionDetail`: Detailed processing information

### 3. Service Layer

**Image Service**: Core image processing operations
- Image validation and processing
- Format conversion and optimization
- Storage operations abstraction

**Instruction Service**: Instruction lifecycle management
- CRUD operations for instructions
- Status management
- Business rule enforcement

**Processing Service**: Background operations
- NATS message handling
- Asynchronous processing
- Cleanup operations

### 4. Repository Layer

**Purpose**: Data access abstraction
- MongoDB operations
- Query optimization
- Connection management
- Transaction support

### 5. Handler Layer

**Purpose**: Thin HTTP layer focusing on request/response
- Input validation
- Response formatting
- Error handling
- HTTP status code management

## Dependency Injection Design

### Constructor Pattern

```go
// Service construction with dependency injection
func NewInstructionService(
    repo InstructionRepository,
    detailRepo InstructionDetailRepository,
    imageService ImageService,
    config *Config,
) *InstructionService {
    return &InstructionService{
        repo:         repo,
        detailRepo:   detailRepo,
        imageService: imageService,
        config:       config,
    }
}

// Handler construction with injected services
func NewInstructionHandler(
    instructionService *InstructionService,
    imageService *ImageService,
    validator *RequestValidator,
    responseHelper *ResponseHelper,
) *InstructionHandler {
    return &InstructionHandler{
        instructionService: instructionService,
        imageService:       imageService,
        validator:          validator,
        responseHelper:     responseHelper,
    }
}
```

## Interface Design

### Repository Interfaces

```go
type InstructionRepository interface {
    Create(ctx context.Context, instruction *Instruction) error
    GetByID(ctx context.Context, id string) (*Instruction, error)
    Update(ctx context.Context, instruction *Instruction) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter InstructionFilter) ([]*Instruction, error)
}
```

### Service Interfaces

```go
type ImageService interface {
    ProcessImage(ctx context.Context, image *Image) (*ProcessedResult, error)
    ValidateImage(image *Image) error
    StoreImage(ctx context.Context, image *Image) (string, error)
    RetrieveImage(ctx context.Context, path string) (*Image, error)
}
```

## Error Handling Strategy

### Centralized Error Definitions

```go
type ServiceError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

var (
    ErrInvalidImageFormat = &ServiceError{
        Code:    "INVALID_IMAGE_FORMAT",
        Message: "Image format is not supported",
    }

    ErrInstructionNotFound = &ServiceError{
        Code:    "INSTRUCTION_NOT_FOUND",
        Message: "Instruction not found",
    }
)
```

### Response Format Standardization

```go
type APIResponse struct {
    Message string      `json:"message"`
    Errors  interface{} `json:"errors"`
    Data    interface{} `json:"data"`
}
```

## Validation Strategy

### Multi-Layer Validation

1. **Request Level**: HTTP input validation in handlers
2. **Domain Level**: Business rule validation in services
3. **Persistence Level**: Database constraints

### Validation Components

- `RequestValidator`: HTTP request validation
- `ImageValidator`: Image-specific validation
- `BusinessValidator`: Domain rule validation

## Testing Strategy

### Unit Testing Structure

```
tests/
├── services/
│   ├── image_service_test.go
│   ├── instruction_service_test.go
│   └── processing_service_test.go
├── handlers/
│   └── instruction_handler_test.go
├── repositories/
│   └── instruction_repository_test.go
└── validators/
    └── request_validator_test.go
```

### Mock Strategy

- Interface-based design enables easy mocking
- Dependency injection for test doubles
- Test utilities for common scenarios

## Migration Strategy

### Phase 1: Foundation (Day 1)
1. Create directory structure
2. Move existing files to appropriate locations
3. Set up configuration layer
4. Define basic interfaces and models
5. Ensure service builds and runs

### Phase 2: Service Extraction (Day 2)
1. Extract image processing logic to `ImageService`
2. Create `InstructionService` for instruction management
3. Move background processing to `ProcessingService`
4. Update handlers to use new services
5. Verify all functionality works

### Phase 3: Enhancement (Day 3)
1. Add comprehensive validation
2. Implement error handling standardization
3. Add response helpers
4. Update documentation
5. Add basic tests

## Performance Considerations

### Optimizations
- Maintain existing performance characteristics
- Database connection pooling
- Efficient image processing pipelines
- Background processing with NATS

### Monitoring
- Preserve existing Prometheus metrics
- Add service-specific metrics
- Error tracking and logging

## Security Considerations

### Maintain Existing Security
- JWT authentication integration
- Input validation and sanitization
- S3 access controls
- Rate limiting

### Enhancements
- Comprehensive input validation
- Error message sanitization
- Secure file upload handling

This design provides a clean, maintainable architecture that aligns with the auth-service patterns while preserving all existing functionality and performance characteristics.