## Context

The auth-service is a critical microservice handling authentication, PIN-based login, OAuth integration, and user management. Currently it has **zero test coverage**, which is unacceptable for a security-critical component. The service uses Fiber framework, MongoDB for persistence, JWT for tokens, and integrates with Google OAuth and email services.

The architecture follows a clean layered pattern with handlers → services → repositories → models, making it highly testable. All dependencies are injected through constructors, enabling easy mocking.

## Goals / Non-Goals

**Goals:**
- Achieve >90% test coverage across all components
- Test all authentication flows (PIN, OAuth, token management)
- Validate security scenarios and edge cases
- Create reusable test utilities and patterns
- Enable safe refactoring and prevent regressions
- Establish testing foundation for future development

**Non-Goals:**
- End-to-end UI testing (handled separately)
- Performance/load testing (different scope)
- Integration with other services (focused on unit tests)
- Production monitoring or observability setup

## Decisions

### Testing Framework Strategy
- **Decision**: Use Go's built-in `testing` package with `testify` for assertions and mocks
- **Rationale**: `testify` is already in dependencies, provides powerful assertions, test suites, and mocking capabilities
- **Alternatives considered**: Pure Go testing (more verbose), Ginkgo/Gomega (BDD style, adds complexity)

### Test Organization Pattern
- **Decision**: Co-locate test files with source files (`handler_test.go` next to `handler.go`)
- **Rationale**: Easy discovery, clear relationship between code and tests
- **Alternatives considered**: Separate test directories (harder to navigate)

### Mock Strategy
- **Decision**: Interface-based mocking using `testify/mock` for external dependencies
- **Rationale**: Clean interfaces already exist, enables isolated unit testing
- **Alternatives considered**: Test containers (heavy for unit tests), manual fakes (more maintenance)

### Database Testing
- **Decision**: Use MongoDB memory containers for repository integration tests
- **Rationale**: Realistic testing without external dependencies
- **Alternatives considered**: Pure mocking (less realistic), separate test database (complex setup)

## Risks / Trade-offs

- **Risk**: Test setup complexity due to MongoDB and external service dependencies
  - **Mitigation**: Create comprehensive test utilities and helper functions
- **Risk**: Slow test suite due to database operations
  - **Mitigation**: Use parallel tests where possible, optimize setup/teardown
- **Trade-off**: Test maintenance overhead vs. code safety
  - **Acceptance**: Essential for security-critical authentication service
- **Risk**: Test flakiness from timing issues in async operations
  - **Mitigation**: Use deterministic testing patterns, avoid real timing dependencies

## Migration Plan

1. **Phase 1**: Create test utilities and mock implementations
2. **Phase 2**: Test models and validators (isolated, easy wins)
3. **Phase 3**: Test repositories with in-memory MongoDB
4. **Phase 4**: Test services with mocked repositories
5. **Phase 5**: Test handlers with mocked services
6. **Phase 6**: Add integration tests for complete flows
7. **Phase 7**: Security-focused edge case testing

Each phase can be validated independently, enabling gradual progress and early feedback.

## Testing Architecture

### Layer-Specific Testing Patterns

**Handler Tests**:
- Use Fiber's testing utilities (`*fiber.App.Test()`)
- Mock service dependencies
- Test HTTP request/response cycles
- Validate status codes and response formats
- Cover authentication middleware integration

**Service Tests**:
- Mock repositories and external helpers
- Focus on business logic validation
- Test error handling and edge cases
- Validate transaction boundaries
- Cover authentication flows and security rules

**Repository Tests**:
- Use in-memory MongoDB instances
- Test data access operations
- Validate query logic and data transformations
- Test error scenarios and connection handling
- Cover all CRUD operations

**Model/Validator Tests**:
- Pure unit tests, no external dependencies
- Test validation rules and business logic
- Cover edge cases and boundary conditions
- Test method implementations and state transitions

### Security Testing Requirements

- **Authentication Bypass Attempts**: Invalid tokens, malformed credentials
- **Session Management**: Token expiry, refresh token rotation, multi-device scenarios
- **PIN Security**: Brute force protection, expiry handling, comparison security
- **OAuth Security**: Token validation, state parameter handling, callback verification
- **Input Validation**: SQL injection, XSS prevention, malformed request handling

## Test Data Management

**Factory Pattern**: Use factory functions for creating test users, sessions, tokens
**Data Cleanup**: Automatic cleanup between tests to ensure isolation
**Fixtures**: Predefined test scenarios for common authentication flows
**Randomization**: Use randomized test data where appropriate to avoid brittle tests

## Continuous Integration Integration

**Test Execution**: All tests run on every PR and merge
**Coverage Requirements**: Enforce minimum coverage thresholds
**Performance Gates**: Prevent test suite slowdown over time
**Security Scanning**: Include security-focused tests in CI pipeline

## Open Questions

- **Test Database Strategy**: Should we use Docker-in-Docker for MongoDB tests in CI?
- **Mock Granularity**: How much should we mock external HTTP clients vs. using real test servers?
- **Parallel Test Execution**: Which tests can safely run in parallel?
- **Test Data Seeding**: Should we use fixed test data or generate programmatically?
- **Timeout Values**: Appropriate timeouts for OAuth and external service calls in tests?