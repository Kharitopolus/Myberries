package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kharitopolus/Myberries/auth_service/internal/config"
	"github.com/Kharitopolus/Myberries/auth_service/internal/handlers"
	"github.com/Kharitopolus/Myberries/auth_service/internal/infrastructure/auth"
	"github.com/Kharitopolus/Myberries/auth_service/internal/infrastructure/database"
	_ "github.com/lib/pq"
)

func TestAuthHandlers(t *testing.T) {
	mux := getMux()

	var rr *httptest.ResponseRecorder

	t.Run("login for not existing user", func(t *testing.T) {
		rr = makeLoginReq(mux, "example@gmail.com", "hackme")
		assertStatus(t, rr.Code, http.StatusUnauthorized)
		assertBody(
			t,
			rr.Body.String(),
			`{"error":"Incorrect email or password"}`,
		)
	})

	t.Run("create user", func(t *testing.T) {
		rr = makeRegisterReq(mux, "example@gmail.com", "example", "hackme")
		assertStatus(t, rr.Code, http.StatusCreated)
		assertBodyContains(
			t,
			rr.Body.String(),
			"user_id",
			`"email":"example@gmail.com"`,
			`"name":"example"`,
			"created_at",
			"updated_at",
		)
	})
	exampleBody := rr.Body.String()

	t.Run("try create user with duplicated email", func(t *testing.T) {
		rr = makeRegisterReq(mux, "example@gmail.com", "example", "hackme")
		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertBody(t, rr.Body.String(), `{"error":"email already taken"}`)
	})

	t.Run("login for existing user", func(t *testing.T) {
		rr = makeLoginReq(mux, "example@gmail.com", "hackme")
		assertStatus(t, rr.Code, http.StatusOK)
		assertBodyContains(t, rr.Body.String(), "access_token", "refresh_token")
	})
	exampleAccessToken := getBodyValue(t, rr.Body.Bytes(), "access_token")
	exampleRefreshToken := getBodyValue(t, rr.Body.Bytes(), "refresh_token")

	t.Run("get me with good token", func(t *testing.T) {
		rr = makeMeReq(mux, exampleAccessToken)
		assertStatus(t, rr.Code, http.StatusOK)
		assertBody(t, rr.Body.String(), exampleBody)
	})

	t.Run("get me with bad token", func(t *testing.T) {
		rr = makeMeReq(mux, exampleRefreshToken)
		assertStatus(t, rr.Code, http.StatusUnauthorized)
		assertBody(t, rr.Body.String(), `{"error":"Couldn't validate JWT"}`)
	})

	t.Run("get new access token with good refresh token", func(t *testing.T) {
		rr = makeRefreshReq(mux, exampleRefreshToken)
		assertStatus(t, rr.Code, http.StatusOK)
		assertBodyContains(t, rr.Body.String(), "token")
	})
	newExampleAccessToken := getBodyValue(t, rr.Body.Bytes(), "token")

	t.Run("get new access token with bad refresh token", func(t *testing.T) {
		rr = makeRefreshReq(mux, exampleAccessToken)
		assertStatus(t, rr.Code, http.StatusUnauthorized)
		assertBody(t, rr.Body.String(), `{"error":"Couldn't validate JWT"}`)
	})

	t.Run("get me with refreshed token", func(t *testing.T) {
		rr = makeMeReq(mux, newExampleAccessToken)
		assertStatus(t, rr.Code, http.StatusOK)
		assertBody(t, rr.Body.String(), exampleBody)
	})
}

func makeRegisterReq(
	mux *http.ServeMux,
	email, name, password string,
) *httptest.ResponseRecorder {
	reqBody := []byte(
		fmt.Sprintf(
			`{"email": "%v", "name": "%v", "password": "%v"}`,
			email,
			name,
			password,
		),
	)
	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		bytes.NewBuffer(
			reqBody,
		),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	return rr
}

func makeLoginReq(
	mux *http.ServeMux,
	email, password string,
) *httptest.ResponseRecorder {
	reqBody := []byte(
		fmt.Sprintf(`{"email": "%v", "password": "%v"}`, email, password),
	)
	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer(
			reqBody,
		),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	return rr
}

func makeMeReq(
	mux *http.ServeMux,
	accessToken string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(
		http.MethodGet,
		"/auth/me",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	return rr
}

func makeRefreshReq(
	mux *http.ServeMux,
	refreshToken string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/refresh",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+refreshToken)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	return rr
}

func getMux() *http.ServeMux {
	cfg := config.Get(".test.env")

	mux := http.NewServeMux()

	db := database.New(cfg.DBUrl)
	db.TruncateUsers(context.Background())
	pm := auth.PasswordManager{}
	tm := auth.TokenManager{
		TokenSecret:         cfg.JWTSecret,
		AccessTokenExpTime:  cfg.AccessTokenExpTimeHours,
		RefreshTokenExpTime: cfg.RefreshTokenExpTimeHours,
	}

	uh := handlers.NewUsersHandlersImpl(db, pm, tm)

	uh.Router(mux)

	return mux
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertBodyContains(t testing.TB, got string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf(
				"response body is wrong, got %q, but want contains %q",
				got,
				w,
			)
		}
	}
}

func getBodyValue(t testing.TB, bodyData []byte, path string) string {
	var body any
	err := json.Unmarshal(bodyData, &body)
	if err != nil {
		t.Errorf("failed to parse response body: %v", bodyData)
	}

	keys := strings.Split(path, ".")
	current := body

	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			t.Errorf("expected JSON object at '%s'", key)
		}

		val, exists := m[key]
		if !exists {
			t.Errorf("key '%s' not found", key)
		}

		current = val
	}

	strVal, ok := current.(string)
	if !ok {
		t.Errorf("value at path '%s' is not a string", path)
	}

	return strVal
}
