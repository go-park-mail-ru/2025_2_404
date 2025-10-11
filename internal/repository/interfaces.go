package repository

import (
	"2025_2_404/internal/domain/models/user"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *user.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}
