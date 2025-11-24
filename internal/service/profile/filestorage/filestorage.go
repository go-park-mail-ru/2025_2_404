package filestorage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type UseCase struct{
	baseDir	string
}

func New(baseDir string) *UseCase{
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create base dir %s: %v\n", baseDir, err)
	}
	
	return &UseCase{
		baseDir: baseDir,
	}
}

func (u *UseCase) Save(uploadPath string, data io.Reader) error {
	fullPath := filepath.Join(u.baseDir, uploadPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
    	return fmt.Errorf("failed to create directory: %w", err)
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
	fullPath := filepath.Join(u.baseDir, path)
	
    data, err := os.ReadFile(fullPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read image file %q: %w", fullPath, err)
    }
    return data, nil
}

func (u *UseCase) Delete(path string) error {
	if path == ""{
		return nil
	}
	
	fullPath := filepath.Join(u.baseDir, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete file %q: %w", fullPath, err)
	}

	return nil
}