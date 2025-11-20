# Design: Simplified Auth Service Architecture

## Context

The auth-service currently has a massive `user_handler.go` file (17,000+ lines) that mixes HTTP handling with business logic. This makes it hard to test, maintain, and understand. The proposed complex layered architecture is over-engineered for the current needs and would add unnecessary complexity.

**Current Issues:**
- Single massive handler file with too many responsibilities
- Business logic mixed with HTTP concerns
- Hard to test individual components in isolation
- File organization doesn't reflect logical groupings
- No clear separation between data access and business logic

**Goals:**
- Make code readable and easy to understand
- Enable easy testing with simple setup
- Avoid over-engineering and complex abstractions
- Keep the learning curve low for new developers
- Maintain existing functionality without breaking changes

## Decisions

### Decision 1: Feature-Based Organization (Not Layered)

**What:** Organize code by feature rather than complex architectural layers:

```
internal/
├── handlers/          # HTTP handlers (thin, focused)
├── services/          # Business logic services
├── repositories/      # Data access layer
├── models/            # Domain models
├── validators/        # Input validation
├── helpers/           # Utility functions
└── config/            # Configuration
```

**Why:**
- Easier to find related code
- Less cognitive overhead than complex layered architecture
- Groups code by what it does, not by technical concerns
- Each directory has a clear, single purpose

### Decision 2: Simple Constructor Injection (No Frameworks)

**What:** Use plain constructor injection for dependencies:

```go
type AuthService struct {
    userRepo repositories.UserRepositoryInterface
    emailService helpers.EmailSender
}

func NewAuthService(userRepo repositories.UserRepositoryInterface, emailService helpers.EmailSender) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        emailService: emailService,
    }
}
```

**Why:**
- No dependency on complex DI frameworks
- Easy to understand and debug
- Explicit dependency relationships
- Simple to test with manual mocks

### Decision 3: Focused Service Classes

**What:** Create small, focused service classes instead of large, complex ones:

- `AuthService` - Login, logout, token management
- `PinService` - PIN generation and validation
- `OAuthService` - Google OAuth flows
- `UserService` - User management and profile

**Why:**
- Each service has single responsibility
- Easier to test in isolation
- Clear boundaries between different concerns
- Smaller files are easier to navigate

### Decision 4: Thin HTTP Handlers

**What:** Handlers should only handle HTTP-specific concerns:

```go
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // 1. Parse and validate request
    req, err := h.validator.ParseLoginRequest(c)
    if err != nil {
        return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
    }

    // 2. Call service
    result, err := h.authService.Login(req.Email, req.Pin)
    if err != nil {
        return h.sendErrorResponse(c, fiber.StatusBadRequest, err.Error())
    }

    // 3. Send response
    return h.sendSuccessResponse(c, result)
}
```

**Why:**
- Clear separation of HTTP and business concerns
- Easy to test by mocking services
- Handlers focus on request/response handling only
- Business logic is reusable and testable

### Decision 5: Clean Repository Interface

**What:** Simple, clear repository interfaces:

```go
type UserRepositoryInterface interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    FindByID(id string) (*models.User, error)
    Update(user *models.User) error
    AddRefreshToken(userID, token string) error
    RemoveRefreshToken(userID, token string) error
    ClearAllRefreshTokens(userID string) error
}
```

**Why:**
- Clear contract for data operations
- Easy to mock for testing services
- Single responsibility for data access
- Can swap implementations easily

## Target Directory Structure

```
auth-service/
├── main.go                    # Simple setup with constructor injection
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration with validation
│   ├── models/
│   │   └── user.go           # User domain model
│   ├── handlers/
│   │   ├── auth_handler.go   # Auth endpoints (login, logout, refresh)
│   │   ├── user_handler.go   # User endpoints (profile)
│   │   ├── oauth_handler.go  # OAuth endpoints (Google)
│   │   └── pin_handler.go    # PIN endpoints (send-pin)
│   ├── services/
│   │   ├── auth_service.go   # Login, logout, tokens
│   │   ├── pin_service.go    # PIN generation & validation
│   │   ├── oauth_service.go  # Google OAuth flow
│   │   └── user_service.go   # User management
│   ├── repositories/
│   │   └── user_repository.go # Data access layer
│   ├── validators/
│   │   └── request_validator.go # Input validation
│   └── helpers/
│       ├── email_helper.go   # Email sending
│       ├── jwt_helper.go     # JWT token generation
│       └── utils_helper.go   # Utility functions
└── tests/                    # Integration tests
```

## Migration Strategy

### Phase 1: Setup Structure (Day 1)
1. Create new directory structure
2. Move and organize existing files
3. Update imports

### Phase 2: Extract Services (Day 2-3)
1. Create service interfaces and implementations
2. Extract business logic from handlers
3. Test services independently

### Phase 3: Simplify Handlers (Day 3-4)
1. Rewrite handlers to use services
2. Add input validation layer
3. Update HTTP response handling

### Phase 4: Testing & Cleanup (Day 4-5)
1. Add comprehensive tests
2. Update documentation
3. Remove old code

## Benefits Over Complex Architecture

**Simpler Learning Curve:**
- New developers can understand structure quickly
- No complex architectural patterns to learn
- Clear organization by feature

**Better Testability:**
- Services can be tested in isolation
- Simple constructor injection
- Minimal mocking requirements

**Easier Maintenance:**
- Small, focused files
- Clear responsibilities
- Less boilerplate code

**Faster Development:**
- No complex DI framework setup
- Simple dependency management
- Direct code organization

## When to Evolve

This simple approach is suitable for the current auth-service size and complexity. Consider evolving to more complex architecture only when:

- Service grows significantly (>50 endpoints)
- Multiple databases or external systems are added
- Complex business rules emerge
- Team size grows large
- Performance requires sophisticated caching

For now, keep it simple and focus on code quality and readability.