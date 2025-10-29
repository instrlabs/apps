---
description: Test image compression function by calling internal.ImageService.Compress() directly
---

# Image Compression Test

Test `internal.ImageService.Compress()` directly without HTTP/NATS. Accepts image URL or uses existing files.

## Workflow

1. **Prepare images**: If URL provided, download to `/tmp/test_images/`. Otherwise use existing images in that directory.

2. **Run test**: Execute `/tmp/test_image_compress.sh` which:
   - Creates `/tmp/test_compress_main.go` (Go program calling ImageService.Compress)
   - Compiles and runs test on all images in `/tmp/test_images/`
   - Displays: original size, compressed size, reduction %, format, status
   - Saves compressed outputs to `/tmp/compressed_<filename>`

3. **Results**: Show compression metrics for each image tested.