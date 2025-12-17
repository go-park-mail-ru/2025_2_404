// Package profile provides use case implementations for managing user profiles.
package profile

import (
	modelclient "2025_2_404/internal/service/profile/domain"
	"context"
	"fmt"
)

func (u *UseCase) ShowBalance(ctx context.Context, clientID modelclient.ID) (uint32, error) {
	return u.repo.ShowBalance(ctx, clientID)
}

func (u *UseCase) AddBalance(ctx context.Context, clientID modelclient.ID, addAmount uint32) error {
	return u.repo.AddBalance(ctx, clientID, addAmount)
}

func (u *UseCase) SubtractBalance(ctx context.Context, clientID modelclient.ID, subAmount uint32) error {
	balanc, err := u.repo.ShowBalance(ctx, clientID)
	if err != nil {
		return err
	}
	if balanc < subAmount {
		return fmt.Errorf("Subtract more balance")
	}
	return u.repo.SubtractBalance(ctx, clientID, subAmount)
}
