package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Storage struct {
	UploadDir string
}

func NewStorage(uploadDir string) *Storage {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(err)
	}
	return &Storage{UploadDir: uploadDir}
}

func (s *Storage) SaveImage(file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	if !isAllowedImageType(ext) {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(s.UploadDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return the relative path
	return filename, nil
}

func isAllowedImageType(ext string) bool {
	ext = strings.ToLower(ext)
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return allowedTypes[ext]
}
