package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMuPDFService_ConvertPDFToJPG(t *testing.T) {
	// Create a temporary PDF file for testing
	tempDir, err := os.MkdirTemp("", "mupdf-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple PDF file for testing
	// Note: In a real test, you would use a real PDF file
	// For this test, we'll just check if the function handles errors correctly
	pdfPath := filepath.Join(tempDir, "test.pdf")
	err = os.WriteFile(pdfPath, []byte("This is not a valid PDF file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test PDF file: %v", err)
	}

	// Create the MuPDFService
	service := NewMuPDFService()
	defer service.Close()

	// Test with invalid PDF
	_, err = service.ConvertPDFToJPG(pdfPath)
	if err == nil {
		t.Error("Expected error when converting invalid PDF, but got nil")
	}

	// Test with non-existent file
	_, err = service.ConvertPDFToJPG("non-existent-file.pdf")
	if err == nil {
		t.Error("Expected error when converting non-existent file, but got nil")
	}

	// Note: To test with a valid PDF, you would need to include a real PDF file
	// and check that the output JPG file exists and is valid
	t.Log("To fully test this functionality, you need to use a real PDF file")
}

func TestMuPDFService_ConvertPDFToJPGMultiPage(t *testing.T) {
	// Create a temporary PDF file for testing
	tempDir, err := os.MkdirTemp("", "mupdf-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple PDF file for testing
	// Note: In a real test, you would use a real PDF file
	pdfPath := filepath.Join(tempDir, "test.pdf")
	err = os.WriteFile(pdfPath, []byte("This is not a valid PDF file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test PDF file: %v", err)
	}

	// Create the MuPDFService
	service := NewMuPDFService()
	defer service.Close()

	// Test with invalid PDF
	_, err = service.ConvertPDFToJPGMultiPage(pdfPath)
	if err == nil {
		t.Error("Expected error when converting invalid PDF, but got nil")
	}

	// Test with non-existent file
	_, err = service.ConvertPDFToJPGMultiPage("non-existent-file.pdf")
	if err == nil {
		t.Error("Expected error when converting non-existent file, but got nil")
	}

	// Note: To test with a valid PDF, you would need to include a real PDF file
	// and check that the output JPG files exist and are valid
	t.Log("To fully test this functionality, you need to use a real PDF file")
}
