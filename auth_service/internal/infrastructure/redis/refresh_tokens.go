package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r RClient) SetRefreshToken(ctx context.Context, token string, user_id uuid.UUID) error {
	err := r.client.Set(ctx, token, user_id.String(), time.Duration(r.refreshToeknExpTimeHOurs)*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RClient) GetUserIDByRefresToken(ctx context.Context, token string) (uuid.UUID, error) {
	val, err := r.client.Get(ctx, token).Result()
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(val)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, err
}

func (r RClient) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.client.Del(ctx, token).Err()
	if err != nil {
		return err
	}
	return nil
}
