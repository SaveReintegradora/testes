package controllers_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"minha-api/models"
	"minha-api/repositories"
	"minha-api/routes"
	"minha-api/tests/testutils"
	"minha-api/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFiles(t *testing.T) {
	r := testutils.SetupRouterWithMocks()
	req, _ := http.NewRequest("GET", "/files", nil)
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Esperado status 200, obteve %d", w.Code)
	}
}

func TestSendFiles(t *testing.T) {
	r := testutils.SetupRouterWithMocks()

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
}

func TestDownloadFile(t *testing.T) {
	bookRepo := repositories.NewBookRepositoryMock()
	fileRepo := repositories.NewFileProcessRepositoryMock()
	// Garante que o arquivo de teste tem FilePath preenchido e ID válido (UUID)
	fileRepo.Files["2bce990f-5e9d-40df-8133-6b323fec8cbe"] = models.FileProcess{
		ID:         "2bce990f-5e9d-40df-8133-6b323fec8cbe",
		FileName:   "mock.txt",
		FilePath:   "https://mock-s3.local/mock.txt",
		Status:     "recebido",
		ReceivedAt: fileRepo.Files["1"].ReceivedAt, // mantém o timestamp original se existir
	}
	s3mock := &utils.MockS3Uploader{}
	s3presign := &utils.MockS3Presigner{}
	r := routes.SetupRoutesWithReposAndS3(bookRepo, fileRepo, s3mock, s3presign)
	fileID := "2bce990f-5e9d-40df-8133-6b323fec8cbe"

	req, _ := http.NewRequest("GET", "/files/"+fileID+"/download", nil)
	req.Header.Set("X-API-Key", "minha-chave-secreta")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Esperado status 302 (redirect), obteve %d", w.Code)
	}
	location := w.Header().Get("Location")
	if location == "" {
		t.Errorf("Esperado header Location com URL do arquivo, mas veio vazio")
	}
}
