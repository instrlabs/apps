# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.


## Services

**Gateway Service** (`gateway-service/`)
- Handle proxies requests to downstream services
- Handle Middlewares (CORS, Rate Limit, CSRF Protection, and Authentication)

**Auth Service** (`auth-service/`)
- Handle Authentication and Authorization
- Handle OAuth & JWT Flows
- Session-Based Token Binding with device validation (IP + User-Agent hash)
- Support multiple concurrent sessions with per-device revocation

**Image Service** (`image-service/`)
- Handle Image processing requests (via NATS, API, Scheduling)

**Notification Service** (`notification-service/`)
- Handle Server-Sent Events (SSE) for real-time notifications to web clients

**Web Application** (`web/`)
- Next.js 15 with App Router and React 19

## Environment Configuration

Environment variables are separated into two categories for clarity:

### Public Exposed URLs (`*_URL`)
These are services accessible from the frontend and external clients:
- `WEB_URL` - Next.js web application (e.g., `http://localhost:3000`)
- `API_URL` - Backend API gateway (e.g., `http://localhost:3001`)
- `NOTIFICATION_URL` - SSE notification service (e.g., `http://localhost:3002`)

**Used by:** Web frontend and browser-based clients

### Internal Service URLs (`*_SERVICE`)
These are services that can only be accessed within the Docker network:
- `AUTH_SERVICE` - Authentication service (e.g., `http://auth-service:3000`)
- `IMAGE_SERVICE` - Image processing service (e.g., `http://image-service:3000`)
- `GATEWAY_SERVICE` - Backend API gateway service (e.g., `http://gateway-service:3000`)
- `NOTIFICATION_SERVICE` - SSE notification service (e.g., `http://notification-service:3000`)

**Used by:** Backend services (Gateway, Notification) for internal communication

### Configuration Files
- **Root `.env`** - Contains all public and internal URL definitions
- **Service `.env` files** - Reference root variables using `${VAR_NAME}` syntax
  - `gateway-service/.env` - Uses `AUTH_SERVICE` and `IMAGE_SERVICE`
  - `auth-service/.env` - Uses `API_URL` and `WEB_URL`
  - `notification-service/.env` - Uses only public `ORIGINS_ALLOWED`
  - `web/.env.local` - Uses `API_URL` and `NOTIFICATION_URL`

See `.env.example` for complete configuration template.

## Code Patterns

### Service Structure
```
<service>/
├── Dockerfile
├── go.mod
├── main.go          # Entry point with setup
├── internal/
│   ├── config.go    # Environment config loader
│   ├── *_handler.go # HTTP handlers (Fiber)
│   ├── *_repository.go # Database access
│   └── *.go         # Domain models
└── static/          # Static assets (if any)
```

### Origin Header Pattern (`x-user-origin`)

**IMPORTANT:** All backend services use `x-user-origin` for origin validation, NOT the standard `Origin` header.

**Key Points:**
- Frontend code NEVER manually sets `x-user-origin` - it's automatic via Next.js middleware
- All Go services check `c.Get("x-user-origin")` for origin validation
- For testing, manually set `x-user-origin` to simulate production flow
- Allowed origins configured in `gateway-service/.env` under `ORIGINS_ALLOWED`

## Guidelines
- Always get environment configuration from `.env`
- If Testing always refers to the .ai-guides directory in the service folder
- **Use `x-user-origin` header for all origin validation in backend services**
- **All test scripts must be created in `/tmp/` directory** - Keep test files temporary and out of version control

## AI Workflow Documentation Standards

**IMPORTANT:** All complex multi-step workflows MUST follow this documentation pattern to maintain context continuity.

### Documentation Structure

For any significant implementation task, create a master planning document:

**Format:** `.tmp/{JOB_NAME}_PLANS.md`

**Example:** `.tmp/AUTHENTICATION_REFACTOR_PLANS.md`, `.tmp/API_MIGRATION_PLANS.md`

### Master Document Requirements

The master planning document MUST contain:

1. **Project Header**
   - Task name and date
   - Current status
   - Last update timestamp

2. **Overview Section**
   - Business context/requirements
   - Key information from user
   - High-level architecture

3. **File Impact Map**
   - List all files to be modified/created
   - Current vs. target state
   - Dependency relationships

4. **Implementation Cycles**
   - One section per cycle
   - Format: `## Cycle N: {Description}`
   - For each cycle include:
     - File(s) affected
     - What to change (exact line numbers)
     - Implementation details
     - **Status field: [PENDING], [IN_PROGRESS], [COMPLETED]**
     - Build verification step

5. **Progress Tracking Table**
   - Simple table showing all cycles
   - Update status after each cycle completion
   - Include lines added/modified count

6. **Context Management Section**
   - When to read this document again
   - Key assumptions to remember
   - Important links/references

### Cycle Progress Updates

**AFTER COMPLETING EACH CYCLE:**

1. Update the `Status` field for that cycle: `[COMPLETED]`
2. Add completion timestamp
3. Update the progress tracking table
4. Document any issues encountered
5. Note any context discovered for next cycles

### Context Token Management

**When approaching maximum context tokens:**

1. **STOP** - Do not continue with implementation
2. **READ** - Re-read the master planning document completely
3. **VERIFY** - Check which cycles are completed vs. pending
4. **RESUME** - Continue from next pending cycle
5. **REFERENCE** - Always cite the document when resuming

**Pattern to resume work:**
```
Reading from .tmp/{JOB_NAME}_PLANS.md:
- Completed cycles: [list]
- Current cycle: [number and description]
- Next action: [specific task]
```

### Documentation Examples

See `.tmp/` directory for examples:
- `REVISED_IMPLEMENTATION_PLAN.md` - Master plan structure
- `CYCLE_*.md` - Individual cycle details
- `IMPLEMENTATION_COMPLETE.md` - Final summary

### Benefits

✅ No context loss between sessions
✅ Easy to resume interrupted work
✅ Clear progress tracking
✅ Documentation serves as reference
✅ Reduces mistakes and rework
✅ Helpful for code reviews

### Checklist for Documentation

- [ ] Create `.tmp/{JOB_NAME}_PLANS.md` at start of work
- [ ] Document all files to be changed
- [ ] Break implementation into numbered cycles
- [ ] Add status field to each cycle
- [ ] Update status after completing each cycle
- [ ] Add timestamp for each update
- [ ] Include verification step for each cycle
- [ ] Create summary document at end
- [ ] Reference document when resuming work