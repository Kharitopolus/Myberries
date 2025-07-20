package handlers

import (
	"net/http"

	"github.com/Kharitopolus/Myberries/auth_service/internal/service"
)

func (h UsersHanlersImpl) Refresh(
	w http.ResponseWriter,
	r *http.Request,
) {
	refreshToken, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	accessToken, err := service.GetNewAccessTokenByRefresh(
		h.tm.ValidateRefreshToken,
		h.tm.MakeAccessToken,
		refreshToken,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't validate token",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshResp{Token: accessToken})
}
