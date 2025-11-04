package profile

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
)

type repositoryI interface{
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error)
	Update(ctx context.Context, client modeluser.User) error
	Delete(ctx context.Context, clientID modeluser.ID) error
}

type fileStorageI interface{
	Save(uploadPath string, data io.Reader) error
	ReadImageAsBytes(path string) ([]byte, error)
	Delete(path string) error
}

type UseCase struct{
	repo repositoryI
	storage fileStorageI
}

func New(repo repositoryI, storage fileStorageI) *UseCase{
	return &UseCase{
		repo: repo,
		storage: storage,
	}
}

func (u *UseCase) Update(ctx context.Context, client modeluser.User, file io.Reader, ext string) error {
	filename := uuid.New().String() + ext

	uploadPath := filepath.Join("ad/", filename)

	if err := u.storage.Save(uploadPath, file); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	client.ImagePath = uploadPath


	return u.repo.Update(ctx, client)
}

func (u *UseCase) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, []byte, error){
	
	client, err := u.repo.Show(ctx, clientID)
	if err != nil {
		return modeluser.User{}, nil, fmt.Errorf("problem with show in repo")
	}
	imgBytes, err := u.storage.ReadImageAsBytes(client.ImagePath)
	if err != nil {
		return  modeluser.User{}, nil, fmt.Errorf("problem in convert")
	}

	return client, imgBytes, nil
}

func (u *UseCase) Delete(ctx context.Context, clientID modeluser.ID) error {
	client, err := u.repo.Show(ctx, clientID)
	if err != nil {
		return fmt.Errorf("failed to get user for deletion: %w", err)
	}

	if client.ImagePath != "" {
		if err := u.storage.Delete(client.ImagePath); err != nil {
			return fmt.Errorf("failed to delete user photo: %w", err)
		}
	}

	if err := u.repo.Delete(ctx, clientID); err != nil {
		return fmt.Errorf("failed to delete user from repository: %w", err)
	}

	return nil
}