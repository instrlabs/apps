# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based application with a Next.js frontend and four Go backend services communicating via NATS message broker. All services use the Fiber web framework and share common infrastructure through `github.com/instrlabs/shared`.

### Services

**Gateway Service** (`gateway-service/`) - Port 3000
- API gateway that proxies requests to downstream services
- Handles CORS, rate limiting, CSRF protection, and JWT authentication
- Extracts JWT tokens from cookies and injects user metadata headers (`x-authenticated`, `x-user-id`, `x-user-roles`)
- Routes are dynamically configured via `Services` array in `internal/config.go`
- Health check endpoint at `/health` for service availability monitoring

**Auth Service** (`auth-service/`) - Port 3001
- User authentication with email/PIN and Google OAuth
- JWT token generation and refresh token management
- MongoDB for user storage with bcrypt password hashing
- Unauthenticated endpoints: `/login`, `/refresh`, `/send-pin`, `/google`, `/google/callback`
- User profile retrieval requires authentication

**Image Service** (`image-service/`) - Port 3002
- Image processing instruction management with full CRUD operations
- S3/MinIO storage integration for input/output files
- NATS subscriber for async image processing (`image.requests` subject)
- Background cleanup job that runs every 30 minutes to remove processed files
- MongoDB collections: `products`, `instructions`, `files`
- File repository pattern for managing S3 metadata

**Notification Service** (`notification-service/`) - Port 3003
- Server-Sent Events (SSE) for real-time notifications
- NATS subscriber for notification messages (`notifications.sse` subject)
- All endpoints require authentication
- Streams notifications to connected clients

**Web Application** (`web/`) - Port 8000
- Next.js 15.4.7 with App Router and React 19
- Route groups: `(non-auth)` for public pages, `(site)` for authenticated, `debug/` for dev testing, `api/` for routes
- Server-side fetch utilities for API integration
- Tailwind CSS v4 with Geist font and strict TypeScript
- Comprehensive overlay system and SSE integration for real-time notifications

## Development Commands

### Backend Services

**Run a service locally (with hot reload):**
```bash
cd <service-name>
cp .env.example .env  # Create config first (edit with your values)
go run main.go
```

**Build a service binary:**
```bash
cd <service-name>
go build -o app .
```

**Run tests:**
```bash
cd <service-name>
go test ./...              # Run all tests in service
go test ./internal/...     # Run only internal package tests
go test -v ./...           # Verbose output with test names
go test -run TestName ...  # Run specific test by name
```

**Check code quality:**
```bash
cd <service-name>
go fmt ./...     # Format code
go vet ./...     # Run linter
go mod tidy      # Clean up dependencies
```

**Full Docker Compose stack (includes all services + dependencies):**
```bash
# Build and start all services with hot reload
# All services read from root .env file
docker-compose up --build

# Available services:
# - Gateway: http://localhost:3000
# - Auth: http://localhost:3001
# - Image: http://localhost:3002
# - Notification: http://localhost:3003
# - Web: http://localhost:8000

# Stop all services
docker-compose down

# View logs for specific service
docker-compose logs -f <service-name>

# Environment configuration:
# All services use single root .env file
# No service-level .env files (centralized secrets)
```

### Frontend

**Install dependencies:**
```bash
cd web
npm install
```

**Development server (with hot reload):**
```bash
cd web
npm run dev  # Starts on port 8000 with Turbopack
```

**Production build and start:**
```bash
cd web
npm run build  # Create optimized production bundle
npm start      # Run production server on port 8000
```

**Code quality:**
```bash
cd web
npm run lint           # Run ESLint
npm run format         # Format with Prettier
npm run format:check   # Check formatting without writing
```

## Environment Configuration

**Single source of truth:** All configuration is centralized in the root `.env` file. All services (backend and frontend) read from this single file.

### `.env` - Root Configuration (NEVER Commit)
**Purpose:** Central configuration file for entire application - all services read from here.

**Location:** `/.env` (root of project)

**Characteristics:**
- Listed in `.gitignore` (never committed to git)
- Contains ALL secrets and configuration for all services
- Single file that all services read (Docker Compose, local development, CI/CD)
- Each environment (dev/staging/prod) has its own `.env` file

**No service-level .env files** - Configuration is NOT duplicated across services

### Configuration Structure

All variables are defined in a single `.env` file and used by all services:

**Service URLs (Frontend Integration):**
```
WEB_URL=http://localhost:8000
GATEWAY_URL=http://localhost:3000
NOTIFICATION_URL=http://instrlabs-notification-service:3000
AUTH_SERVICE_URL=http://instrlabs-auth-service:3000
IMAGE_SERVICE_URL=http://instrlabs-image-service:3000
```

**Shared Secrets (Must Match Across Services):**
```
JWT_SECRET=<secret-key>
```

**Database & Message Bus:**
```
MONGO_URI=mongodb://user:pass@host:27017
MONGO_DB=instrlabs
NATS_URI=nats://host:4222
NATS_SUBJECT_IMAGE_REQUESTS=image.requests
NATS_SUBJECT_NOTIFICATIONS_SSE=notifications.sse
```

**Email Configuration (Auth Service):**
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=465
SMTP_USERNAME=email@example.com
SMTP_PASSWORD=app-password
EMAIL_FROM=noreply@example.com
```

**S3/Object Storage (Image Service):**
```
S3_ENDPOINT=storage.example.com
S3_REGION=us-east-1
S3_ACCESS_KEY=key
S3_SECRET_KEY=secret
S3_BUCKET=instrlabs
S3_USE_SSL=true
```

**Google OAuth (Auth Service):**
```
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=secret
```

**Feature Flags:**
```
PIN_FLAG=false
CSRF_ENABLED=true
```

### How Services Access Configuration

**Docker Compose:**
- All services load from `.env`: `env_file: - ./.env`
- All services share same secrets automatically

**Local Development (Go services):**
```bash
cd auth-service
go run main.go  # Reads from root .env via config loader
```

**Local Development (Next.js):**
```bash
cd web
npm run dev  # Next.js reads from root .env
```

**Configuration Loading Priority:**
1. **Environment variables** (highest priority - CI/CD, system env)
2. **`.env` file** (root level - your local configuration)
3. **Hardcoded defaults** (lowest priority - fallbacks in config.go)

### Benefits of Centralized Configuration

✅ **No duplication** - Single source of truth
✅ **Consistency** - All services use same values
✅ **Secrets management** - All secrets in one place
✅ **Easy deployment** - Single `.env` per environment
✅ **Docker friendly** - All containers read same file
✅ **Clear structure** - All variables visible in one place

### Critical Shared Variables (Must Match Across Services)

- **`JWT_SECRET`** - Secret key for JWT signing/validation (gateway and auth services use this)
- **`MONGO_URI`** - MongoDB connection URI
- **`MONGO_DB`** - Database name
- **`NATS_URI`** - NATS message broker connection

## Shared Infrastructure Module

All services depend on `github.com/instrlabs/shared` (v0.0.10+) which provides initialization helpers and common patterns.

**Initialization helpers** (`initx` package):
- `NewMongo(config)` - MongoDB client setup with connection pooling
- `NewS3(config)` - MinIO/S3 client setup (image-service only)
- `NewNats(uri)` - NATS message broker connection
- `SetupPrometheus(app)` - Prometheus metrics endpoint at `/metrics`
- `SetupServiceHealth(app)` - Health check at `/health` with dependency status
- `SetupLogger(app)` - Request logging middleware with structured output
- `SetupServiceSwagger(app)` - Swagger UI at `/swagger/*`
- `SetupAuthenticated(app, exclusions)` - Authentication middleware that checks `x-authenticated` header

**Standard service initialization pattern:**
```go
// Load config from environment
cfg := internal.LoadConfig()

// Initialize dependencies
mongo := initx.NewMongo(&initx.MongoConfig{
    MongoURI: cfg.MongoURI,
    MongoDB: cfg.MongoDB,
})
defer mongo.Close()

nats := initx.NewNats(cfg.NatsURI)
defer nats.Close()

// Create Fiber app and setup middleware
app := fiber.New(fiber.Config{})
initx.SetupPrometheus(app)
initx.SetupLogger(app)
initx.SetupServiceSwagger(app)
initx.SetupServiceHealth(app)
initx.SetupAuthenticated(app, []string{"/public-endpoint"})

// Register routes and start server
log.Fatal(app.Listen(cfg.Port))
```

**Updating the shared module:**
```bash
# In any service directory:
go get github.com/instrlabs/shared@v0.0.X  # Update to specific version
go mod tidy                                  # Clean up dependencies
```

## Authentication Flow

The system uses JWT tokens stored in HTTP-only cookies with header-based metadata forwarding through the gateway.

**Login flow:**
1. Client POSTs credentials to `/auth/login`
2. Auth service validates and returns `access_token` (short-lived) + `refresh_token` (long-lived) in Set-Cookie headers
3. Frontend stores tokens in HTTP-only cookies automatically via fetch utilities

**Request flow:**
1. Client sends request with cookies to gateway
2. Gateway middleware (`gateway-service/internal/middleware.go`):
   - Extracts and validates JWT from `access_token` cookie
   - Decodes token using `JWT_SECRET` (must match auth service)
   - Strips original `cookie` header from downstream request
   - Injects derived headers into request:
     - `x-authenticated: true/false` - Whether token was valid
     - `x-user-id: <id>` - User MongoDB ObjectID
     - `x-user-roles: <comma-separated>` - User roles
     - `x-gateway: true` - Indicates request came through gateway
   - Returns 401 if token exists but is invalid
3. Service middleware (`SetupAuthenticated`):
   - Checks `x-authenticated` header
   - Returns 401 if false and endpoint is not in public exclusion list
   - No 401 for missing tokens on public endpoints
4. Handler accesses user info from headers via `c.Get()` or middleware context

**Token refresh flow:**
- If `access_token` cookie missing but `refresh_token` exists, frontend middleware (`web/middleware.ts`) automatically calls `/auth/refresh`
- Gateway forwards `refresh_token` as `x-user-refresh` header on this endpoint only
- Auth service validates refresh token and returns new `access_token` + `refresh_token` pair
- Frontend stores new tokens in cookies for next request

**Important security notes:**
- `JWT_SECRET` must be identical across gateway and auth services (used for token validation)
- Tokens stored in HTTP-only cookies are not accessible to JavaScript
- Gateway strips cookies before forwarding to prevent token leakage to internal services
- Public endpoints must be explicitly listed in `SetupAuthenticated()` exclusion list

## Message Bus (NATS)

NATS is used for asynchronous inter-service communication and event broadcasting.

**Subject definitions and flow:**

**`image.requests`** - Image processing task queue
- **Publisher**: Image service handlers (`POST /instructions/:id/details`)
- **Subscriber**: Image service worker (main.go goroutine)
- **Payload**: ObjectID of instruction input (hex-encoded in message data)
- **Flow**: Handler publishes input ID → Worker receives message → Calls `RunInstructionMessage()` → Processes image
- **Pattern**: Fire-and-forget; multiple workers can subscribe for scalability

**`notifications.sse`** - Real-time notification events
- **Publisher**: Image service (`POST /instructions/:id/details` completion, status updates)
- **Subscriber**: Notification service for SSE broadcasting
- **Payload**: JSON-encoded notification message with user/event data
- **Flow**: Image service publishes event → Notification service receives → Broadcasts to connected SSE clients
- **Pattern**: One-to-many broadcast

**Standard subscription pattern:**
```go
// In main.go during service initialization
nats := initx.NewNats(cfg.NatsURI)
defer nats.Close()

// Subscribe to message subject
nats.Conn.Subscribe(cfg.NatsSubjectImageRequests, func(m *natsgo.Msg) {
    // Handle message asynchronously
    handler.ProcessMessage(m.Data)
})

// Start service (subscribers remain active)
log.Fatal(app.Listen(cfg.Port))
```

**Publishing messages:**
```go
// Publish message to subject
payload := []byte("message-content")
if err := h.nats.Conn.Publish(cfg.NatsSubjectNotificationsSSE, payload); err != nil {
    log.Errorf("Failed to publish: %v", err)
}
```

**Configuration:**
- `NATS_URI` - Connection string (e.g., `nats://localhost:4222`)
- `NATS_SUBJECT_IMAGE_REQUESTS` - Subject name for image processing tasks
- `NATS_SUBJECT_NOTIFICATIONS_SSE` - Subject name for notification events

## MongoDB Collections

Services use MongoDB for persistence following the repository pattern with standard CRUD methods.

**Auth Service:**
- **`users`** - User accounts
  - Fields: `_id` (ObjectID), `email`, `password_hash` (bcrypt), `roles` (array), `created_at`, `updated_at`
  - Unique index on email
  - Password hashing handled by bcrypt in user handler

**Image Service:**
- **`products`** - Available image processing products
  - Fields: `_id`, `name`, `description`, `version`
  - Reference data populated during setup

- **`instructions`** - Image processing jobs
  - Fields: `_id`, `user_id` (foreign key to users), `product_id`, `status`, `created_at`, `updated_at`
  - Status values: `pending`, `processing`, `completed`, `failed`
  - Indexed on `user_id` and `status` for efficient queries

- **`files`** - Input/output file metadata
  - Fields: `_id`, `instruction_id` (foreign key), `type` (`input`/`output`), `file_path` (S3 key), `size`, `mime_type`, `status`, `created_at`
  - Actual file content stored in S3 bucket
  - S3 key format: `instructions/{instruction_id}/{type}/{filename}`

**Repository pattern:**
```go
// Each collection has a repository with standard methods:
func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Model, error)
func (r *Repository) Create(ctx context.Context, model *Model) error
func (r *Repository) Update(ctx context.Context, model *Model) error
func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error
func (r *Repository) List(ctx context.Context, filter bson.M) ([]*Model, error)
```

**Usage in handlers:**
```go
// Access repository from handler
user, err := h.userRepo.FindByID(context.Background(), userID)
if err != nil {
    return c.Status(400).JSON(fiber.Map{"errors": err.Error()})
}
```

## API Response Format

All API responses follow a standard envelope structure for consistent error handling and type safety.

**Successful response (2xx):**
```json
{
  "success": true,
  "message": "Resource created successfully",
  "data": { "id": "...", "name": "..." },
  "errors": null
}
```

**Error response (4xx/5xx):**
```json
{
  "success": false,
  "message": "Invalid input",
  "data": null,
  "errors": {
    "email": "Email already exists",
    "password": "Password too short"
  }
}
```

**Handler implementation pattern:**
```go
func (h *Handler) CreateResource(c *fiber.Ctx) error {
    var req CreateRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Invalid request",
            "errors": map[string]string{"body": err.Error()},
            "data": nil,
        })
    }

    // Business logic...
    resource, err := h.repo.Create(context.Background(), &model)

    return c.Status(201).JSON(fiber.Map{
        "success": true,
        "message": "Resource created successfully",
        "data": resource,
        "errors": nil,
    })
}
```

**Frontend integration:**
- All fetch utilities (`web/utils/fetch.ts`) automatically parse and handle this response format
- `fetchPOST()`, `fetchGET()`, etc. validate the `success` field
- Error responses are automatically typed with form field error handling
- Type-safe TypeScript interfaces provided for each endpoint

## Docker Build Pattern

All Go services use multi-stage Docker builds for optimized container images.

**Build stages:**
1. **Builder stage** - `golang:1.24` image
   - Sets `CGO_ENABLED=0` for static binary compilation (no system library dependencies)
   - Compiles Go code with full optimizations
   - Results in statically-linked binary that runs anywhere

2. **Release stage** - `alpine:latest` image
   - Minimal base image (~5MB) for small container size
   - Includes `dumb-init` for proper signal handling and zombie process prevention
   - Copies binary from builder stage

**Port configuration:**
- All services expose port `3000` internally within the container
- Docker Compose maps to different host ports:
  - Gateway: 3000 → 3000
  - Auth: 3000 → 3001
  - Image: 3000 → 3002
  - Notification: 3000 → 3003

**Build locally:**
```bash
cd <service-name>
docker build -t instrlabs-<service>:latest .
docker run -p 3001:3000 instrlabs-<service>:latest
```

## CI/CD

**GitHub Actions workflow** (`.github/workflows/image-builder.yml`):
- Manual workflow dispatch (not automatic on push)
- Selectable service and environment
- Builds Docker images for `linux/amd64` platform
- Pushes to Docker Hub with tags: `<dockerhub-user>/instrlabs-<service>:<git-sha>-<environment>`
- Supports incremental build caching

**To trigger a build:**
1. Go to repository Actions tab
2. Select "image-builder" workflow
3. Click "Run workflow"
4. Select service and environment from dropdowns
5. Confirm build

## Code Patterns

### Service Project Structure
```
<service>/
├── Dockerfile                # Multi-stage build configuration
├── go.mod                    # Module definition and dependencies
├── go.sum                    # Dependency checksums
├── main.go                   # Entry point: initialization and route setup
├── internal/
│   ├── config.go            # Environment variable parsing with validation
│   ├── *_handler.go         # HTTP request handlers (Fiber framework)
│   ├── *_repository.go      # MongoDB/database access layer
│   ├── *.go                 # Domain models and business logic
│   └── model.go             # Data structures
└── static/                  # Static assets (if any)
```

### Handler Method Pattern
```go
// All handlers follow this pattern with Fiber context
func (h *Handler) CreateResource(c *fiber.Ctx) error {
    // 1. Parse and validate request
    var req CreateRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid request body",
            "errors": map[string]string{"body": err.Error()},
            "data": nil,
        })
    }

    // 2. Validate input
    if err := req.Validate(); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Validation failed",
            "errors": err.ValidationErrors(),
            "data": nil,
        })
    }

    // 3. Extract user from headers (set by gateway middleware)
    userID := c.Get("x-user-id")

    // 4. Business logic using repositories
    result, err := h.repo.Create(context.Background(), &req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to create resource",
            "errors": nil,
            "data": nil,
        })
    }

    // 5. Return response
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Resource created successfully",
        "data": result,
        "errors": nil,
    })
}
```

### Adding New Routes

**Step 1: Create handler method**
```go
// internal/resource_handler.go
func (h *ResourceHandler) GetResource(c *fiber.Ctx) error {
    // Implementation
}
```

**Step 2: Register route in main.go**
```go
// For auth/image/notification services:
app.Get("/resources/:id", resourceHandler.GetResource)

// For gateway, routes are auto-registered from config
```

**Step 3: Add to public exclusion list if needed**
```go
// If endpoint should be public, add to SetupAuthenticated exclusion list
initx.SetupAuthenticated(app, []string{
    "/login",
    "/refresh",
    "/public-resource",  // New public endpoint
})
```

### Gateway Route Configuration

The gateway proxies requests to downstream services based on configured prefixes:

```go
// gateway-service/internal/router.go - Dynamic route registration
for _, srv := range config.Services {
    app.All(srv.Prefix+"/*", proxyHandler)
}
```

**Service prefix configuration** (in gateway `.env`):
- `/auth/*` → Auth service at `http://auth-service:3000`
- `/images/*` → Image service at `http://image-service:3000`
- `/notifications/*` → Notification service at `http://notification-service:3000`

**Request flow through gateway:**
1. Client → `GET /auth/profile` to gateway
2. Gateway middleware processes request, injects headers
3. Gateway strips cookies, forwards to `/profile` on auth-service
4. Service processes with injected headers (x-authenticated, x-user-id, etc.)
5. Response returned to client with Set-Cookie headers

## Web Application Architecture

### Frontend Tech Stack
- **Framework:** Next.js 15.4.7 with App Router and React 19
- **Styling:** Tailwind CSS v4 with Geist font family
- **Type Safety:** TypeScript with strict mode enabled
- **Forms:** react-hook-form for form state management
- **Utilities:** date-fns, lodash.debounce

### Route Structure and Layout Organization

The web application uses Next.js route groups to organize pages with different layouts:

**`(non-auth)` group** - Public pages without authentication
- `/login`, `/register`, `/forgot-password`, `/reset-password`
- Unauthenticated layout without navigation or context providers
- Middleware redirects authenticated users away from these routes

**`(site)` group** - Authenticated pages with full layout
- All protected routes after login
- Nested provider hierarchy:
  ```
  ProfileProvider → ProductProvider → SSEProvider → NotificationProvider →
  ModalProvider → OverlayProvider → [OverlayTop + OverlayContent + NotificationWidget]
  ```
- Layout file: `app/(site)/layout.tsx`

**`debug/` group** - Development and testing pages
- Internal testing endpoints
- Not included in production builds

**`api/` routes** - Backend API integration
- `/api/sse` - Server-Sent Events proxy that forwards to private notification service
- Forwards cookies and headers to upstream services
- Handles connection errors and client disconnects

### Frontend Authentication

**Middleware-based flow** (`middleware.ts`):
1. All routes except whitelist (`/login`, `/register`, `/forgot-password`, `/reset-password`, `/`) require authentication
2. Checks for `access_token` cookie; if missing but `refresh_token` exists, calls `/auth/refresh`
3. Forwards client metadata headers to gateway:
   - `x-user-ip` - Client IP address
   - `x-user-agent` - Browser user agent
   - `x-user-host` - Request host
   - `x-user-origin` - Request origin
4. Redirects to `/login` on authentication failure

**Token management:**
- `access_token` and `refresh_token` stored in HTTP-only cookies (not accessible to JavaScript)
- Tokens set automatically by `fetchPOST()` on successful login/refresh calls
- Short-lived access tokens trigger automatic refresh flow

### API Integration Pattern

**Server-side fetch utilities** (`utils/fetch.ts` - all are server actions with `"use server"`):
- `fetchGET(url, options?)` - GET requests, returns parsed JSON
- `fetchPOST(url, body, options?)` - POST requests with JSON body
- `fetchPUT(url, body, options?)` - PUT requests
- `fetchPATCH(url, body, options?)` - PATCH requests
- `fetchGETBytes(url, options?)` - Binary data retrieval (e.g., images, files)
- `fetchPOSTFormData(url, formData, options?)` - Multipart form uploads

**All fetch utilities:**
- Route through `GATEWAY_URL` environment variable
- Automatically forward client metadata headers from middleware
- Automatically handle cookie-based authentication
- Automatically set response tokens on successful login/refresh
- Parse response and validate `success` field
- Type-safe TypeScript response handling

**Response structure** (standard across all endpoints):
```typescript
{
  success: boolean;
  message: string;
  data: T | null;
  errors: {
    [field: string]: string | string[];
  } | null;
}
```

### Overlay System

**Centralized sliding overlay management** (`hooks/useOverlay.tsx`):
- Overlays slide in from left or right edges of screen
- Registration-based system with string keys (e.g., `"left:navigation"`, `"right:profile"`)
- Context-based state management shared across all components

**Available actions:**
```typescript
const { openLeft, closeLeft, openRight, closeRight } = useOverlay();
openLeft("left:navigation");    // Open navigation overlay
closeLeft();                      // Close left overlay
openRight("right:profile");      // Open profile overlay
closeRight();                     // Close right overlay
```

**Registering new overlays:**
- Define overlay component in `components/overlays/<name>.tsx`
- Register in `registerOverlays()` function with unique key
- Reference key in `useOverlay()` calls

### Component Organization

Components organized by function in subdirectories under `components/`:
- **`actions/`** - Interactive buttons, clickable elements, CTAs
- **`cards/`** - Card layouts, containers, list items
- **`feedback/`** - Loading states, spinners, alerts, toasts, skeletons
- **`icons/`** - SVG icon components (use `currentColor` for styling)
- **`inputs/`** - Form inputs, text fields, select menus, checkboxes
- **`layouts/`** - Page structure (OverlayTop, OverlayContent, main content area)
- **`navigation/`** - Navigation menus, breadcrumbs, tabs
- **`overlays/`** - Overlay content panels (navigation, notifications, profile menus)

### Custom Hooks

- **`useOverlay`** - Control left/right sliding overlay state and actions
- **`useModal`** - Modal dialog open/close and state management
- **`useNotification`** - Toast/notification display system
- **`useSSE`** - Server-Sent Events connection management and message handling
- **`useProduct`** - Access product list from ProductProvider context
- **`useProfile`** - Access authenticated user profile from ProfileProvider context
- **`useMediaQuery`** - Responsive breakpoint detection for mobile/tablet/desktop

### Styling System

**Core principle:** Use Tailwind's built-in utility classes exclusively. No CSS custom properties, no `@utility` definitions, no custom CSS classes.

**Color system with opacity:**
- `bg-white`, `bg-black` - Full opacity backgrounds
- `bg-white/90` - 90% opacity (for elevated elements)
- `bg-white/8` - 8% opacity (for subtle backgrounds)
- `text-white`, `text-white/80` - Text colors with opacity
- `border-white/10` - Border colors with opacity

**Spacing & typography:**
- **Spacing:** `gap-2` (8px), `gap-3` (12px), `p-2`, `p-3`, `p-6`
- **Typography:** `text-xs` (12px), `text-sm` (14px), `text-base` (16px)
- **Line Height:** `leading-5` (20px), `leading-6` (24px)
- **Font Weight:** `font-normal` (400), `font-medium` (500), `font-semibold` (600)
- **Border Radius:** `rounded` (4px), `rounded-lg` (8px), `rounded-full`

**Component styling pattern:**
```tsx
"use client";

type ComponentProps = React.HTMLAttributes<HTMLDivElement> & {
  size?: "sm" | "base" | "lg";
  variant?: "primary" | "secondary";
};

export default function Component({
  size = "base",
  variant = "primary",
  className = "",
  ...rest
}: ComponentProps) {
  const baseClasses = "flex items-center rounded transition-colors";

  const sizeConfig = {
    sm: "gap-2 p-2 text-sm leading-5",
    base: "gap-2 p-2 text-base leading-6",
    lg: "gap-3 p-3 text-base leading-6",
  };

  const variantConfig = {
    primary: "bg-white text-black hover:bg-white/90 disabled:opacity-60",
    secondary: "bg-white/8 border border-white/10 text-white hover:bg-white/12",
  };

  return (
    <div
      className={[baseClasses, sizeConfig[size], variantConfig[variant], className]
        .filter(Boolean)
        .join(" ")}
      {...rest}
    />
  );
}
```

**DO ✅:**
- Use TypeScript to type all props and config objects
- Use pure Tailwind only - no CSS variables or custom classes
- Use opacity modifiers for color variants
- Use pseudo-class modifiers (`hover:`, `focus:`, `disabled:`)
- Support `className` prop and `{...rest}` spread for customization
- Use array filtering pattern for conditional classes

**DON'T ❌:**
- Don't create CSS variables - use Tailwind exclusively
- Don't use inline styles
- Don't manage focus/hover in JavaScript
- Don't hardcode pixel values - use Tailwind's spacing scale
- Don't use external class name libraries

### Icon System

**SVG icons stored in `components/svgs/` as React components:**
```tsx
// components/svgs/search.tsx
export default function SearchSvg(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" {...props}>
      <circle cx="10.5" cy="10.5" r="7.5" stroke="currentColor" strokeWidth="2" />
      <path d="M16 16L21 21" stroke="currentColor" strokeWidth="2" />
    </svg>
  );
}
```

**Key practices:**
- Use `currentColor` so icons inherit parent text color
- Support `{...props}` for size/className customization
- Register icons in `components/icon.tsx` for easy lookup

### Build Configuration

**Next.js config** (`next.config.ts`):
- Standalone output mode for Docker deployment
- Server actions with 5MB body size limit
- Console logs stripped in production (warn/error logs retained)
- Turbopack enabled for dev performance

**Environment variables** (`.env`):
- `NODE_ENV` - Environment (development/production)
- `GATEWAY_URL` - Backend API gateway base URL
- `NOTIFICATION_URL` - SSE notification service URL (server-side only, proxied via `/api/sse`)

**Path aliases** (`tsconfig.json`):
- `@/*` maps to project root for clean imports

### SSE Proxy Pattern

The `/api/sse` route proxies Server-Sent Events from the private notification service:
- Forwards cookies and origin headers to upstream
- Maintains persistent connection to client
- Streams notification messages directly to browser
- Handles graceful disconnects and error recovery
- Required because notification service is not publicly accessible

## Testing

No test infrastructure is currently present in the codebase.

**To add Go service tests:**
```bash
# Create test file alongside source (same package)
# internal/handler_test.go, internal/repository_test.go, etc.

# Run all tests in a service
cd <service-name>
go test ./...

# Run specific test with verbose output
go test -v -run TestHandlerCreate ./internal

# Run with coverage
go test -cover ./...
```

**Test file template:**
```go
// internal/handler_test.go
package internal

import (
    "testing"
)

func TestCreateResource(t *testing.T) {
    // Arrange: Setup test fixtures
    handler := NewHandler(mockRepo)

    // Act: Execute test
    result := handler.CreateResource(mockContext)

    // Assert: Verify results
    if result == nil {
        t.Fatal("expected result, got nil")
    }
}
```

**For frontend tests:** See `web/CLAUDE.md` for Next.js testing setup.

## Development Workflow

### Local Development Setup

**Quick start (all services in Docker):**
```bash
# 1. Clone repository
git clone <repo>
cd instrlabs-apps

# 2. Create .env with your configuration (see .env.example)
# Add database credentials, API keys, OAuth secrets, etc.
cp .env.example .env
nano .env  # Edit with your secrets

# 3. Start all services
docker-compose up --build

# Services available at:
#   Gateway: http://localhost:3000
#   Auth: http://localhost:3001
#   Image: http://localhost:3002
#   Notification: http://localhost:3003
#   Web: http://localhost:8000
```

**For local service development (without Docker):**
```bash
# 1. Create .env with your configuration
cp .env.example .env
nano .env  # Edit with your secrets

# 2. Ensure dependencies are running:
# - MongoDB: mongod (local or docker)
# - NATS: nats-server (local or docker)
# - MinIO (for image-service): minio server /data

# 3. Run each service in separate terminals:
# All services read from the root .env file

# Terminal 1: Auth service
cd auth-service && go run main.go

# Terminal 2: Image service
cd image-service && go run main.go

# Terminal 3: Notification service
cd notification-service && go run main.go

# Terminal 4: Gateway service
cd gateway-service && go run main.go

# Terminal 5: Web application
cd web && npm install && npm run dev
```

**Understanding .env file:**
```bash
# Single root .env file for entire project
# All services (backend and frontend) read from this file
# No service-level .env files - configuration is centralized

# Setup:
cp .env.example .env    # Start with template
nano .env               # Add your secrets:
                        # - Real MongoDB connection
                        # - S3 credentials
                        # - OAuth secrets
                        # - Email credentials

# IMPORTANT: .env is in .gitignore - never commit it!
```

### Adding a New Microservice
1. Copy structure from existing service (e.g., `auth-service`)
2. Update `go.mod` module name to `github.com/instrlabs/<new-service>`
3. Create `.env` file from `.env.example`
4. Add service configuration to `docker-compose.yaml`
5. Register routes in `gateway-service/internal/config.go` `Services` array
6. Update gateway `.env` with new service URL
7. Run `docker-compose up --build` to test

### Adding New NATS Message Subjects
1. Define constants in each service's `internal/config.go`:
   ```go
   type Config struct {
       // ... existing fields
       NatsSubjectNewTopic string `envconfig:"NATS_SUBJECT_NEW_TOPIC"`
   }
   ```
2. Add environment variables to `.env.example` and all service `.env` files
3. Implement publisher in one service:
   ```go
   h.nats.Conn.Publish(h.cfg.NatsSubjectNewTopic, payload)
   ```
4. Implement subscriber in consuming service (main.go):
   ```go
   nats.Conn.Subscribe(cfg.NatsSubjectNewTopic, func(m *natsgo.Msg) {
       handler.ProcessMessage(m.Data)
   })
   ```

### Updating the Shared Module
When changes are made to `github.com/instrlabs/shared`:
```bash
# Update in all services
cd auth-service && go get github.com/instrlabs/shared@v0.0.X && go mod tidy && cd ..
cd gateway-service && go get github.com/instrlabs/shared@v0.0.X && go mod tidy && cd ..
cd image-service && go get github.com/instrlabs/shared@v0.0.X && go mod tidy && cd ..
cd notification-service && go get github.com/instrlabs/shared@v0.0.X && go mod tidy && cd ..

# Or use a loop
for dir in */; do
    cd "$dir" && go get github.com/instrlabs/shared@v0.0.X && go mod tidy && cd ..
done
```

### Common Issues and Solutions

**Gateway can't reach service:**
- Verify service is running on correct port
- Check gateway `.env` SERVICE_* URLs match docker-compose service names
- For Docker: use `http://service-name:3000`, for localhost: use `http://localhost:3001`

**Authentication fails:**
- Verify `JWT_SECRET` is identical in gateway and auth service `.env`
- Check cookies are being sent (browser DevTools → Application → Cookies)
- Verify token hasn't expired (`access_token` expires, triggers refresh flow)

**MongoDB connection errors:**
- Verify `MONGO_URI` is correct (default: `mongodb://localhost:27017`)
- Check MongoDB is running (`docker ps` or `mongod` locally)
- For Docker Compose: use `mongodb://mongo:27017` (MongoDB service name)

**S3/MinIO connection errors:**
- Verify S3 credentials in image-service `.env`
- Check bucket exists: `aws s3 ls s3://bucket-name` or MinIO console
- For local MinIO: ensure it's running on correct endpoint