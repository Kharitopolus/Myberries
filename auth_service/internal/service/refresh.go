package service

import (
	"github.com/google/uuid"
)

type (
	validateRefreshToken func(string) (uuid.UUID, error)
	makeAccessToken      func(uuid.UUID) (string, error)
)

func GetNewAccessTokenByRefresh(
	vr validateRefreshToken,
	ma makeAccessToken,
	refreshToken string,
) (string, error) {
	userID, err := vr(refreshToken)
	if err != nil {
		return "", err
	}
	accessToken, err := ma(userID)

	return accessToken, err
}
