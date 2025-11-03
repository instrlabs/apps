package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

// Compress processes a PDF file using pdfcpu optimization
func (s *PDFService) Compress(file []byte) ([]byte, error) {
	// Validate the PDF first
	err := s.Validate(file)
	if err != nil {
		return nil, fmt.Errorf("invalid PDF file: %w", err)
	}

	// Create temporary files for pdfcpu processing
	tempDir := "/tmp"
	inputFile := filepath.Join(tempDir, "input.pdf")
	outputFile := filepath.Join(tempDir, "output.pdf")

	// Write input file
	err = os.WriteFile(inputFile, file, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	// Use pdfcpu to optimize/compress the PDF
	err = api.OptimizeFile(inputFile, outputFile, nil)
	if err != nil {
		log.Errorf("Failed to compress PDF: %v", err)
		return nil, fmt.Errorf("failed to compress PDF: %w", err)
	}

	// Read the compressed file
	compressed, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed file: %w", err)
	}

	// Check if compression actually reduced size
	if len(compressed) >= len(file) {
		log.Infof("PDF compression did not reduce size, returning original")
		return file, nil
	}

	log.Infof("PDF compressed successfully: %d bytes -> %d bytes (%.1f%% reduction)",
		len(file), len(compressed), float64(len(file)-len(compressed))/float64(len(file))*100)

	return compressed, nil
}

// Validate checks if the provided data is a valid PDF
func (s *PDFService) Validate(file []byte) error {
	reader := bytes.NewReader(file)

	// Try to read PDF context to validate
	_, err := api.ReadContext(reader, nil)
	if err != nil {
		return fmt.Errorf("invalid PDF file: %w", err)
	}

	return nil
}
