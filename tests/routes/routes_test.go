package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minha-api/tests/testutils"
)

func setupRouter() http.Handler {
	return testutils.SetupRouter()
}

func setupRouterWithMocks() http.Handler {
	return testutils.SetupRouterWithMocks()
}

func TestBooksRoutes(t *testing.T) {
	r := setupRouterWithMocks()
	// Cria livro e obtém ID
	postBody := `{"title":"Livro Teste","author":"Autor"}`
	req, _ := http.NewRequest("POST", "/books", strings.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("POST /books: esperado status 201/200, obteve %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	id, _ := resp["id"].(string)
	if id == "" {
		t.Fatalf("POST /books: resposta sem id")
	}
}
