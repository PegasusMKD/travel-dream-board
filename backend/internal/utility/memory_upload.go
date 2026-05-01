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

// MemoriesDir returns the on-disk directory where memory image files are stored.
// Defaults to "./memories" but can be overridden with MEMORIES_DIR (e.g. a Railway
// volume mounted at /memories). Memories are NOT statically served — they're
// streamed through an authenticated endpoint that verifies board access.
func MemoriesDir() string {
	dir := os.Getenv("MEMORIES_DIR")
	if dir == "" {
		dir = "./memories"
	}
	return dir
}

// SaveMemoryUpload writes the uploaded file under MEMORIES_DIR/<uuid>.<ext> and
// returns the absolute path on disk. The path is meant to be persisted in the DB
// `image_url` column; clients fetch the bytes via `GET /api/v1/memories/:uuid/image`.
func SaveMemoryUpload(fileHeader *multipart.FileHeader) (string, error) {
	dir := MemoriesDir()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create memories dir: %w", err)
	}

	name := uuid.New().String() + ext
	full := filepath.Join(dir, name)

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(full)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return full, nil
}
