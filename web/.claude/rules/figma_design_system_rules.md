# Figma Design System Integration Rules

This document provides comprehensive guidelines for integrating Figma designs into the InstrLabs codebase using the Model Context Protocol (MCP).

## 1. Styling with Tailwind CSS

### Core Principle

**Use Tailwind's built-in utility classes exclusively.** No CSS custom properties, no `@utility` definitions, no custom CSS classes.

### Color System

Use Tailwind's color palette with opacity modifiers:

**Base Colors:**
- `bg-black` - Black background
- `bg-white` - White background
- `text-white` - White text
- `text-black` - Black text
- `border-white` - White border

**Opacity Modifiers (using `/` syntax):**
- `bg-white/90` - White background at 90% opacity (rgba(255, 255, 255, 0.9))
- `bg-white/10` - White background at 10% opacity (rgba(255, 255, 255, 0.1))
- `text-white/80` - White text at 80% opacity
- `text-white/50` - White text at 50% opacity
- `border-white/10` - White border at 10% opacity

**Common Patterns:**
```tsx
// Primary button
<button className="bg-white/90 text-black hover:bg-white disabled:opacity-60">

// Secondary button
<button className="bg-white/4 text-white border border-white/10 hover:bg-white/8">

// Input field
<div className="bg-white/4 border border-white/10 text-white/50 focus-within:text-white">

// Text variants
<p className="text-white">Primary text</p>
<p className="text-white/80">Secondary text</p>
<p className="text-white/30">Muted text</p>
```

### Spacing & Typography

Use Tailwind's built-in scale:

- **Spacing:** `gap-2` (8px), `gap-3` (12px), `p-2`, `p-3`, `p-6`
- **Typography:** `text-sm` (14px), `text-base` (16px), `text-xs` (12px)
- **Line Height:** `leading-5` (20px), `leading-6` (24px)
- **Font Weight:** `font-normal` (400), `font-medium` (500), `font-semibold` (600), `font-light` (300)
- **Border Radius:** `rounded` (4px), `rounded-lg` (8px), `rounded-full`

### Pseudo-Class Modifiers

Use Tailwind's state modifiers:

- `hover:` - Hover state
- `focus:` - Focus state
- `focus-within:` - Focus-within state (for parent containers)
- `active:` - Active state
- `disabled:` - Disabled state
- `has-[selector]:` - Parent state based on child (e.g., `has-[input:disabled]:opacity-50`)
- `group-hover:` - Hover on parent with `group` class

[Full list →](https://tailwindcss.com/docs/hover-focus-and-other-states)
---

## 2. Component Library

### Location

Components are organized in **`components/`** directory

### Component Architecture

**Core Principles:**

1. **Use `"use client"` Directive**
   - All interactive components must start with `"use client"`
   - Required for components using hooks, event handlers, or state

2. **Component Naming**
   - Use PascalCase (e.g., `Button`, `InputField`, `FileDropzone`)
   - File names should match component names: `button.tsx`, `input-field.tsx`

3. **Type-Safe Props**
   - Extend native HTML element types
   - Define props interface with TypeScript
   - Use `React.ButtonHTMLAttributes`, `React.InputHTMLAttributes`, etc.
   - Support `className` prop for customization
   - Support `{...rest}` spread for native attributes

4. **Organize Classes by Purpose**
   - **Base classes**: Layout, structure, common styles
   - **Size config**: Spacing, typography variations
   - **Variant config**: Color schemes, visual styles
   - **State modifiers**: Hover, focus, disabled states

5. **Configuration Objects**
   - Use typed `Record<>` objects for variants
   - Keep class strings readable and organized
   - Group related utilities together

**Example Structure:**

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
  // Base classes - structure and layout
  const baseClasses = "flex items-center rounded transition-colors";

  // Size configuration - spacing and typography
  const sizeConfig: Record<"sm" | "base" | "lg", string> = {
    sm: "gap-2 p-2 text-sm leading-5",
    base: "gap-2 p-2 text-base leading-6",
    lg: "gap-3 p-3 text-base leading-6",
  };

  // Variant configuration - colors and states
  const variantConfig: Record<"primary" | "secondary", string> = {
    primary: "bg-white/90 text-black hover:bg-white disabled:opacity-60",
    secondary: "bg-white/4 text-white border border-white/10 hover:bg-white/8",
  };

  return (
    <div
      className={[
        baseClasses,
        sizeConfig[size],
        variantConfig[variant],
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      {...rest}
    />
  );
}
```

---

## 3. Frameworks & Libraries

### UI Framework
- **React 19** with Next.js 15.4.7 App Router
- All components are React Server Components by default
- Use `"use client"` for interactive components

### Styling
- **Tailwind CSS v4** - Utility-first CSS framework exclusively
- **No custom CSS** - Pure Tailwind classes only
- **Opacity modifiers** - Use `/` syntax for color variations (e.g., `bg-white/10`)

### Type Safety
- **TypeScript** with strict mode enabled
- Type all component props and configuration objects

### Forms
- **react-hook-form** - For form state management and validation

### Build System
- **Next.js** with Turbopack (dev)
- **Standalone output mode** for Docker deployment

---

## 4. Asset Management

### Best Practices
- **DO NOT download Figma assets** - Use provided URLs during development
- **Convert to SVG components** - For icons, create React SVG components in `components/svgs/`
- **Use currentColor** - Make SVG icons inherit text color:

```tsx
<svg viewBox="0 0 24 24" fill="none" {...props}>
  <circle stroke="currentColor" strokeWidth="2" />
</svg>
```

### Asset Optimization
- SVGs are inlined as React components
- No CDN configuration needed
- Images use Next.js Image component when appropriate

---

## 5. Icon System

### Location
Icons are stored as React SVG components in **`components/svgs/`**

### Structure
Each icon is a separate file exporting a React component:

**`components/svgs/search.tsx`**
```tsx
import React from "react";

export default function SearchSvg(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" {...props}>
      <circle cx="10.5" cy="10.5" r="7.5" stroke="currentColor" strokeWidth="2" />
      <path d="M16 16L21 21" stroke="currentColor" strokeWidth="2" />
    </svg>
  );
}
```

### Icon Component (`components/icon.tsx`)

Central icon registry for easy icon usage:

```tsx
import SearchSvg from "./svgs/search";
import VisibleSvg from "./svgs/visible";
// ... more imports

const iconMap: Record<string, React.ComponentType<React.SVGProps<SVGSVGElement>>> = {
  search: SearchSvg,
  visible: VisibleSvg,
  // ... more icons
};

export default function Icon({ name, size = 24, className }: IconProps) {
  const IconComponent = iconMap[name];
  return <IconComponent width={size} height={size} className={className} />;
}
```

### Usage
```tsx
<Icon name="search" size={24} />
<Icon name="search" size={24} className="text-white/50 hover:text-white transition-colors" />
```

**Note:** Icons inherit color from parent text color via `currentColor`. Use text color utilities to style icons.

### Naming Convention
- **File names:** kebab-case (e.g., `circle-success.tsx`)
- **Component names:** PascalCase with "Svg" suffix (e.g., `CircleSuccessSvg`)
- **Icon registry keys:** kebab-case (e.g., `"circle-success"`)

### Adding New Icons from Figma

1. Export icon as SVG from Figma (or use MCP URL temporarily)
2. Create new file in `components/svgs/{icon-name}.tsx`
3. Convert SVG to React component with `props: React.SVGProps<SVGSVGElement>`
4. Use `currentColor` for stroke/fill to inherit text color
5. Register in `components/icon.tsx` iconMap
6. Import and add to iconMap object

---

## 6. Styling Approach

### Methodology
**Pure Tailwind CSS v4**

### Core Principles

1. **Use Tailwind Exclusively:**
   - Colors: `bg-white`, `text-black`, `bg-white/10`
   - Spacing: `gap-2`, `p-3`, `mt-4`
   - Layout: `flex`, `grid`, `items-center`
   - Typography: `text-sm`, `font-medium`, `leading-6`
   - Border: `rounded`, `border`, `border-white/10`

2. **Opacity Modifiers for Variants:**
   - Use `/` syntax: `bg-white/90`, `bg-white/10`, `text-white/80`
   - No CSS variables needed
   - All color variations handled by Tailwind

3. **Pseudo-Classes with Modifiers:**
   - Use Tailwind pseudo-class modifiers (`hover:`, `focus:`, `focus-within:`, `disabled:`)
   - Avoid JavaScript state for focus/hover when possible

### Example Patterns

**Button Component:**
```tsx
// Primary button
const primaryClasses = [
  "bg-white/90",
  "text-black",
  "hover:bg-white",
  "disabled:opacity-60",
  "disabled:cursor-not-allowed",
  "transition-colors",
  "gap-2",
  "p-2",
  "rounded",
].join(" ");

// Secondary button
const secondaryClasses = [
  "bg-white/4",
  "text-white",
  "border",
  "border-white/10",
  "hover:bg-white/8",
  "disabled:opacity-60",
  "disabled:cursor-not-allowed",
  "transition-colors",
  "gap-2",
  "p-2",
  "rounded",
].join(" ");
```

**Input Component:**
```tsx
const inputClasses = [
  "bg-white/4",
  "border",
  "border-white/10",
  "text-white/50",
  "focus-within:text-white",
  "has-[input:disabled]:opacity-50",
  "transition-colors",
  "gap-2",
  "p-2",
  "rounded",
].join(" ");
```

### Class Name Composition

Use array filtering pattern for conditional classes:

```tsx
const className = [
  "base-class",
  "always-applied",
  condition && "conditional-class",
  anotherCondition ? "class-a" : "class-b",
  customClassName,
]
  .filter(Boolean)
  .join(" ");
```

### Global Styles

Located in `app/globals.css`:
- Custom animations (`@keyframes` for notifications, modals, etc.)
- Global body styles
- Font imports (Geist font family)

### Responsive Design

Use Tailwind breakpoints:
- `sm:` - 640px
- `md:` - 768px
- `lg:` - 1024px
- `xl:` - 1280px

Example:
```tsx
<div className="flex-col md:flex-row gap-4 md:gap-6">
```

---

## 7. Project Structure

### Directory Organization

```
web/
├── app/                     # Next.js App Router
│   ├── (non-auth)/          # Public pages (login, register)
│   ├── (site)/              # Authenticated pages
│   ├── api/                 # API routes
│   ├── debug/               # Development pages
│   ├── globals.css          # Global styles & design tokens
│   └── layout.tsx           # Root layout
├── components/              # React components (organized by type)
│   ├── button.tsx           # Primary button
│   ├── input.tsx            # Primary input
│   ├── svgs/                # SVG icon components
│   └── ...
├── hooks/                    # Custom React hooks
├── services/                 # API service functions
├── utils/                    # Utility functions
├── middleware.ts             # Auth & routing middleware
├── CLAUDE.md                # Claude Code guidance
├── .cursor/
│   └── rules/
│       └── figma_design_system_rules.md  # This file
└── package.json
```

### Route Groups

Next.js route groups organize pages with different layouts:

1. **`(non-auth)/`** - Public pages without auth (login, register)
2. **`(site)/`** - Authenticated pages with full layout
3. **`debug/`** - Development/testing pages
4. **`api/`** - API routes

### Component Feature Organization

**By Type, Not By Feature:**
- ✅ `components/inputs/input.tsx`
- ✅ `components/cards/apps-card.tsx`
- ❌ `components/login/LoginInput.tsx`

---

## 8. Figma to Code Conversion Workflow

### Step-by-Step Process

#### 1. **Analyze Figma Design**
   - Identify component properties (size, color, state variants)
   - Note spacing values (8px = gap-2, 12px = gap-3)
   - Extract color values and convert to Tailwind classes
   - Check typography (14px = text-sm, 16px = text-base)
   - Map opacity values (80% = /80, 10% = /10)

#### 2. **Map Colors to Tailwind**

Convert Figma colors to Tailwind classes with opacity modifiers:

| Figma Color | Tailwind Class |
|-------------|----------------|
| `rgba(255, 255, 255, 1)` | `bg-white` or `text-white` |
| `rgba(255, 255, 255, 0.9)` | `bg-white/90` or `text-white/90` |
| `rgba(255, 255, 255, 0.8)` | `bg-white/80` or `text-white/80` |
| `rgba(255, 255, 255, 0.5)` | `bg-white/50` or `text-white/50` |
| `rgba(255, 255, 255, 0.1)` | `bg-white/10` or `text-white/10` |
| `rgba(255, 255, 255, 0.04)` | `bg-white/4` or `text-white/4` |
| `rgba(0, 0, 0, 1)` | `bg-black` or `text-black` |
| `rgba(0, 0, 0, 0.6)` | `bg-black/60` or `text-black/60` |

#### 3. **Build Component**

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
  const sizeConfig: Record<"sm" | "base" | "lg", { spacing: string; font: string }> = {
    sm: { spacing: "gap-2 p-2", font: "text-sm leading-5" },
    base: { spacing: "gap-2 p-2", font: "text-base leading-6" },
    lg: { spacing: "gap-3 p-3", font: "text-base leading-6" },
  };

  const variantConfig: Record<"primary" | "secondary", string> = {
    primary: "bg-white/90 text-black hover:bg-white",
    secondary: "bg-white/4 text-white border border-white/10 hover:bg-white/8",
  };

  return (
    <div
      className={[
        "flex",
        "items-center",
        "rounded",
        "transition-colors",
        sizeConfig[size].spacing,
        sizeConfig[size].font,
        variantConfig[variant],
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      {...rest}
    />
  );
}
```

#### 4. **Handle Icons**

If component uses icons:
- Extract icon SVGs from Figma
- Create React SVG components in `components/svgs/`
- Register in `components/icon.tsx`
- Use `hasLeftIcon`/`hasRightIcon` props pattern

#### 5. **Update Imports**

Ensure components import from correct locations:
```tsx
import Button from "@/components/button";
import Input from "@/components/input";
import Icon from "@/components/icon";
```

---

## 9. Figma to Tailwind Mapping

### Color Mapping

| Figma Color | Opacity | Tailwind Class |
|-------------|---------|----------------|
| White (#FFFFFF) | 100% | `bg-white`, `text-white`, `border-white` |
| White | 90% | `bg-white/90`, `text-white/90` |
| White | 80% | `bg-white/80`, `text-white/80` |
| White | 50% | `bg-white/50`, `text-white/50` |
| White | 30% | `bg-white/30`, `text-white/30` |
| White | 10% | `bg-white/10`, `border-white/10` |
| White | 8% | `bg-white/8` |
| White | 4% | `bg-white/4` |
| Black (#000000) | 100% | `bg-black`, `text-black` |
| Black | 60% | `bg-black/60`, `text-black/60` |

### Spacing Mapping

| Figma Value | Tailwind Class |
|-------------|----------------|
| 8px | `gap-2`, `p-2`, `m-2` |
| 12px | `gap-3`, `p-3`, `m-3` |
| 16px | `gap-4`, `p-4`, `m-4` |
| 24px | `gap-6`, `p-6`, `m-6` |

### Typography Mapping

| Figma Size | Line Height | Weight | Tailwind Classes |
|-----------|-------------|---------|------------------|
| 12px | 16px | Regular | `text-xs leading-4 font-normal` |
| 14px | 20px | Regular | `text-sm leading-5 font-normal` |
| 14px | 20px | Medium | `text-sm leading-5 font-medium` |
| 16px | 24px | Regular | `text-base leading-6 font-normal` |
| 16px | 24px | Medium | `text-base leading-6 font-medium` |
| 16px | 24px | SemiBold | `text-base leading-6 font-semibold` |

### Border Radius Mapping

| Figma Value | Tailwind Class |
|-------------|----------------|
| 4px | `rounded` |
| 6px | `rounded-md` |
| 8px | `rounded-lg` |
| 999px / Full | `rounded-full` |

---

## 10. Best Practices

### DO ✅

- **Use TypeScript:** Type all props and config objects
- **Use Pure Tailwind:** Only use Tailwind's built-in utilities
- **Use Opacity Modifiers:** `bg-white/90`, `text-white/50` for color variants
- **Use Pseudo-Class Modifiers:** `hover:`, `focus:`, `focus-within:`, `disabled:` etc.
- **Spread Props:** Support `{...rest}` for native HTML attributes
- **Array Filtering:** Use `.filter(Boolean).join(" ")` for className composition
- **Type Config Objects:** Use `Record<>` for size/variant configurations
- **Name Consistently:** Follow existing naming patterns
- **Organize by Type:** Place components in appropriate subdirectories

### DON'T ❌

- **Don't create CSS variables** - Use Tailwind classes exclusively
- **Don't create @utility classes** - Use Tailwind utilities only
- **Don't use inline styles** - Use className composition
- **Don't manage focus/hover in JS** - Use Tailwind modifiers
- **Don't hardcode pixel values** - Use Tailwind's spacing scale
- **Don't nest deeply** - Keep component structure flat
- **Don't use external class name libraries** - Use array filtering pattern

---

## 11. Common Patterns

### Size Configuration Pattern

```tsx
const sizeConfig: Record<
  "sm" | "base" | "lg",
  { spacing: string; font: string; iconSize: number }
> = {
  sm: { spacing: "gap-2 p-2", font: "text-sm", iconSize: 20 },
  base: { spacing: "gap-2 p-2", font: "text-base", iconSize: 24 },
  lg: { spacing: "gap-3 p-3", font: "text-base", iconSize: 24 },
};
```

### Variant Configuration Pattern

```tsx
const variantConfig: Record<"primary" | "secondary", string> = {
  primary: "bg-white/90 text-black hover:bg-white disabled:opacity-60",
  secondary: "bg-white/4 text-white border border-white/10 hover:bg-white/8 disabled:opacity-60",
};
```

### Class Composition Pattern

```tsx
const className = [
  // Base styles
  "flex",
  "items-center",
  "rounded",
  "transition-colors",
  // Size-based styles
  sizeConfig[size].spacing,
  sizeConfig[size].font,
  // Variant styles
  variantConfig[variant],
  // Custom className prop
  customClassName,
]
  .filter(Boolean)
  .join(" ");
```

### Icon Rendering Pattern

```tsx
const renderIcon = (iconName: string | null) => {
  if (!iconName) return null;
  return (
    <span className="relative shrink-0">
      <Icon name={iconName} size={currentSize.iconSize} />
    </span>
  );
};
```

---

## 12. Quick Reference

### Figma Design → Code Checklist

- [ ] Extract color values → Map to Tailwind classes with opacity modifiers
- [ ] Map spacing to Tailwind scale (8px=gap-2, 12px=gap-3, 16px=gap-4)
- [ ] Map typography to Tailwind classes (14px=text-sm, 16px=text-base)
- [ ] Map opacity values (90%=/90, 10%=/10, 4%=/4)
- [ ] Create component file in appropriate `components/` subdirectory
- [ ] Type all props with TypeScript
- [ ] Create typed size/variant config objects using `Record<>`
- [ ] Use array filtering pattern for className composition
- [ ] Add pseudo-class modifiers (`hover:`, `focus-within:`, `disabled:`)
- [ ] Add icons to `components/svgs/` and register in `icon.tsx`
- [ ] Use `{...rest}` spread for native HTML attributes
- [ ] Test all size/variant/state combinations

### File Template Locations

- **Button-like components:** `components/button.tsx`
- **Input-like components:** `components/input.tsx`
- **Card components:** `components/cards/`
- **SVG icons:** `components/svgs/`
- **Global styles:** `app/globals.css` (animations only)

---

## 13. Advanced Patterns

### Multi-Line Template Literal Pattern (FileDropzone)

For complex components with dynamic inline styles, use `useMemo` with template literals:

```tsx
const baseClass = useMemo(() => (
  `
  group cursor-pointer outline-none
  flex w-full flex-col items-center justify-center
  gap-2 p-6
  rounded-lg border border-dashed border-white/10
  bg-transparent

  transition-colors focus-visible:ring-2 focus-visible:ring-white/20
  ${isDragging ? "bg-white/8" : "hover:bg-white/5"}
  ${className || ""}
  `
), [isDragging, className]);
```

**When to use:**
- Complex dynamic styles based on state
- Readability over single-line arrays
- Grouping related utilities visually

**Important:** Format multi-line template literals with logical grouping (layout, spacing, colors, states).

### Overlay System Pattern

The codebase uses a registration-based overlay system (see `hooks/useOverlay.tsx:31-60`):

```tsx
// Register overlays with keys
function registerOverlays() {
  registerOverlay("left:navigation", {
    side: "left",
    render: () => <NavigationOverlay />,
  });

  registerOverlay("right:notifications", {
    side: "right",
    render: () => <NotificationOverlay />,
  });
}

// Use in components
const { openLeft, closeLeft, openRight, closeRight } = useOverlay();

// Open overlays
openLeft("left:navigation");
openRight("notifications"); // Note: "right:" prefix added automatically
```

**Key Pattern:** Side overlays slide in from left/right with registration-based content management.

### Context Provider Pattern

Components that need global state use context providers with custom hooks:

```tsx
// Provider
export function OverlayProvider({ children }: { children: React.ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);

  const value = useMemo(() => ({
    isOpen,
    open: () => setIsOpen(true),
    close: () => setIsOpen(false),
  }), [isOpen]);

  return (
    <OverlayContext.Provider value={value}>
      {children}
    </OverlayContext.Provider>
  );
}

// Hook
export function useOverlay() {
  const ctx = useContext(OverlayContext);
  if (!ctx) throw new Error('useOverlay must be used within OverlayProvider');
  return ctx;
}
```

### Accessibility Patterns

The codebase implements comprehensive accessibility:

```tsx
<div
  role="button"
  tabIndex={0}
  aria-label="Upload files. Allowed: .png, .jpg. Max size: 5MB."
  aria-describedby={helperId}
  onKeyDown={(e) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      openFileDialog();
    }
  }}
>
  {/* content */}
</div>

<span id={helperId} className="text-xs text-white/50">
  Total file size allowed is 5MB...
</span>
```

**Required attributes:**
- `role` for semantic roles
- `tabIndex={0}` for keyboard focus
- `aria-label` for screen readers
- `aria-describedby` for detailed descriptions
- `onKeyDown` handlers for Enter/Space activation

### Notification Hook Pattern

Use the notification system for user feedback:

```tsx
import useNotification from "@/hooks/useNotification";

function MyComponent() {
  const { showNotification } = useNotification();

  const handleAction = () => {
    // Show error
    showNotification({
      type: "error",
      message: "Invalid file type or file size."
    });

    // Show success
    showNotification({
      type: "success",
      message: "File uploaded successfully!"
    });
  };
}
```

**Available types:** `"success"`, `"error"`, `"info"`, `"warning"`

### Custom Animations

Define custom animations in `app/globals.css`:

```css
@keyframes notificationIn {
  from { transform: translateY(100%); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.animate-notification-in {
  animation: notificationIn 0.3s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}
```

**Existing animations:**
- `animate-notification-in` / `animate-notification-out` - Toast notifications
- `animate-modal-in` / `animate-modal-out` - Modal dialogs
- `animate-flash` - Flashing effect for loading states

### Avatar Color Buckets

The Avatar component uses a bucketing system for consistent user colors:

```tsx
function getBucket(name: string | undefined) {
  const firstChar = safeName[0]?.toLowerCase() ?? "u";
  const alphaIndex = firstChar.charCodeAt(0) - 97;
  return ((alphaIndex % 26) + 26) % 8;
}

const bgPaletteCls = [
  "bg-blue-500",
  "bg-green-500",
  "bg-red-500",
  "bg-yellow-500",
  "bg-purple-500",
  "bg-teal-500",
  "bg-orange-500",
  "bg-slate-500",
];
```

**Pattern:** Deterministic color assignment based on first letter for consistent user identification.

---

## 14. Hooks Reference

### Available Hooks

| Hook | Purpose | Usage |
|------|---------|-------|
| `useOverlay` | Control left/right sliding overlays | `const { openLeft, closeLeft } = useOverlay();` |
| `useModal` | Modal dialog state management | `const { showModal, closeModal } = useModal();` |
| `useNotification` | Toast/notification system | `const { showNotification } = useNotification();` |
| `useSSE` | Server-Sent Events connection | `const { isConnected } = useSSE();` |
| `useProduct` | Access product list from context | `const { products } = useProduct();` |
| `useProfile` | Access user profile from context | `const { profile } = useProfile();` |
| `useMediaQuery` | Responsive breakpoint detection | `const isMobile = useMediaQuery('(max-width: 768px)');` |

### Hook Integration in Components

When building Figma components that need global state:

```tsx
"use client";

import useNotification from "@/hooks/useNotification";
import useOverlay from "@/hooks/useOverlay";

export default function MyComponent() {
  const { showNotification } = useNotification();
  const { openRight } = useOverlay();

  const handleClick = () => {
    openRight("profile");
    showNotification({ type: "success", message: "Profile opened!" });
  };

  return <button onClick={handleClick}>Open Profile</button>;
}
```

---

## 15. Utility Functions

### Common Utilities

The codebase provides utility functions in `utils/`:

- **`bytesToString(bytes: number)`** - Convert bytes to human-readable format (e.g., "5MB")
- **`acceptsToExtensions(accepts: string[])`** - Convert MIME types to file extensions

**Example:**
```tsx
import { bytesToString } from "@/utils/bytesToString";
import { acceptsToExtensions } from "@/utils/acceptsToExtensions";

const maxSize = 5242880; // 5MB in bytes
const accepts = ["image/png", "image/jpeg"];

console.log(bytesToString(maxSize)); // "5MB"
console.log(acceptsToExtensions(accepts)); // [".png", ".jpg"]
```

---

## 16. Prettier Configuration

Code formatting is enforced via Prettier with Tailwind plugin:

```json
{
  "semi": true,
  "singleQuote": false,
  "trailingComma": "all",
  "printWidth": 100,
  "tabWidth": 2,
  "plugins": ["prettier-plugin-tailwindcss"]
}
```

**Important:**
- Always use double quotes (not single quotes)
- Include trailing commas
- 100 character line width
- Tailwind classes are automatically sorted

**Usage:**
```bash
npm run format        # Format all files
npm run format:check  # Check formatting without writing
```

---

## 17. Component State Management

### Local State Pattern

Use `useState` for simple local state:

```tsx
const [isDragging, setIsDragging] = useState(false);
const [isOpen, setIsOpen] = useState(false);
```

### Callback Optimization

Use `useCallback` for event handlers passed to child components:

```tsx
const handleClick = useCallback(() => {
  // Handler logic
}, [dependency1, dependency2]);
```

### Memoization

Use `useMemo` for expensive computations:

```tsx
const baseClass = useMemo(() => (
  // Complex className composition
), [dependencies]);
```

### Reference Pattern

Use `useRef` for DOM element references:

```tsx
const inputRef = useRef<HTMLInputElement>(null);

const openFileDialog = useCallback(() => {
  inputRef.current?.click();
}, []);

<input ref={inputRef} type="file" />
```

### ID Generation

Use `useId` for accessible form element IDs:

```tsx
import { useId } from "react";

function MyComponent() {
  const helperId = useId();

  return (
    <>
      <input aria-describedby={helperId} />
      <span id={helperId}>Helper text</span>
    </>
  );
}
```

---

## End of Document

This comprehensive guide should be used as the primary reference when converting Figma designs to code using the Model Context Protocol. Always follow these patterns and conventions to maintain consistency across the codebase.
