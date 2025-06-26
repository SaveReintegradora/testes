package controllers_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFiles(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/files", nil)
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Esperado status 200, obteve %d", w.Code)
	}
}

func TestSendFiles(t *testing.T) {
	r := setupRouter()

	// Cria um buffer multipart para simular upload de arquivo
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fileWriter, err := w.CreateFormFile("nomeArquivo", "teste.txt")
	if err != nil {
		t.Fatalf("Erro ao criar form file: %v", err)
	}
	io.Copy(fileWriter, bytes.NewReader([]byte("conteudo de teste")))
	w.Close()

	req, _ := http.NewRequest("POST", "/files/sendFiles", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK && resp.Code != http.StatusCreated {
		t.Errorf("Esperado status 200 ou 201, obteve %d", resp.Code)
	}
}
