package modules

import (
	modeluser "2025_2_404/internal/service/profile/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type key string

const(
	UserIDKey key = "userID"
)

func Set(ctx context.Context, userID modeluser.ID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func Get(ctx context.Context) (modeluser.ID, error) {
	userID, ok := ctx.Value(UserIDKey).(modeluser.ID)
	if !ok {
		return modeluser.ID(uuid.Nil), fmt.Errorf("problem in get")
	}
	return userID, nil
}