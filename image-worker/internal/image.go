package internal

import (
	"bytes"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

type ImageService struct{}

func NewImageService() *ImageService { return &ImageService{} }

func (s *ImageService) Compress(file []byte) ([]byte, error) {
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
