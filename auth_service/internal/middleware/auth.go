package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const UserIDCtxKey = "userID"

type validateAccessToken func(string) (uuid.UUID, error)

type AuthMiddleware struct {
	ValidateAccessToken validateAccessToken
}

func (am AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Replace(
			r.Header.Get("Authorization"),
			fmt.Sprintf("%s ", "Bearer"),
			"",
			1,
		)

		userID, err := am.ValidateAccessToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			dat, _ := json.Marshal(struct {
				Error string `json:"error"`
			}{Error: "Couldn't validate JWT"})
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(dat)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDCtxKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
