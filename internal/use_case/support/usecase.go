package support

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"io"
)

type fileStorageI interface{
	Save(uploadPath string, data io.Reader) error
	ReadImageAsBytes(path string) ([]byte, error)
	Delete(path string) error
}

