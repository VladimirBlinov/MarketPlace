package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServer_HandleAllProducts(t *testing.T) {
	s := New(NewConfig())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/allproducts", nil)

	s.HandleGetAllProducts().ServeHTTP(rec, req)
	assert.Equal(t, rec.Body.String(), "All Products here")
}
