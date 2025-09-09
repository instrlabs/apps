package internal

import (
	"bytes"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

// ImageService provides image-related operations.
// It wraps lower-level functions to facilitate dependency injection and testing.
type ImageService struct{}

// NewImageService constructs a new ImageService instance.
func NewImageService() *ImageService { return &ImageService{} }

// Compress compresses a JPEG/PNG image to a reasonable quality.
func (s *ImageService) Compress(file []byte) ([]byte, error) { return ImageCompress(file) }

// ImageCompress is a helper function that performs JPEG encoding at a set quality.
func ImageCompress(file []byte) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 70}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
