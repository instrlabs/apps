# Tasks: Image Service Architecture Refactoring

## Overview

This document outlines the specific tasks required to refactor the image-service to follow the clean architecture patterns established in the auth-service.

## Phase 1: Foundation Setup (Day 1)

### Task 1.1: Create Directory Structure
- [x] Create new directory structure matching auth-service pattern
- [x] Create `internal/config/` directory
- [x] Create `internal/models/` directory
- [x] Create `internal/handlers/` directory
- [x] Create `internal/services/` directory
- [x] Create `internal/repositories/` directory
- [x] Create `internal/validators/` directory
- [x] Create `internal/helpers/` directory
- [x] Create `internal/middleware/` directory

**Validation**: Verify all directories are created and basic files can be placed correctly. ✅

### Task 1.2: Move Existing Files to Appropriate Locations
- [ ] Move `config.go` → `internal/config/config.go`
- [ ] Move `instruction_repository.go` → `internal/repositories/instruction_repository.go`
- [ ] Move `instruction_detail_repository.go` → `internal/repositories/instruction_detail_repository.go`
- [ ] Move `product_client.go` → `internal/helpers/product_client.go`
- [ ] Keep `pkg/utils/` files in existing location

**Validation**: Ensure all files are moved and imports are updated correctly.

### Task 1.3: Setup Configuration Layer
- [ ] Enhance `internal/config/config.go` with comprehensive configuration
- [ ] Add configuration validation
- [ ] Add environment-specific settings
- [ ] Add default values for development
- [ ] Create `internal/config/validation.go` for config validation

**Validation**: Configuration loads successfully with proper validation and defaults.

### Task 1.4: Define Domain Models
- [ ] Create `internal/models/instruction.go` with instruction domain model
- [ ] Create `internal/models/image.go` with image domain model
- [ ] Create `internal/models/processing_status.go` with status enums
- [ ] Add validation methods to all models
- [ ] Add JSON marshaling/unmarshaling support

**Validation**: Models compile successfully and basic validation works.

### Task 1.5: Verify Basic Build
- [x] Update `main.go` imports for new structure
- [x] Ensure service builds without errors
- [x] Run basic smoke test to verify service starts

**Validation**: Service builds and starts successfully with new directory structure. ✅

## Phase 2: Service Layer Extraction (Day 2)

### Task 2.1: Create Image Service
- [x] Create `internal/services/image_service.go`
- [x] Extract image processing logic from `instruction_handler.go`
- [x] Create `ImageService` interface with core operations
- [x] Implement image validation methods
- [x] Add image format conversion operations
- [x] Integrate with S3 storage operations

**Validation**: Image service compiles and basic operations work correctly. ✅

### Task 2.2: Create Instruction Service
- [ ] Create `internal/services/instruction_service.go`
- [ ] Extract instruction management logic from handler
- [ ] Create `InstructionService` interface
- [ ] Implement CRUD operations for instructions
- [ ] Add status management methods
- [ ] Integrate with repositories

**Validation**: Instruction service handles all instruction operations correctly.

### Task 2.3: Create Processing Service
- [ ] Create `internal/services/processing_service.go`
- [ ] Extract background processing logic from `instruction_processor.go`
- [ ] Create `ProcessingService` interface
- [ ] Implement NATS message handling
- [ ] Add asynchronous processing operations
- [ ] Add cleanup operations

**Validation**: Processing service handles background operations correctly.

### Task 2.4: Refactor Instruction Handler
- [ ] Move `instruction_handler.go` → `internal/handlers/instruction_handler.go`
- [ ] Extract business logic to services (reduce from 13,618 lines)
- [ ] Implement dependency injection for services
- [ ] Add request validation using validators
- [ ] Standardize response format using helpers
- [ ] Add proper error handling

**Validation**: Handler is significantly smaller and delegates to services correctly.

### Task 2.5: Create Health Handler
- [ ] Create `internal/handlers/health_handler.go`
- [ ] Implement comprehensive health checks
- [ ] Add database connectivity check
- [ ] Add external service health checks
- [ ] Add service metrics endpoint

**Validation**: Health endpoints return proper service status information.

### Task 2.6: Update Main Application
- [ ] Refactor `main.go` to use new architecture
- [ ] Implement proper dependency injection
- [ ] Wire up all services and handlers
- [ ] Ensure graceful shutdown
- [ ] Maintain all existing functionality

**Validation**: Application starts and all endpoints work correctly.

## Phase 3: Enhancement and Cleanup (Day 3)

### Task 3.1: Implement Validation Layer
- [ ] Create `internal/validators/request_validator.go`
- [ ] Create `internal/validators/image_validator.go`
- [ ] Add comprehensive input validation
- [ ] Add image format validation
- [ ] Add business rule validation
- [ ] Integrate validation in handlers

**Validation**: All inputs are properly validated with helpful error messages.

### Task 3.2: Create Response Helpers
- [ ] Create `internal/helpers/response_helper.go`
- [ ] Standardize API response format
- [ ] Add success response helpers
- [ ] Add error response helpers
- [ ] Add pagination helpers
- [ ] Update all handlers to use helpers

**Validation**: All API responses follow consistent format.

### Task 3.3: Implement Error Handling
- [ ] Create `internal/errors.go` with service error definitions
- [ ] Add error categorization (validation, business, system)
- [ ] Implement proper HTTP status code mapping
- [ ] Add error logging and monitoring
- [ ] Update error handling in all layers

**Validation**: Errors are handled consistently with proper status codes.

### Task 3.4: Create Storage Helpers
- [ ] Create `internal/helpers/s3_helper.go`
- [ ] Abstract S3 operations from services
- [ ] Add error handling and retry logic
- [ ] Add upload/download utilities
- [ ] Update services to use helper

**Validation**: S3 operations work reliably with proper error handling.

### Task 3.5: Create Message Bus Helpers
- [ ] Create `internal/helpers/nats_helper.go`
- [ ] Abstract NATS operations
- [ ] Add connection management
- [ ] Add message publishing utilities
- [ ] Update processing service to use helper

**Validation**: Message bus operations work reliably with proper error handling.

### Task 3.6: Add Authentication Middleware
- [ ] Create `internal/middleware/auth_middleware.go`
- [ ] Implement JWT validation
- [ ] Add user context injection
- [ ] Add rate limiting
- [ ] Integrate with existing auth system

**Validation**: Authentication works correctly across all protected endpoints.

### Task 3.7: Update Documentation
- [ ] Update README.md with new architecture description
- [ ] Update API documentation in `static/swagger.json`
- [ ] Add code comments for all public methods
- [ ] Create development setup guide
- [ ] Add testing documentation

**Validation**: Documentation accurately reflects the new architecture.

### Task 3.8: Add Basic Tests
- [ ] Create test files for each service
- [ ] Add unit tests for image service
- [ ] Add unit tests for instruction service
- [ ] Add unit tests for validators
- [ ] Add integration tests for key workflows

**Validation**: Tests pass and provide good coverage of core functionality.

### Task 3.9: Performance Validation
- [ ] Run performance tests to ensure no regression
- [ ] Validate memory usage is acceptable
- [ ] Test concurrent request handling
- [ ] Validate background processing performance
- [ ] Monitor database query performance

**Validation**: Performance meets or exceeds current benchmarks.

### Task 3.10: Final Integration Testing
- [ ] Run full integration test suite
- [ ] Test all existing API endpoints
- [ ] Verify file upload/download functionality
- [ ] Test background processing workflows
- [ ] Validate error handling scenarios

**Validation**: All existing functionality works correctly with new architecture.

## Validation and Testing

### Per-Task Validation
Each task includes specific validation criteria to ensure successful completion.

### Phase Gates
- **Phase 1 → 2**: Service must build and start with new structure
- **Phase 2 → 3**: All business logic must be extracted to services
- **Phase 3 → Complete**: All functionality must work with new architecture

### Final Acceptance Criteria
- [ ] Service builds without errors
- [ ] All existing API endpoints work correctly
- [ ] File upload/download functionality preserved
- [ ] Background processing works as expected
- [ ] Configuration loads with proper validation
- [ ] Error handling is consistent and helpful
- [ ] Code follows established patterns from auth-service
- [ ] Documentation is accurate and complete
- [ ] Basic test coverage is in place
- [ ] Performance meets current standards

## Risk Mitigation

### Backup Strategy
- [ ] Create backup of current codebase before starting
- [ ] Use git branches for each phase
- [ ] Document rollback procedures

### Testing Strategy
- [ ] Test after each task completion
- [ ] Run integration tests after each phase
- [ ] Perform regression testing before completion

### Deployment Considerations
- [ ] Ensure environment variables are compatible
- [ ] Validate Docker build process
- [ ] Test deployment pipeline compatibility

This task breakdown provides a clear, step-by-step approach to refactoring the image-service while maintaining functionality and minimizing risk.