# Change: Simplify Auth Service Structure for Better Readability and Testability

## Why

The current auth-service implementation mixes business logic with HTTP handlers, making it hard to test and maintain. The massive `user_handler.go` (17,000+ lines) contains too much responsibility, and the directory structure doesn't clearly separate concerns. We need a simpler, more readable structure that makes testing easier while avoiding over-engineering.

## What Changes

- **Simple service layer** - Extract business logic from handlers into clean, focused services
- **Clear separation** - Handlers handle HTTP, Services handle business logic, Repositories handle data
- **Dependency injection** - Use simple constructor injection (no complex frameworks)
- **Organize by feature** - Group related code together instead of complex layered directories
- **Better testing** - Make each component independently testable with simple mocks
- **Cleaner handlers** - Thin handlers that delegate to services
- **Focused repositories** - Single responsibility for data access

## Impact

- Affected specs: `auth-service` (refactored structure)
- Affected code:
  - `auth-service/main.go` - Simple constructor-based initialization
  - `auth-service/internal/user_handler.go` - Split into thin handler + focused services
  - `auth-service/internal/user_repository.go` - Clean interface and implementation
  - `auth-service/internal/user.go` - Clean domain model
  - `auth-service/internal/config.go` - Add validation
  - `auth-service/internal/utils.go` - Split into focused helper files
- New structure:
  - `auth-service/internal/handlers/` - HTTP handlers (thin, focused)
  - `auth-service/internal/services/` - Business logic services
  - `auth-service/internal/repositories/` - Data access layer
  - `auth-service/internal/models/` - Domain models
  - `auth-service/internal/validators/` - Input validation
  - `auth-service/internal/helpers/` - Utility functions
- Testing improvements: Easy unit testing with simple constructor injection and minimal mocks
- Migration: Gradual refactoring with no breaking changes