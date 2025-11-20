# Proposal: Refactor Image Service Architecture

## Summary

Refactor the image-service to follow the same clean, feature-based architecture patterns established in the auth-service. The current image-service has a basic structure but lacks proper separation of concerns, testing support, and consistent patterns with the rest of the microservices.

**Proposed Change ID:** `refactor-image-service-architecture`

## Current State Analysis

The image-service currently has:

- **Flat internal structure**: All files are in a single `internal/` directory without clear separation of concerns
- **Mixed responsibilities**: `instruction_handler.go` contains business logic, validation, HTTP concerns, and image processing logic (13,618 lines)
- **Monolithic components**: Large files with multiple responsibilities make testing and maintenance difficult
- **Inconsistent patterns**: Doesn't follow the auth-service's clean architecture patterns
- **Missing service layer**: Business logic is embedded directly in handlers and processors
- **Tight coupling**: Dependencies are mixed without clear abstraction layers

## Desired State

Transform image-service to match the auth-service's clean architecture:

```
image-service/
├── main.go                          # Clean setup with constructor injection
├── go.mod                           # Dependencies
├── go.sum                           # Dependency lock
├── Dockerfile                       # Multi-stage build
├── .dockerignore                    # Build optimization
├── internal/
│   ├── config/
│   │   └── config.go                # Configuration with validation
│   ├── models/
│   │   ├── instruction.go           # Instruction domain model
│   │   ├── image.go                 # Image domain model
│   │   └── processing_status.go     # Processing status enums
│   ├── handlers/
│   │   ├── instruction_handler.go   # Thin HTTP handlers
│   │   └── health_handler.go        # Health check handlers
│   ├── services/
│   │   ├── image_service.go         # Image processing business logic
│   │   ├── instruction_service.go   # Instruction management logic
│   │   └── processing_service.go    # Background processing logic
│   ├── repositories/
│   │   ├── instruction_repository.go    # Instruction data access
│   │   └── instruction_detail_repository.go # Detail data access
│   ├── validators/
│   │   ├── request_validator.go     # Input validation
│   │   └── image_validator.go       # Image-specific validation
│   ├── helpers/
│   │   ├── response_helper.go       # HTTP response utilities
│   │   ├── s3_helper.go             # S3 storage utilities
│   │   └── nats_helper.go           # Message bus utilities
│   ├── middleware/
│   │   └── auth_middleware.go       # Authentication middleware
│   └── errors.go                    # Error definitions
├── pkg/
│   └── utils/
│       ├── mime.go                  # MIME type utilities (existing)
│       └── slice.go                 # Slice utilities (existing)
└── static/
    └── swagger.json                 # API documentation
```

## Key Improvements

### 1. Feature-Based Organization
- **Clear separation**: Handlers, services, repositories, models, validators, and helpers in separate directories
- **Single responsibility**: Each directory and file has a focused purpose
- **Consistent with auth-service**: Follows established patterns across the microservice ecosystem

### 2. Service Layer Introduction
- **Extract business logic**: Move logic from `instruction_handler.go` (13,618 lines) into focused service classes
- **Image processing service**: Dedicated `image_service.go` for image processing operations
- **Instruction management**: `instruction_service.go` for instruction lifecycle management
- **Background processing**: `processing_service.go` for NATS-based background operations

### 3. Enhanced Configuration Management
- **Comprehensive config**: Expand beyond basic configuration to include all service settings
- **Environment validation**: Ensure all required configuration is present and valid
- **Default values**: Sensible defaults for development environments

### 4. Improved Testing Support
- **Clean interfaces**: Easy mocking for unit tests
- **Focused services**: Service layer enables isolated business logic testing
- **Thin handlers**: HTTP layer can be tested separately from business logic
- **Test utilities**: Helper functions for common test scenarios

### 5. Consistent Error Handling
- **Standardized errors**: Centralized error definitions and handling
- **Consistent responses**: Same response format as auth-service
- **Proper HTTP status codes**: Appropriate status codes for different error scenarios

### 6. Better Dependency Management
- **Constructor injection**: Clear dependency injection pattern
- **Interface-based design**: Easy to mock and test
- **Loose coupling**: Services depend on abstractions, not concrete implementations

## Scope

### In Scope
- **Directory restructuring**: Match auth-service's clean architecture
- **Service layer extraction**: Break down the 13,618-line handler into focused services
- **Configuration enhancement**: Comprehensive configuration with validation
- **Error handling standardization**: Consistent error patterns
- **Model separation**: Clear domain models with validation
- **Response standardization**: Consistent API response format
- **Testing support**: Clean interfaces for better testability
- **Documentation**: Updated API documentation and code comments

### Out of Scope
- **API contract changes**: Maintain existing API compatibility
- **Database schema changes**: No changes to MongoDB collections
- **External service integrations**: No changes to S3, NATS, or product client
- **Business logic changes**: Preserve existing functionality
- **Docker/container changes**: Keep existing deployment setup

## Benefits

1. **Consistency**: Aligns with auth-service and established microservice patterns
2. **Maintainability**: Smaller, focused files are easier to understand and modify
3. **Testability**: Clean separation enables comprehensive unit testing
4. **Readability**: Organized structure is easier to navigate and understand
5. **Scalability**: Clean architecture supports future enhancements
6. **Developer Experience**: Predictable structure makes onboarding easier
7. **Code Quality**: Better separation reduces technical debt

## Implementation Strategy

### Phase 1: Structure Setup (Day 1)
- Create new directory structure
- Move existing files to appropriate locations
- Set up basic configuration and models
- Ensure the service still builds and runs

### Phase 2: Service Layer Extraction (Day 2)
- Extract image processing logic into `image_service.go`
- Create `instruction_service.go` for instruction management
- Move background processing to `processing_service.go`
- Update handlers to use new services
- Ensure all existing functionality works

### Phase 3: Enhancement and Cleanup (Day 3)
- Add comprehensive validation
- Standardize error handling
- Improve configuration management
- Add proper response helpers
- Update documentation
- Add basic tests for new structure

## Success Criteria

- **Architecture alignment**: Image-service follows auth-service patterns exactly
- **Functionality preserved**: All existing features work unchanged
- **Code quality**: No file exceeds 1000 lines (except generated files)
- **Testability**: Services can be unit tested in isolation
- **Configuration**: All settings are validated with helpful error messages
- **Documentation**: API documentation is current and accurate
- **Build success**: Service builds and runs without errors

## Risk Mitigation

- **Incremental approach**: Changes are made in small, testable increments
- **Backward compatibility**: API contracts remain unchanged
- **Continuous testing**: Service is tested after each phase
- **Rollback ready**: Current code is backed up for quick rollback if needed

## Estimated Timeline

**Total Duration**: 3 days

- **Day 1**: Directory restructuring and basic setup
- **Day 2**: Service layer extraction and integration
- **Day 3**: Enhancement, testing, and documentation

This refactoring will significantly improve the image-service's maintainability, testability, and consistency with the broader microservice architecture while preserving all existing functionality.