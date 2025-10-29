# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

**Development server:**
```bash
npm run dev
# Runs Next.js dev server with Turbopack on port 8000
```

**Production build:**
```bash
npm run build
npm start  # Runs production server on port 8000
```

**Linting and formatting:**
```bash
npm run lint          # Run ESLint
npm run format        # Format code with Prettier
npm run format:check  # Check formatting without writing
```

## Architecture Overview

### Tech Stack
- **Framework:** Next.js 15.4.7 with App Router (React 19)
- **Styling:** Tailwind CSS v4 with Geist font
- **Type Safety:** TypeScript with strict mode
- **Forms:** react-hook-form

### Route Groups and Layout Structure

The application uses Next.js route groups to organize pages with different layouts:

1. **`(non-auth)/`** - Public pages without authentication (e.g., login, register)
2. **`(site)/`** - Authenticated pages with full layout including navigation, overlays, and context providers
3. **`debug/`** - Development/testing pages
4. **`api/`** - API routes (e.g., SSE proxy at `/api/sse`)

The authenticated site layout (`app/(site)/layout.tsx`) wraps all authenticated pages in a nested provider hierarchy:
```
ProfileProvider → ProductProvider → SSEProvider → NotificationProvider →
ModalProvider → OverlayProvider → OverlayTop + OverlayContent + NotificationWidget
```

### Authentication Flow

**Middleware-based authentication** (`middleware.ts:5-69`):
- All routes except whitelisted paths (`/login`, `/register`, `/forgot-password`, `/reset-password`, `/`) require authentication
- Automatically attempts token refresh if `access_token` is missing but `refresh_token` exists
- Forwards client metadata headers (`x-user-ip`, `x-user-agent`, `x-user-host`, `x-user-origin`) to the gateway
- Redirects to `/login` on authentication failure

**Token management:**
- `access_token` and `refresh_token` stored in HTTP-only cookies
- Tokens automatically set by `fetchPOST()` when calling `/auth/login` or `/auth/refresh` endpoints
- Auth functions in `services/auth.ts` must only be called from server components/actions

### API Integration Pattern

**Server-side fetch utilities** (`utils/fetch.ts`):
- All API calls go through `GATEWAY_URL` environment variable
- `fetchGET()`, `fetchPOST()`, `fetchPUT()`, `fetchPATCH()` - JSON endpoints
- `fetchGETBytes()` - Binary data (e.g., images)
- `fetchPOSTFormData()` - Form uploads with multipart/form-data

**Response structure:**
```typescript
{
  success: boolean;
  message: string;
  data: T | null;
  errors: FormErrors | null;
}
```

**Important:** All fetch utilities are server actions (`"use server"`). They automatically:
- Forward client metadata headers from middleware
- Handle cookie-based authentication
- Set tokens on successful login/refresh

### Overlay System

**Centralized overlay management** (`hooks/useOverlay.tsx`):
- Overlays slide in from left or right sides of the screen
- Registration-based system with key-based lookup (e.g., `"left:navigation"`, `"right:notifications"`, `"right:profile"`)
- Actions: `openLeft(key)`, `closeLeft()`, `openRight(key)`, `closeRight()`
- Register new overlays in `registerOverlays()` function

### Component Organization

Components are organized by function in subdirectories:
- `components/actions/` - Interactive buttons and actions
- `components/cards/` - Card layouts and containers
- `components/feedback/` - Loading states, alerts, toasts
- `components/icons/` - SVG icon components
- `components/inputs/` - Form inputs and controls
- `components/layouts/` - Page layout components (OverlayTop, OverlayContent)
- `components/navigation/` - Navigation components
- `components/overlays/` - Overlay content (navigation, notifications, profile)

### Custom Hooks

- `useOverlay` - Control left/right sliding overlays
- `useModal` - Modal dialog state management
- `useNotification` - Toast/notification system
- `useSSE` - Server-Sent Events connection management
- `useProduct` - Access product list from context
- `useProfile` - Access user profile from context
- `useMediaQuery` - Responsive breakpoint detection

### Environment Variables

Required environment variables (see `.env.example`):
- `NODE_ENV` - Environment (development/production)
- `GATEWAY_URL` - Backend API gateway base URL
- `NOTIFICATION_URL` - SSE notification service URL (server-side only)

### Build Configuration

**Next.js config** (`next.config.ts`):
- Standalone output mode for Docker deployment
- Server actions with 5MB body size limit
- Console logs stripped in production (except warn/error)

**Path aliases:**
- `@/*` maps to project root (configured in `tsconfig.json`)

### SSE Proxy Pattern

The `/api/sse` route proxies Server-Sent Events from the private `NOTIFICATION_URL` service:
- Forwards cookies and origin headers to upstream
- Streams responses directly to client
- Handles connection errors and client disconnects
- Required because notification service is not publicly accessible

### Styling Guidelines

- Use native array filtering for conditional classes instead of external libraries
- Tailwind CSS with automatic class sorting via Prettier
- Geist font family applied globally
- No inline styles unless absolutely necessary


## Figma Design System Integration

### Core Principle

**Use Tailwind's built-in utility classes exclusively.** No CSS custom properties, no `@utility` definitions, no custom CSS classes.

### Color System

Use Tailwind's color palette with opacity modifiers:

**Base Colors:**
- `bg-black`, `bg-white` - Black/white backgrounds
- `text-white`, `text-black` - White/black text
- `border-white` - White borders

**Opacity Modifiers (using `/` syntax):**
- `bg-white/90` - 90% opacity
- `bg-white/8` - 8% opacity
- `text-white/80` - 80% opacity
- `border-white/10` - 10% opacity

**Common Patterns:**
```tsx
// Primary button
<button className="bg-white text-black hover:bg-white/90 disabled:opacity-60">

// Secondary button
<button className="bg-white/8 border border-white/10 text-white hover:bg-white/12">

// Transparent variant
<button className="text-white hover:bg-white/8">
```

### Spacing & Typography

Use Tailwind's built-in scale:

- **Spacing:** `gap-2` (8px), `gap-3` (12px), `p-2`, `p-3`, `p-6`
- **Typography:** `text-sm` (14px), `text-base` (16px), `text-xs` (12px)
- **Line Height:** `leading-5` (20px), `leading-6` (24px)
- **Font Weight:** `font-normal` (400), `font-medium` (500), `font-semibold` (600)
- **Border Radius:** `rounded` (4px), `rounded-lg` (8px), `rounded-full`

### Component Architecture

**Core Principles:**

1. **Use `"use client"` Directive** - All interactive components require this
2. **Component Naming** - Use PascalCase (e.g., `Button`, `InputField`)
3. **Type-Safe Props** - Extend native HTML element types with TypeScript
4. **Organize Classes by Purpose:**
   - Base classes: Layout, structure, common styles
   - Size config: Spacing, typography variations
   - Variant config: Color schemes, visual styles
   - State modifiers: Hover, focus, disabled states

**Configuration Objects Pattern:**

```tsx
"use client";

import React from "react";

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

  const sizeConfig: Record<"sm" | "base" | "lg", string> = {
    sm: "gap-2 p-2 text-sm leading-5",
    base: "gap-2 p-2 text-base leading-6",
    lg: "gap-3 p-3 text-base leading-6",
  };

  const variantConfig: Record<"primary" | "secondary", string> = {
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

### Icon System

Icons are stored as React SVG components in **`components/svgs/`**:

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

Use `currentColor` to make icons inherit text color. Icons are registered in `components/icon.tsx`.

### Color Mapping Reference

| Figma Color | Opacity | Tailwind Class |
|-------------|---------|----------------|
| White | 100% | `bg-white`, `text-white` |
| White | 90% | `bg-white/90`, `text-white/90` |
| White | 8% | `bg-white/8` |
| White | 4% | `bg-white/4` |
| White | 10% | `border-white/10` |
| Black | 100% | `bg-black`, `text-black` |

### Best Practices

#### DO ✅
- Use TypeScript to type all props and config objects
- Use Pure Tailwind only - no CSS variables or custom classes
- Use opacity modifiers for color variants (`bg-white/90`, `text-white/50`)
- Use pseudo-class modifiers (`hover:`, `focus:`, `focus-within:`, `disabled:`)
- Support `className` prop and `{...rest}` spread for customization
- Use array filtering pattern for conditional classes
- Type config objects with `Record<>`

#### DON'T ❌
- Don't create CSS variables - use Tailwind exclusively
- Don't use inline styles
- Don't manage focus/hover in JavaScript
- Don't hardcode pixel values - use Tailwind's spacing scale
- Don't use external class name libraries

### Asset Management

- **DO NOT download Figma assets** - Use provided URLs during development
- **Convert to SVG components** - Create React SVG components in `components/svgs/`
- **Use currentColor** - Make SVG icons inherit text color for styling flexibility

