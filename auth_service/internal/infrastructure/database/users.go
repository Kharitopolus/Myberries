package database

import (
	"context"
	"database/sql"

	"github.com/Kharitopolus/Myberries/auth_service/internal/entity"
	"github.com/Kharitopolus/Myberries/auth_service/internal/infrastructure/database/sqlc"
	"github.com/google/uuid"
)

func (d DB) CreateUser(
	ctx context.Context,
	email, name, passwordHash string,
) (entity.User, error) {
	user, err := d.sqlc.CreateUser(
		ctx,
		sqlc.CreateUserParams{
			Email:        email,
			Name:         sql.NullString{String: name, Valid: true},
			PasswordHash: passwordHash,
		},
	)
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		UserID: user.UserID,
		Email:  user.Email,
		Name:   user.Email,
	}, nil
}

func (d DB) GetUserByEmail(
	ctx context.Context,
	email string,
) (entity.User, error) {
	user, err := d.sqlc.GetUserByEmail(
		ctx,
		email,
	)
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		UserID:       user.UserID,
		Email:        user.Email,
		Name:         user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (d DB) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (entity.User, error) {
	user, err := d.sqlc.GetUserByID(
		ctx,
		userID,
	)
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		UserID: user.UserID,
		Email:  user.Email,
		Name:   user.Email,
	}, nil
}
