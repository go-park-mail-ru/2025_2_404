// Package readerfile provides utilities for reading files from HTTP requests.
package readerfile

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
)

func IsMissingFileError(err error) bool {
	return errors.Is(err, http.ErrMissingFile)
}

func ExtractImage(r *http.Request, basePath string, formFieldName string) ([]byte, string, error) {
	_, fileHeader, err := r.FormFile(formFieldName)
	if err != nil {
		if IsMissingFileError(err) {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("invalid image: %w", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", fmt.Errorf("cannot open image: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("cannot read image: %w", err)
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := basePath + uuid.New().String() + ext

	return fileBytes, filename, nil
}
