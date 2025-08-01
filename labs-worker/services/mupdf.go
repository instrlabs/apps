package services

// #cgo CFLAGS: -I/go/src/labs-worker/libraries/mupdf/include
// #cgo LDFLAGS: -L/go/src/labs-worker/libraries/mupdf/build/release -lmupdf -lmupdf-third -lm
// #cgo LDFLAGS: -lfreetype -ljbig2dec -ljpeg -lz -lopenjp2
// #include <stdlib.h>
// #include <string.h>
// #include "mupdf/fitz.h"
// #include "mupdf/pdf.h"
import "C"
import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type MuPDFService struct {
}

func NewMuPDFService() *MuPDFService {
	return &MuPDFService{}
}

func (m *MuPDFService) Close() {
	// MuPDF context is created and destroyed for each operation,
	// so no global cleanup is needed
}

// ConvertPDFToJPG converts a PDF file to JPG images
// It returns the path to the output JPG file or an error
func (m *MuPDFService) ConvertPDFToJPG(pdfPath string) (string, error) {
	// Convert Go string to C string
	cPdfPath := C.CString(pdfPath)
	defer C.free(unsafe.Pointer(cPdfPath))

	// Create a context
	ctx := C.fz_new_context(nil, nil, C.FZ_STORE_UNLIMITED)
	if ctx == nil {
		return "", errors.New("failed to create context")
	}
	defer C.fz_drop_context(ctx)

	// Register document handlers
	C.fz_register_document_handlers(ctx)

	// Create an error context
	var err C.fz_error_context
	C.fz_try(ctx, &err)

	// Open the document
	doc := C.fz_open_document(ctx, cPdfPath)
	if doc == nil {
		return "", errors.New("failed to open document")
	}
	defer C.fz_drop_document(ctx, doc)

	// Get the number of pages
	pageCount := C.fz_count_pages(ctx, doc)
	if pageCount == 0 {
		return "", errors.New("document has no pages")
	}

	// Create output directory
	outputDir := filepath.Join(os.TempDir(), filepath.Base(pdfPath)+"_jpg")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Process the first page only for now
	// In a real implementation, you might want to process all pages
	pageNum := C.int(0)
	page := C.fz_load_page(ctx, doc, pageNum)
	if page == nil {
		return "", errors.New("failed to load page")
	}
	defer C.fz_drop_page(ctx, page)

	// Get the bounds of the page
	bounds := C.fz_bound_page(ctx, page)

	// Create a pixmap with the bounds
	pixmap := C.fz_new_pixmap_with_bbox(ctx, C.fz_device_rgb(ctx), &bounds, nil, 1)
	if pixmap == nil {
		return "", errors.New("failed to create pixmap")
	}
	defer C.fz_drop_pixmap(ctx, pixmap)

	// Clear the pixmap with white
	C.fz_clear_pixmap_with_value(ctx, pixmap, 0xff)

	// Create a device for drawing
	dev := C.fz_new_draw_device(ctx, C.fz_identity, pixmap)
	if dev == nil {
		return "", errors.New("failed to create draw device")
	}
	defer C.fz_drop_device(ctx, dev)

	// Run the page
	C.fz_run_page(ctx, page, dev, C.fz_identity, nil)

	// Output path for the JPG
	outputPath := filepath.Join(outputDir, "page_1.jpg")
	cOutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cOutputPath))

	// Save the pixmap as a JPG
	C.fz_save_pixmap_as_jpeg(ctx, pixmap, cOutputPath, 90)

	return outputPath, nil
}

// Compress compresses a PDF file
// It returns the path to the compressed PDF file or an error
func (m *MuPDFService) Compress(pdfPath string) (string, error) {
	// Convert Go string to C string
	cPdfPath := C.CString(pdfPath)
	defer C.free(unsafe.Pointer(cPdfPath))

	// Create a context
	ctx := C.fz_new_context(nil, nil, C.FZ_STORE_UNLIMITED)
	if ctx == nil {
		return "", errors.New("failed to create context")
	}
	defer C.fz_drop_context(ctx)

	// Register document handlers
	C.fz_register_document_handlers(ctx)

	// Create an error context
	var err C.fz_error_context
	C.fz_try(ctx, &err)

	// Open the document
	doc := C.pdf_open_document(ctx, cPdfPath)
	if doc == nil {
		return "", errors.New("failed to open document")
	}
	defer C.pdf_drop_document(ctx, doc)

	// Output path for the compressed PDF
	outputPath := filepath.Join(os.TempDir(), "compressed_"+filepath.Base(pdfPath))
	cOutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cOutputPath))

	// Create options for saving
	opts := C.pdf_write_options{}
	opts.do_compress = 1
	opts.do_compress_images = 1
	opts.do_garbage = 1

	// Save the document with compression
	C.pdf_save_document(ctx, doc, cOutputPath, &opts)

	return outputPath, nil
}

// Merge merges multiple PDF files
// It returns the path to the merged PDF file or an error
func (m *MuPDFService) Merge(pdfPaths []string) (string, error) {
	if len(pdfPaths) == 0 {
		return "", errors.New("no PDF files to merge")
	}

	// Create a context
	ctx := C.fz_new_context(nil, nil, C.FZ_STORE_UNLIMITED)
	if ctx == nil {
		return "", errors.New("failed to create context")
	}
	defer C.fz_drop_context(ctx)

	// Register document handlers
	C.fz_register_document_handlers(ctx)

	// Create an error context
	var err C.fz_error_context
	C.fz_try(ctx, &err)

	// Create a new PDF document
	doc := C.pdf_create_document(ctx)
	if doc == nil {
		return "", errors.New("failed to create document")
	}
	defer C.pdf_drop_document(ctx, doc)

	// Process each input PDF
	for _, pdfPath := range pdfPaths {
		// Convert Go string to C string
		cPdfPath := C.CString(pdfPath)
		defer C.free(unsafe.Pointer(cPdfPath))

		// Open the source document
		srcDoc := C.pdf_open_document(ctx, cPdfPath)
		if srcDoc == nil {
			return "", fmt.Errorf("failed to open document: %s", pdfPath)
		}
		defer C.pdf_drop_document(ctx, srcDoc)

		// Get the number of pages
		pageCount := C.pdf_count_pages(ctx, srcDoc)

		// Add each page to the destination document
		for i := C.int(0); i < pageCount; i++ {
			// Load the page
			page := C.pdf_load_page(ctx, srcDoc, i)
			if page == nil {
				return "", fmt.Errorf("failed to load page %d from %s", i, pdfPath)
			}
			defer C.pdf_drop_page(ctx, page)

			// Add the page to the destination document
			C.pdf_insert_page(ctx, doc, -1, C.fz_page(page))
		}
	}

	// Output path for the merged PDF
	outputPath := filepath.Join(os.TempDir(), "merged.pdf")
	cOutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cOutputPath))

	// Save the document
	opts := C.pdf_write_options{}
	C.pdf_save_document(ctx, doc, cOutputPath, &opts)

	return outputPath, nil
}

// Split splits a PDF file into multiple PDF files
// It returns the directory containing the split PDF files or an error
func (m *MuPDFService) Split(pdfPath string) (string, error) {
	// Convert Go string to C string
	cPdfPath := C.CString(pdfPath)
	defer C.free(unsafe.Pointer(cPdfPath))

	// Create a context
	ctx := C.fz_new_context(nil, nil, C.FZ_STORE_UNLIMITED)
	if ctx == nil {
		return "", errors.New("failed to create context")
	}
	defer C.fz_drop_context(ctx)

	// Register document handlers
	C.fz_register_document_handlers(ctx)

	// Create an error context
	var err C.fz_error_context
	C.fz_try(ctx, &err)

	// Open the document
	doc := C.pdf_open_document(ctx, cPdfPath)
	if doc == nil {
		return "", errors.New("failed to open document")
	}
	defer C.pdf_drop_document(ctx, doc)

	// Get the number of pages
	pageCount := C.pdf_count_pages(ctx, doc)
	if pageCount == 0 {
		return "", errors.New("document has no pages")
	}

	// Create output directory
	outputDir := filepath.Join(os.TempDir(), filepath.Base(pdfPath)+"_split")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Process each page
	for i := C.int(0); i < pageCount; i++ {
		// Create a new PDF document for this page
		newDoc := C.pdf_create_document(ctx)
		if newDoc == nil {
			return "", fmt.Errorf("failed to create document for page %d", i)
		}
		defer C.pdf_drop_document(ctx, newDoc)

		// Load the page
		page := C.pdf_load_page(ctx, doc, i)
		if page == nil {
			return "", fmt.Errorf("failed to load page %d", i)
		}
		defer C.pdf_drop_page(ctx, page)

		// Add the page to the new document
		C.pdf_insert_page(ctx, newDoc, -1, C.fz_page(page))

		// Output path for this page
		outputPath := filepath.Join(outputDir, fmt.Sprintf("page_%d.pdf", i+1))
		cOutputPath := C.CString(outputPath)
		defer C.free(unsafe.Pointer(cOutputPath))

		// Save the document
		opts := C.pdf_write_options{}
		C.pdf_save_document(ctx, newDoc, cOutputPath, &opts)
	}

	return outputDir, nil
}
