package internal

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CompressJPEG(instructionHex string) error {
	files, _, err := w.s3.DownloadAllForInstruction(instructionHex)
	if err != nil {
		return err
	}

	for idx, data := range files {
		compressed, err := compressJPEGFromReader(bytes.NewReader(data), 70)
		if err != nil {
			return err
		}

		objectName := fmt.Sprintf("images/%s-%d-compressed.jpg", instructionHex, idx)
		if err := w.s3.UploadBytes(compressed, objectName, "image/jpeg"); err != nil {
			return err
		}

		if w.fileRepo != nil {
			if id, e := primitive.ObjectIDFromHex(instructionHex); e == nil {
				_, _ = w.fileRepo.Create(&File{InstructionID: id, Type: FileTypeResponse})
			}
		}
	}
	return nil
}

func compressJPEGFromReader(r io.Reader, quality int) ([]byte, error) {
	img, err := imaging.Decode(r)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	if quality <= 0 || quality > 100 {
		quality = 70
	}
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
