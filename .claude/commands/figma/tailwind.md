---
description: Convert extracted Figma design information into Tailwind CSS classes and configuration
---

# Figma to Tailwind CSS Conversion

Convert comprehensive Figma design information into Tailwind CSS classes, configuration files, and utility classes.

## Task

1. **Input**: Accept a path to a Figma design extraction file (`/tmp/figma_extract_*.txt`)
   - The file should contain the output from the `/figma:extract` command
   - Extract design tokens from the file (colors, spacing, typography, border-radius, etc.)

2. **Parse Design Tokens**:
   - Extract all color variables and map to Tailwind color palette
   - Extract spacing/sizing values and map to Tailwind spacing scale
   - Extract typography (font-family, font-size, font-weight, line-height)
   - Extract border-radius values
   - Extract shadows if present
   - Identify all component variants and their properties

3. **Generate Tailwind Configuration**:
   - Create a `tailwind.config.ts` snippet with:
     - Extended color palette from Figma variables
     - Custom spacing values
     - Typography/font configuration
     - Border-radius customizations
     - Any custom theme values

4. **Generate Component Classes**:
   - For each component variant found, generate:
     - Base component classes
     - Size variant classes (sm, base, lg, etc.)
     - Color variant classes (primary, secondary, etc.)
     - State variant classes (default, hover, disabled, etc.)
     - CSS custom properties (variables) mapping
   - Format as Tailwind utility classes with proper syntax

5. **Output Format**: Generate a comprehensive report including:
   - **Tailwind Config File** (`tailwind.config.ts`)
   - **CSS Variables File** (CSS custom properties for dynamic theming)
   - **Component Classes** (organized by component and variant)
   - **Color Palette** (Tailwind color scale mapping)
   - **Spacing Scale** (Tailwind spacing mapping)
   - **Typography System** (font classes and utility classes)
   - **Usage Examples** (how to use generated classes in React/HTML)

6. **Save Output**:
   - Create `/tmp/figma_tailwind_<timestamp>.ts` for Tailwind config
   - Create `/tmp/figma_tailwind_<timestamp>.css` for CSS variables
   - Create `/tmp/figma_tailwind_<timestamp>.md` for complete documentation
   - Display file paths and summary to user

## Parameters

Accept the file path as input:
```
/figma:tailwind /tmp/figma_extract_*.txt
```

Or accept a Figma URL (will look for most recent extraction):
```
/figma:tailwind https://figma.com/design/abc123/MyDesign?node-id=1-2
```

## Color Mapping Strategy

Convert Figma color values to Tailwind:
- `#ffffff` → `white` or `colors/white`
- `#000000` → `black` or `colors/black`
- `rgba(255,255,255,0.96)` → CSS custom property with fallback
- Named variables like `colors/button/primary` → Tailwind palette key

## Example Output Structure

### tailwind.config.ts
```typescript
export default {
  theme: {
    extend: {
      colors: {
        'button-primary': '#fffffff5',
        'button-primary-hover': '#ffffff',
        'button-secondary': '#ffffff0a',
        'text-black': '#000000',
        // ... more colors
      },
      spacing: {
        'compact': '8px',   // spacing/2
        'standard': '12px', // spacing/3
        // ... more spacing
      },
      borderRadius: {
        'DEFAULT': '4px',
        // ... more radius
      },
      fontFamily: {
        'geist': ['Geist', 'sans-serif'],
      },
      fontSize: {
        'sm': ['14px', { lineHeight: '20px' }],
        'base': ['16px', { lineHeight: '24px' }],
        // ... more sizes
      },
      fontWeight: {
        'medium': 500,
        'semibold': 600,
      },
    },
  },
}
```

### Component Classes
```css
/* Button Component */
.btn-sm { @apply h-9 px-2 py-2 rounded gap-2; }
.btn-base { @apply h-10 px-2 py-2 rounded gap-2; }
.btn-lg { @apply h-12 px-3 py-3 rounded gap-3; }

.btn-primary { @apply bg-button-primary text-black; }
.btn-primary:hover { @apply bg-button-primary-hover; }
.btn-primary:disabled { @apply opacity-60 cursor-not-allowed; }

.btn-secondary { @apply bg-button-secondary text-white border border-border-primary; }
.btn-secondary:hover { @apply bg-button-secondary-hover; }
```

## Features

- **Automatic scaling**: Converts all numeric values to Tailwind scale
- **Variable mapping**: Links Figma variables to Tailwind config
- **Component variants**: Generates classes for all variants found
- **CSS Variables**: Creates CSS custom properties for theme switching
- **Responsive**: Includes breakpoint-aware variants if applicable
- **Accessibility**: Includes focus states and disabled states
- **Type safety**: Generates TypeScript-compatible config files

## Output Includes

- ✅ Complete Tailwind configuration
- ✅ CSS custom properties (variables)
- ✅ Component-specific utility classes
- ✅ Color palette mapping
- ✅ Spacing system
- ✅ Typography system
- ✅ Usage examples
- ✅ Implementation guide

---

**Note**: This command requires a valid Figma design extraction file from `/figma:extract` command.
