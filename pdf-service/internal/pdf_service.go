package internal

import (
	"bytes"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

// Compress processes a PDF file (simplified version - just validates for now)
func (s *PDFService) Compress(file []byte) ([]byte, error) {
	// For now, just validate the file and return it
	// TODO: Implement actual compression once API is properly understood
	err := s.Validate(file)
	if err != nil {
		return nil, fmt.Errorf("invalid PDF file: %w", err)
	}

	// Return the file as-is for now
	// This will be replaced with actual compression later
	return file, nil
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
