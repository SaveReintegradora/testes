package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"minha-api/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	endpoints := []struct {
		method string
		url    string
		body   string
	}{
		{"GET", "/books", ""},
		{"GET", "/books/" + id, ""},
		{"PUT", "/books/" + id, `{"title":"Novo Título"}`},
		{"DELETE", "/books/" + id, ""},
	}
	for _, ep := range endpoints {
		var req *http.Request
		if ep.body != "" {
			req, _ = http.NewRequest(ep.method, ep.url, strings.NewReader(ep.body))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, _ = http.NewRequest(ep.method, ep.url, nil)
		}
		req.Header.Set("X-API-Key", "minha-chave-secreta")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK && w.Code != http.StatusCreated && w.Code != http.StatusNoContent {
			t.Errorf("%s %s: esperado status 200/201/204, obteve %d", ep.method, ep.url, w.Code)
		}
	}
}

func TestFilesRoutes(t *testing.T) {
	r := setupRouterWithMocks()
	// Cria arquivo e obtém ID
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fileWriter, err := mw.CreateFormFile("nomeArquivo", "teste.txt")
	if err != nil {
		t.Fatalf("Erro ao criar form file: %v", err)
	}
	io.Copy(fileWriter, bytes.NewReader([]byte("conteudo de teste")))
	mw.Close()

	req, _ := http.NewRequest("POST", "/files/sendFiles", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("POST /files/sendFiles: esperado status 201/200, obteve %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	id, _ := resp["id"].(string)
	if id == "" {
		t.Fatalf("POST /files/sendFiles: resposta sem id")
	}

	endpoints := []struct {
		method      string
		url         string
		body        string
		isMultipart bool
	}{
		{"GET", "/files", "", false},
		{"GET", "/files/" + id, "", false},
		{"PUT", "/files/" + id, `{"name":"Novo Nome"}`, false},
		{"DELETE", "/files/" + id, "", false},
	}
	for _, ep := range endpoints {
		var req *http.Request
		if ep.isMultipart {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			fileWriter, err := w.CreateFormFile("nomeArquivo", "teste.txt")
			if err != nil {
				t.Fatalf("Erro ao criar form file: %v", err)
			}
			io.Copy(fileWriter, bytes.NewReader([]byte(ep.body)))
			w.Close()
			req, _ = http.NewRequest(ep.method, ep.url, &b)
			req.Header.Set("Content-Type", w.FormDataContentType())
		} else if ep.body != "" {
			req, _ = http.NewRequest(ep.method, ep.url, strings.NewReader(ep.body))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, _ = http.NewRequest(ep.method, ep.url, nil)
		}
		req.Header.Set("X-API-Key", "minha-chave-secreta")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK && w.Code != http.StatusCreated && w.Code != http.StatusNoContent {
			t.Errorf("%s %s: esperado status 200/201/204, obteve %d", ep.method, ep.url, w.Code)
		}
	}
}
