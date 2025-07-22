package handlers

import (
	"context"
	"net/http"

	"github.com/Kharitopolus/Myberries/auth_service/internal/entity"
	"github.com/Kharitopolus/Myberries/auth_service/internal/middleware"
	"github.com/google/uuid"
)

type usersDB interface {
	CreateUser(
		ctx context.Context,
		email, name, password string,
	) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

type passwordManager interface {
	GetHash(password string) (string, error)
	CheckHash(hash, password string) error
}

type tokenManager interface {
	MakeAccessToken(userID uuid.UUID) (string, error)
	MakeRefreshToken(userID uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
	ValidateRefreshToken(tokenString string) (uuid.UUID, error)
}

type UsersHandlersImpl struct {
	db           usersDB
	pm           passwordManager
	tm           tokenManager
	userIDCtxKey string
}

func NewUsersHandlersImpl(
	db usersDB,
	pm passwordManager,
	tm tokenManager,
) UsersHandlersImpl {
	return UsersHandlersImpl{
		db: db,
		pm: pm,
		tm: tm,
	}
}

func (h *UsersHandlersImpl) Router(mux *http.ServeMux) {
	am := middleware.AuthMiddleware{
		ValidateAccessToken: h.tm.ValidateAccessToken,
	}
	h.userIDCtxKey = middleware.UserIDCtxKey

	mux.Handle("GET /auth/me", am.Authenticate(http.HandlerFunc(h.Me)))
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)
}
