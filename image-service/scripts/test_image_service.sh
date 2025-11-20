#!/bin/bash

# Test internal.ImageService.Compress() directly
# This script creates a Go program that imports and tests the actual ImageService

set -e

IMAGE_URL="$1"
TEST_DIR="/tmp/test_images"
OUTPUT_DIR="/tmp"

echo "🔧 Setting up Image Service compression test..."

# Create test directories
mkdir -p "$TEST_DIR"
mkdir -p "$OUTPUT_DIR"

# Download image if URL provided
if [ -n "$IMAGE_URL" ]; then
    echo "📥 Downloading image from: $IMAGE_URL"
    filename=$(basename "$IMAGE_URL")
    # Ensure it has an image extension
    if [[ ! "$filename" =~ \.(jpg|jpeg|png|gif|webp|bmp)$ ]]; then
        # Try to detect from Content-Type or default to jpg
        filename="test.jpg"
    fi

    curl -L -o "$TEST_DIR/$filename" "$IMAGE_URL"
    echo "✅ Downloaded: $TEST_DIR/$filename"
fi

# Find all image files in test directory
IMAGE_FILES=($(find "$TEST_DIR" -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.gif" -o -iname "*.webp" -o -iname "*.bmp" \)))

if [ ${#IMAGE_FILES[@]} -eq 0 ]; then
    echo "❌ No image files found in $TEST_DIR"
    echo "💡 Provide an image URL or place image files in $TEST_DIR"
    exit 1
fi

echo "📄 Found ${#IMAGE_FILES[@]} image file(s) to test"

# Create test Go program that uses the actual ImageService
cat > "$OUTPUT_DIR/test_image_service_main.go" << 'EOF'
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	// Import the internal package from the image-service
	"github.com/instrlabs/image-service/internal"
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
		log.Fatal("Usage: go run test_image_service_main.go <image_file>")
	}

	imagePath := os.Args[1]

	// Read image file
	originalData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Failed to read image file: %v", err)
	}

	// Create image service instance
	service := internal.NewImageService()

	// Compress image using the actual service
	compressedData, err := service.Compress(originalData)
	if err != nil {
		log.Fatalf("Image compression failed: %v", err)
	}

	// Calculate metrics
	originalSize := len(originalData)
	compressedSize := len(compressedData)
	reduction := float64(originalSize-compressedSize) / float64(originalSize) * 100

	// Display results
	filename := filepath.Base(imagePath)
	fmt.Printf("\n📊 Image Service Compression Results: %s\n", filename)
	fmt.Printf("   Original Size:  %s\n", formatBytes(originalSize))
	fmt.Printf("   Compressed Size: %s\n", formatBytes(compressedSize))
	fmt.Printf("   Reduction:      %.1f%%\n", reduction)

	if compressedSize < originalSize {
		fmt.Printf("   Status:         ✅ Compressed successfully\n")

		// Save compressed file
		ext := filepath.Ext(filename)
		nameWithoutExt := filename[:len(filename)-len(ext)]
		outputPath := filepath.Join(filepath.Dir(imagePath), "compressed_"+nameWithoutExt+ext)
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

echo "🧪 Running Image Service compression tests..."
echo

# Test each image file
total_original=0
total_compressed=0

for image_file in "${IMAGE_FILES[@]}"; do
    echo "Testing: $(basename "$image_file")"

    # Run the compression test using the actual Image service
    cd "$OUTPUT_DIR" && go run test_image_service_main.go "$image_file"

    # Get file sizes for summary
    original_size=$(stat -f%z "$image_file" 2>/dev/null || stat -c%s "$image_file" 2>/dev/null || echo "0")

    ext=$(echo "$image_file" | sed 's/.*\.//')
    nameWithoutExt=$(basename "$image_file" .$ext)
    compressed_file="$TEST_DIR/compressed_${nameWithoutExt}.$ext"

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

echo "✅ Image Service compression test completed!"

# Cleanup
rm -f "$OUTPUT_DIR/test_image_service_main.go"
