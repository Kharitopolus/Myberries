package tests

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kharitopolus/Myberries/catalog_service/config"
	"github.com/Kharitopolus/Myberries/catalog_service/handlers"
	"github.com/gin-gonic/gin"
)

func prepare(t *testing.T) *gin.Engine {
	t.Helper()
	cfg := config.ConfigGet(".env.test")
	db := handlers.InitDB(cfg.DBUrl)
	if err := db.Exec("TRUNCATE TABLE categories, products RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatal("cannot truncate table", err)
	}
	rt := handlers.SetupRouter(db)

	return rt
}

func assertStatus(t testing.TB, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := rr.Code
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertBody(t testing.TB, rr *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := rr.Body.String()
	if got != want {
		t.Errorf("response body is wrong, got: '%s' want: '%s'", got, want)
	}
}

func assertBodyContains(
	t testing.TB,
	rr *httptest.ResponseRecorder,
	want ...string,
) {
	t.Helper()
	got := rr.Body.String()
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
