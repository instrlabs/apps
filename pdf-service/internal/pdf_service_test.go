package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPDFService_Compress tests the PDF compression functionality
func TestPDFService_Compress(t *testing.T) {
	service := NewPDFService()

	// Test with minimal valid PDF content
	minimalPDF := `%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj

2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj

3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
/Contents 4 0 R
>>
endobj

4 0 obj
<<
/Length 44
>>
stream
BT
/F1 12 Tf
100 700 Td
(Hello World) Tj
ET
endstream
endobj

xref
0 5
0000000000 65535 f
0000000009 00000 n
0000000058 00000 n
0000000115 00000 n
0000000204 00000 n
trailer
<<
/Size 5
/Root 1 0 R
>>
startxref
299
%%EOF`

	pdfBytes := []byte(minimalPDF)

	t.Run("Valid PDF compression", func(t *testing.T) {
		// First validate the PDF
		err := service.Validate(pdfBytes)
		assert.NoError(t, err)

		// Compress the PDF
		compressed, err := service.Compress(pdfBytes)
		assert.NoError(t, err)
		assert.NotEmpty(t, compressed)

		// The compressed result should still be a valid PDF
		err = service.Validate(compressed)
		assert.NoError(t, err)
	})

	t.Run("Invalid PDF", func(t *testing.T) {
		invalidPDF := []byte("not a pdf")
		err := service.Validate(invalidPDF)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid PDF file")
	})
}

// TestPDFService_Validate tests PDF validation
func TestPDFService_Validate(t *testing.T) {
	service := NewPDFService()

	t.Run("Valid PDF header", func(t *testing.T) {
		validPDF := `%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj

2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj

3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
/Contents 4 0 R
>>
endobj

4 0 obj
<<
/Length 44
>>
stream
BT
/F1 12 Tf
100 700 Td
(Hello World) Tj
ET
endstream
endobj

xref
0 5
0000000000 65535 f
0000000009 00000 n
0000000058 00000 n
0000000115 00000 n
0000000204 00000 n
trailer
<<
/Size 5
/Root 1 0 R
>>
startxref
299
%%EOF`

		err := service.Validate([]byte(validPDF))
		assert.NoError(t, err)
	})

	t.Run("Invalid PDF header", func(t *testing.T) {
		invalidPDF := "This is not a PDF file"
		err := service.Validate([]byte(invalidPDF))
		assert.Error(t, err)
	})

	t.Run("Empty PDF", func(t *testing.T) {
		emptyPDF := []byte("")
		err := service.Validate(emptyPDF)
		assert.Error(t, err)
	})
}
