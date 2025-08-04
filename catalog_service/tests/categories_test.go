package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kharitopolus/Myberries/catalog_service/handlers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategories(t *testing.T) {
	rt := prepare(t)

	var rr *httptest.ResponseRecorder

	t.Run("empty list", func(t *testing.T) {
		rr = makeListCategoriesReq(rt)
		assertStatus(t, rr, 200)
		cl := readCategories(t, rr)
		assert.Equal(t, 0, len(cl))
	})
	//
	// t.Run("empty body to create category", func(t *testing.T) {
	// 	req := httptest.NewRequest("POST", "/categories", nil)
	// 	rr = httptest.NewRecorder()
	// 	rt.ServeHTTP(rr, req)
	// 	assertStatus(t, rr, 400)
	// 	assertBody(t, rr, `{"error":"'name' and 'description' are required"}`)
	// })

	t.Run("still empty list", func(t *testing.T) {
		rr = makeListCategoriesReq(rt)
		assertStatus(t, rr, 200)
		cl := readCategories(t, rr)
		assert.Equal(t, 0, len(cl))
	})

	t.Run("create lala", func(t *testing.T) {
		rr = makeCreateCatalogReq(rt, "lala")
		assertStatus(t, rr, 201)
		assertBodyContains(t, rr, "id", `"name":"lala"`)
	})

	t.Run("list has lala", func(t *testing.T) {
		rr = makeListCategoriesReq(rt)
		assertStatus(t, rr, 200)
		cl := readCategories(t, rr)
		assert.Equal(t, 1, len(cl))
		assert.Equal(t, cl[0].Name, "lala")
	})

	t.Run("create bebe", func(t *testing.T) {
		rr = makeCreateCatalogReq(rt, "bebe")
		assertStatus(t, rr, 201)
		assertBodyContains(t, rr, "id", `"name":"bebe"`)
	})

	t.Run("list has lala and bebe", func(t *testing.T) {
		rr = makeListCategoriesReq(rt)
		assertStatus(t, rr, 200)
		cl := readCategories(t, rr)
		assert.Equal(t, 2, len(cl))
		assert.Equal(t, cl[0].Name, "lala")
		assert.Equal(t, cl[1].Name, "bebe")
	})
	cl := readCategories(t, rr)
	lala := cl[0]
	bebe := cl[1]

	t.Run("get lala", func(t *testing.T) {
		rr = makeGetCategoryReq(rt, lala.ID)
		assertStatus(t, rr, 200)
		c := readCategory(t, rr)
		assert.Equal(t, "lala", c.Name)
	})

	t.Run("get bebe", func(t *testing.T) {
		rr = makeGetCategoryReq(rt, bebe.ID)
		assertStatus(t, rr, 200)
		c := readCategory(t, rr)
		assert.Equal(t, "bebe", c.Name)
	})

	t.Run("update lala to momo", func(t *testing.T) {
		rr = makeUpdateCategoryReq(rt, lala.ID, "momo")
		assertStatus(t, rr, 200)
	})

	t.Run("get lala but now momo", func(t *testing.T) {
		rr = makeGetCategoryReq(rt, lala.ID)
		assertStatus(t, rr, 200)
		c := readCategory(t, rr)
		assert.Equal(t, "momo", c.Name)
	})

	t.Run("lala became momo also in list", func(t *testing.T) {
		rr = makeListCategoriesReq(rt)
		assertStatus(t, rr, 200)
		cl := readCategories(t, rr)
		assert.Equal(t, 2, len(cl))
		assert.Equal(t, cl[0].Name, "momo")
		assert.Equal(t, cl[1].Name, "bebe")
	})

	t.Run("delete lala", func(t *testing.T) {
		t.Log("lala.ID:", lala.ID)
		rr = makeDeleteCategoryReq(rt, lala.ID)
		assertStatus(t, rr, 200)
		t.Log("delete lala body:", rr.Body.String())
	})

	t.Run("get nonexisting lala", func(t *testing.T) {
		rr = makeGetCategoryReq(rt, lala.ID)
		assertStatus(t, rr, 404)
		assert.Equal(t, rr.Body.String(), `{"error":"category does not exist"}`)
	})

	t.Run("update nonexisting lala", func(t *testing.T) {
		rr = makeUpdateCategoryReq(rt, lala.ID, "momo")
		assertStatus(t, rr, 404)
		assert.Equal(t, rr.Body.String(), `{"error":"category does not exist"}`)
	})

	t.Run("list has only bebe", func(t *testing.T) {
		rr = makeListCategoriesReq(rt)
		assertStatus(t, rr, 200)
		cl = readCategories(t, rr)
		assert.Equal(t, 1, len(cl))
		assert.Equal(t, cl[0].Name, "bebe")
	})

	t.Run(
		"deleting category with nonexisting id return 404",
		func(t *testing.T) {
			rr = makeDeleteCategoryReq(rt, lala.ID)
			assertStatus(t, rr, 404)
			assertBody(t, rr, `{"error":"category does not exist"}`)
		},
	)

	t.Run("delete bebe", func(t *testing.T) {
		rr = makeDeleteCategoryReq(rt, bebe.ID)
		assertStatus(t, rr, 200)
	})

	t.Run("list empty after deleting", func(t *testing.T) {
		rr = makeListCategoriesReq(rt)
		assertStatus(t, rr, 200)
		cl = readCategories(t, rr)
		assert.Equal(t, 0, len(cl))
	})
}

func makeGetCategoryReq(
	rt http.Handler,
	id uuid.UUID,
) *httptest.ResponseRecorder {
	path := "/categories" + "/" + id.String()
	req := httptest.NewRequest("GET", path, nil)
	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, req)

	return rr
}

func makeDeleteCategoryReq(
	rt http.Handler,
	id uuid.UUID,
) *httptest.ResponseRecorder {
	path := "/categories" + "/" + id.String()
	req := httptest.NewRequest("DELETE", path, nil)
	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, req)

	return rr
}

func makeListCategoriesReq(rt http.Handler) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/categories", nil)
	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, req)

	return rr
}

func makeUpdateCategoryReq(
	rt http.Handler,
	id uuid.UUID,
	name string,
) *httptest.ResponseRecorder {
	path := "/categories" + "/" + id.String()
	reqBody := []byte(
		fmt.Sprintf(`{"name": "%v"}`, name),
	)
	req := httptest.NewRequest("PUT", path, bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, req)

	return rr
}

func makeCreateCatalogReq(
	rt http.Handler,
	name string,
) *httptest.ResponseRecorder {
	reqBody := []byte(
		fmt.Sprintf(`{"name": "%v"}`, name),
	)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, req)

	return rr
}

func readCategory(
	t *testing.T,
	rr *httptest.ResponseRecorder,
) handlers.Category {
	var c handlers.Category
	err := json.Unmarshal(rr.Body.Bytes(), &c)
	if err != nil {
		t.Errorf("got bad categories list: %v", rr.Body.String())
	}
	return c
}

func readCategories(
	t *testing.T,
	rr *httptest.ResponseRecorder,
) []handlers.Category {
	var cl []handlers.Category
	err := json.Unmarshal(rr.Body.Bytes(), &cl)
	if err != nil {
		t.Errorf("got bad categories list: %v", rr.Body.String())
	}
	return cl
}
