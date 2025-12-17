// Package filestorage provides use case implementations for file storage operations.
package filestorage

import (
	"2025_2_404/internal/service/storage/config"
	"2025_2_404/pkg/globalerrors"
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

type UseCase struct {
	baseDir string
}

func New(cfg *config.Config) *UseCase {
	uc := &UseCase{
		baseDir: cfg.AppConfig.ImgPath,
	}
	slog.Debug(" FileStorage UseCase created", "base_dir", uc.baseDir)
	return uc
}

func (u *UseCase) BaseDir() string {
	return u.baseDir
}

func (u *UseCase) Create(ctx context.Context, imageData []byte, imagePath string) error {
	if imagePath == "" {
		slog.Warn(" Create called with empty path")
		return globalerrors.ErrInvalidPath
	}

	if filepath.IsAbs(imagePath) || filepath.Clean(imagePath) != imagePath {
		slog.Warn(" Suspicious image path (possible traversal)", "path", imagePath)
		return globalerrors.ErrInvalidPath
	}

	fullPath := filepath.Join(u.baseDir, imagePath)
	slog.Debug(" Creating file", "full_path", fullPath, "size", len(imageData))

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		slog.Error(" Failed to create directory", "dir", dir, "error", err)
		return globalerrors.ErrInternal
	}

	file, err := os.Create(fullPath)
	if err != nil {
		slog.Error(" Failed to create file", "path", fullPath, "error", err)
		return globalerrors.ErrFileWrite
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err = file.Write(imageData); err != nil {
		slog.Error(" Failed to write file", "path", fullPath, "error", err)
		return globalerrors.ErrFileWrite
	}

	slog.Debug(" File written successfully", "path", fullPath)
	return nil
}

func (u *UseCase) Get(ctx context.Context, imagePath string) ([]byte, string, error) {
	if imagePath == "" {
		slog.Warn(" Get called with empty path")
		return nil, "", globalerrors.ErrInvalidPath
	}

	// Защита от path traversal
	if filepath.IsAbs(imagePath) || filepath.Clean(imagePath) != imagePath {
		slog.Warn(" Suspicious image path in Get", "path", imagePath)
		return nil, "", globalerrors.ErrInvalidPath
	}

	fullPath := filepath.Join(u.baseDir, imagePath)
	slog.Debug("🔍 Reading file", "full_path", fullPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn(" File not found", "path", fullPath)
			return nil, "", globalerrors.ErrFileNotFound
		}
		slog.Error(" Failed to read file", "path", fullPath, "error", err)
		return nil, "", globalerrors.ErrFileRead
	}

	// Упрощённое определение типа (можно улучшить позже)
	contentType := "image/" + filepath.Ext(fullPath)[1:]
	if contentType == "image/" {
		contentType = "application/octet-stream"
	}

	slog.Debug(" File read successfully", "path", fullPath, "size", len(data))
	return data, contentType, nil
}

func (u *UseCase) Delete(ctx context.Context, imagePath string) error {
	if imagePath == "" {
		slog.Warn(" Delete called with empty path")
		return globalerrors.ErrInvalidPath
	}

	// Защита от path traversal
	if filepath.IsAbs(imagePath) || filepath.Clean(imagePath) != imagePath {
		slog.Warn(" Suspicious image path in Delete", "path", imagePath)
		return globalerrors.ErrInvalidPath
	}

	fullPath := filepath.Join(u.baseDir, imagePath)
	slog.Debug("🗑️ Deleting file", "full_path", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		slog.Warn(" File does not exist, skip delete", "path", fullPath)
		return nil // Успешное "ничего не делать"
	}

	if err := os.Remove(fullPath); err != nil {
		slog.Error(" Failed to delete file", "path", fullPath, "error", err)
		return globalerrors.ErrFileDelete
	}

	slog.Debug(" File deleted successfully", "path", fullPath)
	return nil
}
