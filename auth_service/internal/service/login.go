package service

import (
	"context"

	"github.com/Kharitopolus/Myberries/auth_service/internal/entity"
)

type (
	getUserByEmail    func(ctx context.Context, email string) (entity.User, error)
	checkPasswordHash func(password, passwordHash string) error
)

func CheckCredentials(
	ctx context.Context,
	gu getUserByEmail,
	ch checkPasswordHash,
	email, password string,
) (entity.User, error) {
	user, err := gu(ctx, email)
	if err != nil {
		return entity.User{}, err
	}
	err = ch(user.PasswordHash, password)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
