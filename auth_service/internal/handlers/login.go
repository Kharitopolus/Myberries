package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Kharitopolus/Myberries/auth_service/internal/service"
)

func (h UsersHanlersImpl) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginReq{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}
	user, err := service.CheckCredentials(
		r.Context(),
		h.db.GetUserByEmail,
		h.pm.CheckHash,
		params.Email,
		params.Password,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
			err,
		)
		return
	}

	accessToken, err := h.tm.MakeAccessToken(user.UserID)
	refreshToken, err := h.tm.MakeRefreshToken(user.UserID)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create tokens",
			err,
		)
		return
	}

	var lr LoginResp

	lr.User.UserID = user.UserID
	lr.User.Email = user.Email
	lr.User.Name = user.Name
	lr.User.CreatedAt = user.CreatedAt
	lr.User.UpdatedAt = user.UpdatedAt

	lr.AccessToken = accessToken
	lr.RefreshToken = refreshToken

	respondWithJSON(w, http.StatusOK, lr)
}
