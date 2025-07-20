package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (h UsersHanlersImpl) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.Context().Value(h.userIDCtxKey).(uuid.UUID)

	user, err := h.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user does not exist", err)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		MeResp{
			UserID:    user.UserID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	)
}
