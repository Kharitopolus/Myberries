package service

import (
	"context"

	"github.com/Kharitopolus/Myberries/auth_service/internal/entity"
)

type (
	getHash    func(password string) (string, error)
	createUser func(
		ctx context.Context,
		email, name, passwordHash string,
	) (entity.User, error)
)

func RegisterUser(
	ctx context.Context,
	gh getHash,
	cu createUser,
	email, name, password string,
) (entity.User, error) {
	passwordHash, err := gh(password)
	if err != nil {
		return entity.User{}, err
	}

	user, err := cu(ctx, email, name, passwordHash)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
