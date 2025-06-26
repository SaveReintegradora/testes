package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minha-api/testutils"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	return testutils.SetupRouterWithMocks()
}

func TestGetBooks(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/books", nil)
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Esperado status 200, obteve %d", w.Code)
	}
}

func TestCreateBook(t *testing.T) {
	r := setupRouter()
	json := `{"title":"Livro Teste","author":"Autor"}`
	req, _ := http.NewRequest("POST", "/books", strings.NewReader(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("Esperado status 200 ou 201, obteve %d", w.Code)
	}
}
