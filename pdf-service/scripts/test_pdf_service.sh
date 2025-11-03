#!/bin/bash

# Test internal.PDFService.Compress() directly
# This script creates a Go program that imports and tests the actual PDFService

set -e

PDF_URL="$1"
TEST_DIR="/tmp/test_pdfs"
OUTPUT_DIR="/tmp"

echo "🔧 Setting up PDF Service compression test..."

# Create test directories
mkdir -p "$TEST_DIR"
mkdir -p "$OUTPUT_DIR"

# Download PDF if URL provided
if [ -n "$PDF_URL" ]; then
    echo "📥 Downloading PDF from: $PDF_URL"
    filename=$(basename "$PDF_URL")
    if [[ "$filename" != *.pdf ]]; then
        filename="test.pdf"
    fi

    curl -L -o "$TEST_DIR/$filename" "$PDF_URL"
    echo "✅ Downloaded: $TEST_DIR/$filename"
fi

# Find all PDF files in test directory
PDF_FILES=($(find "$TEST_DIR" -name "*.pdf" -type f))

if [ ${#PDF_FILES[@]} -eq 0 ]; then
    echo "❌ No PDF files found in $TEST_DIR"
    echo "💡 Provide a PDF URL or place PDF files in $TEST_DIR"
    exit 1
fi

echo "📄 Found ${#PDF_FILES[@]} PDF file(s) to test"

# Create test Go program that uses the actual PDFService
cat > "$OUTPUT_DIR/test_pdf_service_main.go" << 'EOF'
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	// Import the internal package from the pdf-service
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
		log.Fatal("Usage: go run test_pdf_service_main.go <pdf_file>")
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
EOF

echo "🧪 Running PDF Service compression tests..."
echo

# Test each PDF file
total_original=0
total_compressed=0

for pdf_file in "${PDF_FILES[@]}"; do
    echo "Testing: $(basename "$pdf_file")"

    # Run the compression test using the actual PDF service
    cd "$OUTPUT_DIR" && go run test_pdf_service_main.go "$pdf_file"

    # Get file sizes for summary
    original_size=$(stat -f%z "$pdf_file" 2>/dev/null || stat -c%s "$pdf_file" 2>/dev/null || echo "0")
    compressed_file="$TEST_DIR/compressed_$(basename "$pdf_file")"

    if [ -f "$compressed_file" ]; then
        compressed_size=$(stat -f%z "$compressed_file" 2>/dev/null || stat -c%s "$compressed_file" 2>/dev/null || echo "0")
        total_original=$((total_original + original_size))
        total_compressed=$((total_compressed + compressed_size))
    else
        total_original=$((total_original + original_size))
        total_compressed=$((total_compressed + original_size))
    fi
done

# Summary
if [ $total_original -gt 0 ]; then
    if command -v bc >/dev/null 2>&1; then
        total_reduction=$(echo "scale=1; ($total_original - $total_compressed) * 100 / $total_original" | bc 2>/dev/null || echo "0")
    else
        total_reduction="0"
    fi
    echo "📈 Summary:"
    echo "   Total Original:    $(formatBytes $total_original)"
    echo "   Total Compressed:  $(formatBytes $total_compressed)"
    echo "   Overall Reduction: ${total_reduction}%"
fi

echo "✅ PDF Service compression test completed!"

# Cleanup
rm -f "$OUTPUT_DIR/test_pdf_service_main.go"