package filestorage

import (
	"2025_2_404/internal/service/storage/config"
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type UseCase struct{
	baseDir	string
}

func New(cfg *config.Config) *UseCase{
	return &UseCase{
		baseDir: cfg.AppConfig.ImgPath,
	}
}

func (u *UseCase) Create(ctx context.Context, imageData []byte, imagePath string) error {

	if err := os.MkdirAll(filepath.Dir(imagePath), os.ModePerm); err != nil {
    return err
	}

	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(imageData)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (u *UseCase) Get(ctx context.Context, imagePath string) ([]byte, string, error) {
	if imagePath == ""{
		return nil, "", nil
	}

    data, err := os.ReadFile(imagePath)
    if err != nil {
        return nil,"", fmt.Errorf("failed to read image file %q: %w", imagePath, err)
    }
    return data,"", nil
}

func (u *UseCase) Delete(ctx context.Context, imagePath string) error {
	if imagePath == ""{
		return nil
	}

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(imagePath)
	if err != nil {
		return fmt.Errorf("failed to delete file %q: %w", imagePath, err)
	}

	return nil
}