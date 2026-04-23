package utility

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// SaveUpload saves the uploaded file and returns the local URL path (e.g. /uploads/category/uuid.ext)
// and the file bytes for further processing (e.g. AI extraction).
func SaveUpload(fileHeader *multipart.FileHeader, category string) (string, []byte, error) {
	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".jpg" // Fallback
	}

	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	categoryDir := filepath.Join(uploadsDir, category)

	if err := os.MkdirAll(categoryDir, os.ModePerm); err != nil {
		return "", nil, fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(categoryDir, newFileName)

	src, err := fileHeader.Open()
	if err != nil {
		return "", nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create local file: %w", err)
	}
	defer dst.Close()

	// Read all bytes while copying
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read uploaded file: %w", err)
	}

	// Write to destination
	if _, err := dst.Write(fileBytes); err != nil {
		return "", nil, fmt.Errorf("failed to write local file: %w", err)
	}

	localUrl := fmt.Sprintf("/uploads/%s/%s", category, newFileName)
	return localUrl, fileBytes, nil
}
