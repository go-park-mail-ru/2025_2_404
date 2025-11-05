package filestorage

import (
	"2025_2_404/internal/config"
	"fmt"
	"io"
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

func (u *UseCase) Save(uploadPath string, data io.Reader) error {
	fullPath := filepath.Join(u.baseDir, uploadPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
    return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (u *UseCase) ReadImageAsBytes(path string) ([]byte, error) {
	if path == ""{
		return nil, nil
	}
	fullpath := u.baseDir + path
    data, err := os.ReadFile(fullpath)
    if err != nil {
        return nil, fmt.Errorf("failed to read image file %q: %w", path, err)
    }
    return data, nil
}