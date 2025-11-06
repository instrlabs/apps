package utils

import (
	"path/filepath"
	"strings"
)

func GetMimeTypeFromName(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".tiff":
		return "image/tiff"
	default:
		return "application/octet-stream"
	}
}
