package services

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lmupdf -lmupdf-third
// #include <mupdf/fitz.h>
// #include <stdlib.h>
// #include <string.h>
//
// // Helper function to render a PDF page to a JPEG file
// int render_page_to_jpeg(const char *pdf_path, const char *jpg_path, int page_num, float zoom, int rotate) {
//     fz_context *ctx = NULL;
//     fz_document *doc = NULL;
//     fz_pixmap *pix = NULL;
//     fz_matrix ctm;
//     fz_page *page = NULL;
//     fz_rect bounds;
//     fz_device *dev = NULL;
//     FILE *f = NULL;
//     int ret = -1;
//
//     // Create a context
//     ctx = fz_new_context(NULL, NULL, FZ_STORE_UNLIMITED);
//     if (!ctx) {
//         return -1;
//     }
//
//     // Register document handlers
//     fz_try(ctx) {
//         fz_register_document_handlers(ctx);
//
//         // Open the PDF
//         doc = fz_open_document(ctx, pdf_path);
//         if (!doc) {
//             fz_throw(ctx, FZ_ERROR_GENERIC, "cannot open document");
//         }
//
//         // Load the page
//         page = fz_load_page(ctx, doc, page_num);
//
//         // Get the page bounds
//         bounds = fz_bound_page(ctx, page);
//
//         // Set up transformation matrix (scale and rotation)
//         ctm = fz_scale(zoom, zoom);
//         ctm = fz_pre_rotate(ctm, rotate);
//
//         // Create a pixmap to hold the rendered page
//         pix = fz_new_pixmap_from_page_contents(ctx, page, ctm, fz_device_rgb(ctx), 0);
//
//         // Save the pixmap as JPEG
//         f = fopen(jpg_path, "wb");
//         if (!f) {
//             fz_throw(ctx, FZ_ERROR_GENERIC, "cannot open output file");
//         }
//
//         fz_write_pixmap_as_jpeg(ctx, f, pix, 90);
//         fclose(f);
//         f = NULL;
//
//         ret = 0; // Success
//     }
//     fz_catch(ctx) {
//         ret = -1;
//     }
//
//     // Clean up
//     if (f) fclose(f);
//     if (pix) fz_drop_pixmap(ctx, pix);
//     if (page) fz_drop_page(ctx, page);
//     if (doc) fz_drop_document(ctx, doc);
//     if (ctx) fz_drop_context(ctx);
//
//     return ret;
// }
import "C"
import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

// MuPDFService handles PDF operations using MuPDF
type MuPDFService struct {
}

// NewMuPDFService creates a new MuPDFService
func NewMuPDFService() *MuPDFService {
	return &MuPDFService{}
}

// Close closes the MuPDFService
func (m *MuPDFService) Close() {
	// MuPDF context is created and destroyed for each operation,
	// so no global cleanup is needed
}

// ConvertPDFToJPG converts a PDF file to a JPG file using MuPDF
func (m *MuPDFService) ConvertPDFToJPG(pdfPath string) (string, error) {
	// Create a temporary file to store the JPG
	tempFile, err := os.CreateTemp("", "jpg-*.jpg")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Convert PDF to JPG using MuPDF
	// For simplicity, we're only converting the first page
	// In a real application, you might want to handle multi-page PDFs differently
	pdfPathC := C.CString(pdfPath)
	jpgPathC := C.CString(tempFile.Name())
	defer C.free(unsafe.Pointer(pdfPathC))
	defer C.free(unsafe.Pointer(jpgPathC))

	// Render the first page (page 0) with zoom=2.0 (200% scale) and no rotation
	result := C.render_page_to_jpeg(pdfPathC, jpgPathC, 0, 2.0, 0)
	if result != 0 {
		return "", fmt.Errorf("failed to convert PDF to JPG using MuPDF")
	}

	log.Printf("Converted PDF to JPG using MuPDF: %s -> %s", pdfPath, tempFile.Name())
	return tempFile.Name(), nil
}

// ConvertPDFToJPGMultiPage converts all pages of a PDF file to JPG files
func (m *MuPDFService) ConvertPDFToJPGMultiPage(pdfPath string) ([]string, error) {
	// This is a placeholder for a more advanced implementation
	// In a real application, you would iterate through all pages of the PDF
	// and convert each one to a separate JPG file

	// For now, we'll just convert the first page as an example
	jpgPath, err := m.ConvertPDFToJPG(pdfPath)
	if err != nil {
		return nil, err
	}

	return []string{jpgPath}, nil
}
