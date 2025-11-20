## MODIFIED Requirements

### Requirement: Auth Service API Documentation
The auth-service SHALL provide comprehensive Swagger/OpenAPI 3.0.3 documentation that accurately reflects all current API endpoints, request/response schemas, and authentication requirements organized according to the refactored handler structure.

#### Scenario: Complete API Documentation Coverage
- **WHEN** developers access the Swagger documentation
- **THEN** all endpoints (/login, /logout, /refresh, /send-pin, /profile, /google, /google/callback, /health) SHALL be documented with correct HTTP methods, parameters, and response schemas

#### Scenario: Handler-aligned Organization
- **WHEN** viewing the API documentation
- **THEN** endpoints SHALL be organized by handler groups (AuthHandler, PinHandler, OAuthHandler, UserHandler) with appropriate tags

#### Scenario: Accurate Request/Response Models
- **WHEN** integrating with auth-service endpoints
- **THEN** request body schemas SHALL match actual validator implementations and response schemas SHALL match actual service responses

#### Scenario: Complete Error Documentation
- **WHEN** API calls result in errors
- **THEN** all possible error responses SHALL be documented with correct status codes, error messages, and response structures

#### Scenario: Authentication Documentation
- **WHEN** accessing protected endpoints
- **THEN** authentication requirements SHALL be clearly documented with bearer token specifications

#### Scenario: Valid Examples
- **WHEN** developers review endpoint examples
- **THEN** all examples SHALL be syntactically correct and reflect actual current API behavior

#### Scenario: OpenAPI Specification Compliance
- **WHEN** validating the Swagger documentation
- **THEN** the documentation SHALL comply with OpenAPI 3.0.3 specification and be accessible via Swagger UI