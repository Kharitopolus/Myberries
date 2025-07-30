package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Kharitopolus/Myberries/auth_service/internal/service"
)

func (h UsersHandlersImpl) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	decoder := json.NewDecoder(r.Body)
	rb := RegisterReq{}
	err := decoder.Decode(&rb)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}

	if len(rb.Email) == 0 || len(rb.Password) == 0 {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Unset email or password",
			err,
		)
		return
	}

	user, err := service.RegisterUser(
		r.Context(),
		h.pm.GetHash,
		h.db.CreateUser,
		rb.Email,
		rb.Name,
		rb.Password,
	)
	if err != nil {
		if strings.Contains(
			err.Error(),
			`violates unique constraint "users_email_key"`,
		) {
			respondWithError(
				w,
				http.StatusBadRequest,
				"email already taken",
				err,
			)
			return
		}
		respondWithError(
			w,
			http.StatusInternalServerError,
			"could not create user",
			err,
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		RegisterResp{UserID: user.UserID, Email: user.Email, Name: user.Name},
	)
}
