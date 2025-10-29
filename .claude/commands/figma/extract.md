---
description: Extract Figma design properties, layouts, colors, sizing, components, and variables from a Figma link
---

# Figma Design Extraction

Extract comprehensive design information from the provided Figma link including:
- Layout properties (spacing, positioning, constraints)
- Color palette and styles
- Sizing and dimensions
- Components and variants
- Variables (colors, numbers, strings, booleans)

## Task

1. Parse the Figma link to extract the `fileKey` and `nodeId`
2. Use the appropriate Figma MCP tools to gather design context:
   - Use `mcp__figma__get_design_context` to get the main design context and code
   - Use `mcp__figma__get_variable_defs` to extract variable definitions
   - Use `mcp__figma__get_metadata` to get structural information if needed
   - Use `mcp__figma__get_screenshot` to capture visual representation
   - Use `mcp__figma__get_code_connect_map` to find code connections
3. Organize and present the extracted information in the following categories:
   - **Layout Properties**: Padding, margins, auto-layout settings, constraints
   - **Colors**: Color palette, fill styles, stroke colors, variable colors
   - **Sizing**: Width, height, min/max dimensions, sizing modes
   - **Components**: Component names, variants, properties, instances
   - **Variables**: Organized by type (color, number, string, boolean) with their values
   - **Code Connections**: Mapped component names and file locations
   - **Metadata**: Node structure, layer types, positions
4. Save ALL extracted information to a `.txt` file in `/tmp/` directory
   - Filename format: `/tmp/figma_extract_<fileKey>_<nodeId>_<timestamp>.txt`
   - Include ALL raw data from Figma API responses
   - Include formatted summary sections
5. Display a summary to the user and provide the file path

## Parameters

Expect the Figma URL in the format:
- `https://figma.com/design/:fileKey/:fileName?node-id=:nodeId`
- Example: `https://figma.com/design/abc123/MyDesign?node-id=1-2`

## Output Format

1. **Save to File**: Write ALL extracted data to `/tmp/figma_extract_<fileKey>_<nodeId>_<timestamp>.txt`
   - Include complete raw JSON responses from all Figma API calls
   - Include formatted sections for easy reading
   - Include metadata about the extraction (timestamp, URL, file/node IDs)

2. **Console Summary**: Display a brief summary including:
   - Number of components found
   - Number of variables extracted
   - Color palette summary
   - File path where complete data was saved

## Example Usage

```
/figma:extract https://figma.com/design/abc123/MyDesign?node-id=1-2
```

---

**Note**: This command requires the Figma MCP server to be configured and authenticated.
