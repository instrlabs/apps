package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/instrlabs/pdf-service/internal"
)

func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run test_compress.go <pdf_file>")
	}

	pdfPath := os.Args[1]

	// Read PDF file
	originalData, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		log.Fatalf("Failed to read PDF file: %v", err)
	}

	// Create PDF service instance
	service := internal.NewPDFService()

	// Compress PDF using the actual service
	compressedData, err := service.Compress(originalData)
	if err != nil {
		log.Fatalf("PDF compression failed: %v", err)
	}

	// Calculate metrics
	originalSize := len(originalData)
	compressedSize := len(compressedData)
	reduction := float64(originalSize-compressedSize) / float64(originalSize) * 100

	// Display results
	filename := filepath.Base(pdfPath)
	fmt.Printf("\n📊 PDF Service Compression Results: %s\n", filename)
	fmt.Printf("   Original Size:  %s\n", formatBytes(originalSize))
	fmt.Printf("   Compressed Size: %s\n", formatBytes(compressedSize))
	fmt.Printf("   Reduction:      %.1f%%\n", reduction)

	if compressedSize < originalSize {
		fmt.Printf("   Status:         ✅ Compressed successfully\n")

		// Save compressed file
		outputPath := filepath.Join(filepath.Dir(pdfPath), "compressed_"+filename)
		err = ioutil.WriteFile(outputPath, compressedData, 0644)
		if err != nil {
			log.Printf("Warning: Failed to save compressed file: %v", err)
		} else {
			fmt.Printf("   Output:         %s\n", outputPath)
		}
	} else {
		fmt.Printf("   Status:         ⚠️  No compression achieved\n")
	}

	fmt.Printf("\n")
}
