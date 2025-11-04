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

func (s *PDFService) Compress(file []byte) ([]byte, error) {
	err := s.Validate(file)
	if err != nil {
		return nil, fmt.Errorf("invalid PDF file: %w", err)
	}

	tempDir := "/tmp"
	inputFile := filepath.Join(tempDir, "input.pdf")
	outputFile := filepath.Join(tempDir, "output.pdf")

	err = os.WriteFile(inputFile, file, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	err = api.OptimizeFile(inputFile, outputFile, nil)
	if err != nil {
		log.Errorf("Failed to compress PDF: %v", err)
		return nil, fmt.Errorf("failed to compress PDF: %w", err)
	}

	compressed, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed file: %w", err)
	}

	if len(compressed) >= len(file) {
		log.Infof("PDF compression did not reduce size, returning original")
		return file, nil
	}

	log.Infof("PDF compressed successfully: %d bytes -> %d bytes (%.1f%% reduction)",
		len(file), len(compressed), float64(len(file)-len(compressed))/float64(len(file))*100)

	return compressed, nil
}

func (s *PDFService) Validate(file []byte) error {
	reader := bytes.NewReader(file)
	_, err := api.ReadContext(reader, nil)
	if err != nil {
		return fmt.Errorf("invalid PDF file: %w", err)
	}

	return nil
}
