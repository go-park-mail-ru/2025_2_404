package filestorage

import (
	"2025_2_404/internal/service/storage/config"
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

type UseCase struct{
	baseDir	string
}

// func New(cfg *config.Config) *UseCase{
// 	return &UseCase{
// 		baseDir: cfg.AppConfig.ImgPath,
// 	}
// }

func New(cfg *config.Config) *UseCase {
	uc := &UseCase{
		baseDir: cfg.AppConfig.ImgPath,
	}
	slog.Debug("📁 FileStorage UseCase created", "base_dir", uc.baseDir)
	return uc
}

func (u *UseCase) BaseDir() string {
	return u.baseDir
}

func (u *UseCase) Create(ctx context.Context, imageData []byte, imagePath string) error {
	fullPath := filepath.Join(u.baseDir, imagePath)
	slog.Debug("📝 Creating file", "full_path", fullPath, "size", len(imageData))

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		slog.Error("❌ Failed to create directory", "dir", dir, "error", err)
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		slog.Error("❌ Failed to create file", "path", fullPath, "error", err)
		return err
	}
	defer file.Close()

	if _, err = file.Write(imageData); err != nil {
		slog.Error("❌ Failed to write file", "path", fullPath, "error", err)
		return err
	}

	slog.Debug("✅ File written successfully", "path", fullPath)
	return nil
}

func (u *UseCase) Get(ctx context.Context, imagePath string) ([]byte, string, error) {
	if imagePath == "" {
		slog.Warn("⚠️ Get called with empty path")
		return nil, "", nil
	}

	fullPath := filepath.Join(u.baseDir, imagePath)
	slog.Debug("🔍 Reading file", "full_path", fullPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		slog.Error("❌ Failed to read file", "path", fullPath, "error", err)
		return nil, "", err
	}

	// Простой способ определить Content-Type (можно улучшить)
	contentType := "application/octet-stream"
	if ext := filepath.Ext(fullPath); ext != "" {
		// TODO: использовать http.DetectContentType или mime.TypeByExtension
		contentType = "image/" + ext[1:]
	}

	slog.Debug("✅ File read successfully", "path", fullPath, "size", len(data))
	return data, contentType, nil
}

func (u *UseCase) Delete(ctx context.Context, imagePath string) error {
	if imagePath == "" {
		slog.Warn("⚠️ Delete called with empty path")
		return nil
	}

	fullPath := filepath.Join(u.baseDir, imagePath)
	slog.Debug("🗑️ Deleting file", "full_path", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		slog.Warn("⚠️ File does not exist, skip delete", "path", fullPath)
		return nil
	}

	if err := os.Remove(fullPath); err != nil {
		slog.Error("❌ Failed to delete file", "path", fullPath, "error", err)
		return err
	}

	slog.Debug("✅ File deleted successfully", "path", fullPath)
	return nil
}