package handlers

import (
	"time"

	"github.com/google/uuid"
)

type RegisterReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type RegisterResp struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginResp struct {
	User struct {
		UserID    uuid.UUID `json:"user_id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MeResp struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RefreshResp struct {
	Token string `json:"token"`
}
