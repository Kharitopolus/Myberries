package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducts(t *testing.T) {
	rt := prepare(t)

	var rr *httptest.ResponseRecorder

	t.Run("empty list categories", func(t *testing.T) {
		rr = makeReq(rt, "GET", "/categories")
		assert.Equal(t, http.StatusOK, rr.Code)
		b := readBody(t, rr).([]any)
		assert.Equal(t, 0, len(b))
	})

	t.Run("empty list products", func(t *testing.T) {
		rr = makeReq(rt, "GET", "/products")
		assert.Equal(t, http.StatusOK, rr.Code)
		b := readBody(t, rr).([]any)
		assert.Equal(t, 0, len(b))
	})

	t.Run("create category food", func(t *testing.T) {
		rb := fmt.Sprintf(
			`{"name":"%v"}`,
			"food",
		)
		rr = makeReqWithBody(rt, "POST", "/categories", rb)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Body.String(), "id")
		assert.Contains(t, rr.Body.String(), `"name":"food"`)
		assert.Contains(t, rr.Body.String(), "created_at")
	})
	foodCategoryID := readBody(t, rr).(map[string]any)["id"]

	t.Run("create product hotdog", func(t *testing.T) {
		rb := fmt.Sprintf(
			`{"name":"%v","description":"%v","price_usd":%v,"category_id":"%v"}`,
			"hotdog",
			"sausage with bread",
			10.00,
			foodCategoryID,
		)
		rr = makeReqWithBody(rt, "POST", "/products", rb)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Body.String(), "id")
		assert.Contains(t, rr.Body.String(), `"name":"hotdog"`)
		assert.Contains(
			t,
			rr.Body.String(),
			`"description":"sausage with bread"`,
		)
		assert.Contains(t, rr.Body.String(), `"price_usd":".00"`)
		assert.Contains(
			t,
			rr.Body.String(),
			fmt.Sprintf(`"category_id":"%v"`, foodCategoryID),
		)
		assert.Contains(t, rr.Body.String(), "created_at")
	})
}

func makeReq(rt http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, req)

	return rr
}

func makeReqWithBody(
	rt http.Handler,
	method, path, body string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer([]byte(body)))
	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, req)

	return rr
}

func readBody(t *testing.T, rr *httptest.ResponseRecorder) any {
	var data any
	err := json.Unmarshal(rr.Body.Bytes(), &data)
	if err != nil {
		t.Error(err)
	}
	return data
}
