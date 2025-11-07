package utils

import (
	"path/filepath"
	"strings"
)

func GetMimeTypeFromName(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
