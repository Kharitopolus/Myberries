package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID
	Email        string
	Name         string
	PasswordHash string
	RoleID       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
