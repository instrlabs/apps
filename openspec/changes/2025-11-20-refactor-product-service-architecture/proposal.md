# Proposal: Refactor Product Service Architecture

## Summary

Refactor the product-service to follow the same clean, feature-based architecture patterns established in the auth-service. The current product-service has a basic structure but lacks proper separation of concerns, testing support, and consistent patterns with the rest of the microservices.

**Proposed Change ID:** `refactor-product-service-architecture`

## Current State Analysis

The product-service currently has:

- **Simple but inconsistent structure**: Code is mixed in a single `internal/` directory without clear separation
- **Monolithic handlers**: `product_handler.go` contains business logic, validation, and HTTP concerns
- **Missing service layer**: Business logic is embedded directly in handlers
- **Limited testing support**: No clear separation makes testing difficult
- **Inconsistent patterns**: Doesn't follow the auth-service architecture patterns

## Desired State

Transform product-service to match the auth-service's clean architecture:

```
product-service/
├── main.go                    # Clean setup with constructor injection
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration with validation
│   ├── models/
│   │   └── product.go        # Product domain model with validation
│   ├── handlers/
│   │   └── product_handler.go # Thin HTTP handlers
│   ├── services/
│   │   └── product_service.go # Business logic services
│   ├── repositories/
│   │   └── product_repository.go # Data access layer
│   ├── validators/
│   │   └── request_validator.go # Input validation
│   └── helpers/
│       └── response_helper.go # HTTP response utilities
```

## Key Improvements

### 1. Feature-Based Organization
- Separate handlers, services, repositories, models, validators, and helpers
- Clear single-responsibility directories
- Consistent with auth-service patterns

### 2. Service Layer Introduction
- Extract business logic from handlers into dedicated service classes
- Enable better testing and reusability
- Clear separation between HTTP and business concerns

### 3. Enhanced Configuration
- Comprehensive configuration with validation
- Environment-specific settings
- Clear error messages for missing configurations

### 4. Improved Testing Support
- Clean interfaces for easy mocking
- Focused service layer for business logic testing
- Separated HTTP handlers for request/response testing

### 5. Consistent Patterns
- Follow auth-service architecture patterns
- Same dependency injection approach
- Consistent error handling and response formats

## Scope

### In Scope
- Restructure directories to match auth-service patterns
- Extract business logic into service layer
- Add proper input validation layer
- Implement consistent configuration management
- Add comprehensive error handling
- Maintain full API backward compatibility

### Out of Scope
- Changing API contracts or response formats
- Adding new business features
- Database schema changes
- External service integrations

## Benefits

1. **Consistency**: Aligns with established microservice patterns
2. **Maintainability**: Clear separation makes code easier to modify
3. **Testability**: Service layer enables focused unit testing
4. **Readability**: Organized structure is easier to navigate
5. **Scalability**: Clean architecture supports future enhancements

## Implementation Timeline

**Estimated Duration**: 2-3 days

- **Day 1**: Directory restructuring and basic setup
- **Day 2**: Service layer extraction and testing
- **Day 3**: Handler refactoring and validation

## Success Criteria

- Product-service follows auth-service architecture patterns
- All existing functionality remains unchanged
- Code is properly separated into logical layers
- Service can be easily tested with minimal setup
- Configuration is comprehensive and validated