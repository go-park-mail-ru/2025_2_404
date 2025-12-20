// Package readerfile provides utilities for reading files from HTTP requests.
package readerfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"image"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"

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

    _, _, err = image.Decode(bytes.NewReader(fileBytes))
    if err != nil {
        return nil, "", fmt.Errorf("invalid image format: not a valid JPEG, PNG or GIF")
    }


    ext := filepath.Ext(fileHeader.Filename)

    switch ext {
    case ".jpg", ".jpeg", ".png", ".gif":
    default:

        contentType := http.DetectContentType(fileBytes)
        switch contentType {
        case "image/jpeg":
            ext = ".jpg"
        case "image/png":
            ext = ".png"
        case "image/gif":
            ext = ".gif"
        default:
            return nil, "", fmt.Errorf("unsupported image type")
        }
    }

    filename := basePath + uuid.New().String() + ext
    return fileBytes, filename, nil
}