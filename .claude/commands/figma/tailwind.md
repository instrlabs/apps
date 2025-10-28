---
description: Convert extracted Figma design information into pure Tailwind CSS utility classes
---

# Figma to Tailwind CSS Conversion

Convert comprehensive Figma design information into pure Tailwind CSS utility classes ready to use in your projects.

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

3. **Convert Tokens to Tailwind Classes**:
   - Transform all parsed design tokens into Tailwind utility classes
   - Generate atomic utility classes from colors, spacing, typography values
   - Create composite classes combining multiple properties
   - Handle responsive variants and states (hover, focus, disabled, etc.)
   - Use standard Tailwind syntax with direct values

4. **Generate Component Classes**:
   - For each component variant found, generate:
     - Base component classes
     - Size variant classes (sm, base, lg, etc.)
     - Color variant classes (primary, secondary, etc.)
     - State variant classes (default, hover, disabled, etc.)
     - CSS custom properties (variables) mapping
   - Format as Tailwind utility classes with proper syntax

5. **Output Format**: Generate a comprehensive report including:
   - **Tailwind Classes CSS** (generated utility and component classes)
   - **Component Classes** (organized by component and variant)
   - **Utility Classes** (color, spacing, typography utilities)
   - **Usage Examples** (how to use generated classes in React/HTML)
   - **Token Mapping** (reference of tokens to classes)

6. **Save Output**:
   - Create `/tmp/figma_tailwind_<timestamp>.css` for all generated Tailwind classes
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

Convert Figma color values to Tailwind utilities:
- `#ffffff` → `bg-white`, `text-white`, `border-white`
- `#000000` → `bg-black`, `text-black`, `border-black`
- `#3b82f6` → `bg-blue-500`, `text-blue-500`, `border-blue-500`
- `rgba(255,255,255,0.96)` → `bg-white/96`, `text-white/96`, etc.
- Named variables like `colors/button/primary` → Map to nearest standard Tailwind color

## Example Output Structure

### Pure Tailwind Classes (figma_tailwind_*.md)

Output classes as Tailwind utility combinations ready to use in HTML/React:

**Color utilities:**
```
text-white, text-black, bg-white, bg-blue-500, etc.
```

**Spacing utilities:**
```
p-2, px-4, py-3, gap-2, gap-3, etc.
```

**Typography utilities:**
```
text-sm, text-base, text-lg, font-medium, font-semibold, leading-5, leading-6, etc.
```

**Component example:**
```
Button size variants:
- Small: h-9 px-3 py-2 text-sm gap-2 rounded
- Base: h-10 px-4 py-2 text-base gap-2 rounded
- Large: h-12 px-4 py-3 text-lg gap-3 rounded

Button color variants:
- Primary: bg-blue-500 text-white hover:bg-blue-600
- Secondary: bg-gray-100 text-gray-900 hover:bg-gray-200
- Danger: bg-red-500 text-white hover:bg-red-600
```

**Component usage in React:**
```jsx
<button className="h-10 px-4 py-2 text-base gap-2 rounded bg-blue-500 text-white hover:bg-blue-600">
  Click me
</button>

<button className="h-12 px-4 py-3 text-lg gap-3 rounded bg-gray-100 text-gray-900">
  Secondary
</button>
```

## Features

- **Token to Tailwind conversion**: Converts Figma design tokens to standard Tailwind utilities
- **Component mapping**: Identifies component variants and maps to Tailwind class combinations
- **Color mapping**: Converts hex/rgba colors to nearest Tailwind colors or suggests custom values
- **Responsive variants**: Generates breakpoint-aware class combinations if applicable
- **State variants**: Includes hover, focus, disabled states for all components
- **Direct usability**: Generates classes ready to use in HTML/React without additional configuration

## Output Includes

- ✅ Pure Tailwind utility classes (no custom CSS)
- ✅ Component-specific class combinations
- ✅ Color palette mapping to Tailwind colors
- ✅ Spacing and sizing utilities
- ✅ Typography utilities
- ✅ Usage examples with React/HTML code
- ✅ Component variant reference

---

**Note**: This command requires a valid Figma design extraction file from `/figma:extract` command.
