package internal

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/disintegration/imaging"
)

type ImageService struct{}

func NewImageService() *ImageService { return &ImageService{} }

func (s *ImageService) Run(productKey string, file []byte) ([]byte, error) {
	switch productKey {
	case "images/compress":
		return s.Compress(file)
	default:
		return nil, fmt.Errorf("unsupported product key: %s", productKey)
	}
}

func (s *ImageService) Compress(file []byte) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(file))
	if err != nil {
		log.Printf("Failed to decode image: %v", err)
		return nil, err
	}

	var buf bytes.Buffer
	format := detectFormat(file)
	log.Printf("Compressing image: format=%s size=%d", format, len(file))

	switch format {
	case "png":
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		if err := enc.Encode(&buf, img); err != nil {
			log.Printf("Failed to encode PNG: %v", err)
			return nil, err
		}
	case "jpeg", "jpg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 60}); err != nil {
			log.Printf("Failed to encode JPEG: %v", err)
			return nil, err
		}
	default:
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 60}); err != nil {
			log.Printf("Failed to encode image (default JPEG): %v", err)
			return nil, err
		}
	}

	out := buf.Bytes()
	log.Printf("Image compressed successfully: originalSize=%d compressedSize=%d ratio=%.2f",
		len(file), len(out), float64(len(out))/float64(len(file)))

	return out, nil
}

func detectFormat(b []byte) string {
	if len(b) >= 8 && bytes.Equal(b[:8], []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}) {
		return "png"
	}
	if len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF {
		return "jpeg"
	}
	if len(b) >= 6 && bytes.Equal(b[:6], []byte{'G', 'I', 'F', '8', '9', 'a'}) ||
		len(b) >= 6 && bytes.Equal(b[:6], []byte{'G', 'I', 'F', '8', '7', 'a'}) {
		return "gif"
	}
	if _, format, err := image.DecodeConfig(bytes.NewReader(b)); err == nil {
		return format
	}
	return ""
}
