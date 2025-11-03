---
description: Test PDF compression function by calling internal.PDFService.Compress() directly
---

# PDF Compression Test

Test `internal.PDFService.Compress()` directly without HTTP/NATS. Accepts PDF URL or uses existing files.

## Usage

```bash
# Test with a PDF URL
/pdf:compress https://example.com/document.pdf

# Test with existing PDF files in /tmp/test_pdfs/
/pdf:compress
```

## Workflow

1. **Prepare PDFs**: If URL provided, download to `/tmp/test_pdfs/`. Otherwise use existing PDFs in that directory.

2. **Run test**: Execute `/tmp/test_pdf_service.sh` which:
   - Creates `/tmp/test_pdf_service_main.go` (Go program importing internal.PDFService)
   - Compiles and runs test on all PDFs in `/tmp/test_pdfs/`
   - Displays: original size, compressed size, reduction %, status
   - Saves compressed outputs to `/tmp/test_pdfs/compressed_<filename>`

3. **Results**: Show compression metrics for each PDF tested and overall summary.

## Implementation Details

The compression uses `pdfcpu` library to optimize PDF files by:
- Removing redundant objects
- Compressing image streams
- Optimizing font data
- Removing metadata and comments

Files are processed in memory and results are saved with `compressed_` prefix.