package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"minha-api/models"
	"minha-api/tests/testutils"

	. "github.com/cucumber/godog"
	"github.com/google/uuid"
)

type apiFeature struct {
	server     http.Handler
	response   *httptest.ResponseRecorder
	request    *http.Request
	body       *bytes.Buffer
	apiKey     string
	lastBookID string // Armazena o último ID de livro criado
	lastFileID string // Armazena o último ID de arquivo criado
}

func (a *apiFeature) reset() {
	a.body = &bytes.Buffer{}
	a.response = httptest.NewRecorder()
	a.apiKey = "minha-chave-secreta"
}

func (a *apiFeature) iAmAuthenticatedWithAPIKey(key string) error {
	a.apiKey = key
	return nil
}

func (a *apiFeature) iMakeARequestToWithBody(method, path, body string) error {
	// Substitui 'ultimo' pelo ID real, se necessário
	if path == "/books/ultimo" {
		path = "/books/" + a.lastBookID
	}
	if path == "/files/ultimo" {
		path = "/files/" + a.lastFileID
	}
	a.body = bytes.NewBufferString(body)
	req, err := http.NewRequest(method, path, a.body)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", a.apiKey)
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}
	a.request = req
	a.response = httptest.NewRecorder()
	a.server.ServeHTTP(a.response, req)
	// Captura o ID retornado se for criação de livro ou arquivo
	if method == "POST" && (path == "/books" || path == "/files/sendFiles") && a.response.Code == 201 {
		var resp map[string]interface{}
		if err := json.Unmarshal(a.response.Body.Bytes(), &resp); err == nil {
			if id, ok := resp["id"]; ok {
				if path == "/books" {
					a.lastBookID = fmt.Sprintf("%v", id)
				} else {
					a.lastFileID = fmt.Sprintf("%v", id)
				}
			}
		}
	}
	return nil
}

func (a *apiFeature) iMakeARequestTo(method, path string) error {
	return a.iMakeARequestToWithBody(method, path, "")
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if a.response.Code != code {
		return fmt.Errorf("esperado status %d, obteve %d", code, a.response.Code)
	}
	return nil
}

func (a *apiFeature) theResponseShouldContain(field, value string) error {
	var resp map[string]interface{}
	if err := json.Unmarshal(a.response.Body.Bytes(), &resp); err != nil {
		return err
	}
	if v, ok := resp[field]; !ok || v != value {
		return fmt.Errorf("esperado %s = %s, obteve %v", field, value, v)
	}
	return nil
}

func InitializeScenario(ctx *ScenarioContext) {
	f := &apiFeature{}
	f.server = testutils.SetupRouterWithMocks() // Função que retorna o handler com mocks
	ctx.Before(func(ctx context.Context, sc *Scenario) (context.Context, error) {
		f.reset()
		if BookRepoMock != nil {
			BookRepoMock.Reset()
		}
		if FileRepoMock != nil {
			FileRepoMock.Reset()
		}
		return ctx, nil
	})
	ctx.Step(`^que estou autenticado com a API Key "([^"]*)"$`, f.iAmAuthenticatedWithAPIKey)
	ctx.Step(`^faço uma requisição (GET|POST|PUT|DELETE) para "([^"]*)"$`, f.iMakeARequestTo)
	ctx.Step(`^faço uma requisição (POST|PUT) para "([^"]*)" com o corpo:$`, f.iMakeARequestToWithBody)
	ctx.Step(`^a resposta deve ter status (\d+)$`, f.theResponseCodeShouldBe)
	ctx.Step(`^a resposta deve conter uma lista de livros$`, func() error {
		var arr []interface{}
		err := json.Unmarshal(f.response.Body.Bytes(), &arr)
		if err != nil {
			return fmt.Errorf("esperado um array JSON, erro: %w", err)
		}
		if len(arr) == 0 {
			return fmt.Errorf("lista de livros está vazia")
		}
		return nil
	})
	ctx.Step(`^a resposta deve conter o livro com ID ".*"$`, func() error {
		var resp map[string]interface{}
		if err := json.Unmarshal(f.response.Body.Bytes(), &resp); err != nil {
			return err
		}
		if v, ok := resp["id"]; !ok || fmt.Sprintf("%v", v) != f.lastBookID {
			return fmt.Errorf("esperado id = %s, obteve %v", f.lastBookID, v)
		}
		return nil
	})
	ctx.Step(`^a resposta deve conter o livro com título "([^"]*)"$`, func(title string) error {
		return f.theResponseShouldContain("title", title)
	})
	ctx.Step(`^a resposta deve conter uma lista de arquivos$`, func() error {
		var arr []interface{}
		err := json.Unmarshal(f.response.Body.Bytes(), &arr)
		if err != nil {
			return fmt.Errorf("esperado um array JSON, erro: %w", err)
		}
		if len(arr) == 0 {
			return fmt.Errorf("lista de arquivos está vazia")
		}
		return nil
	})
	ctx.Step(`^a resposta deve conter o nome do arquivo "([^"]*)"$`, func(name string) error {
		return f.theResponseShouldContain("fileName", name)
	})
	ctx.Step(`^a resposta deve conter o arquivo com ID ".*"$`, func() error {
		var resp map[string]interface{}
		if err := json.Unmarshal(f.response.Body.Bytes(), &resp); err != nil {
			return err
		}
		if v, ok := resp["id"]; !ok || fmt.Sprintf("%v", v) != f.lastFileID {
			return fmt.Errorf("esperado id = %s, obteve %v", f.lastFileID, v)
		}
		return nil
	})
	ctx.Step(`^a resposta deve conter o arquivo com nome "([^"]*)"$`, func(name string) error {
		return f.theResponseShouldContain("fileName", name)
	})
	ctx.Step(`^que existe um livro com ID "([^"]*)"$`, func(id string) error {
		if BookRepoMock == nil {
			return fmt.Errorf("BookRepoMock não inicializado")
		}
		var realID string
		if _, err := uuid.Parse(id); err != nil {
			realID = uuid.New().String()
		} else {
			realID = id
		}
		book := models.Book{
			ID:     realID,
			Title:  "Livro " + realID,
			Author: "Autor",
		}
		BookRepoMock.Books[realID] = book
		f.lastBookID = realID
		return nil
	})
	ctx.Step(`^que existe um arquivo com ID "([^"]*)"$`, func(id string) error {
		if FileRepoMock == nil {
			return fmt.Errorf("FileRepoMock não inicializado")
		}
		var realID string
		if _, err := uuid.Parse(id); err != nil {
			realID = uuid.New().String()
		} else {
			realID = id
		}
		file := models.FileProcess{
			ID:       realID,
			FileName: "arquivo_" + realID + ".txt",
			FilePath: "https://mock-s3.local/arquivo_" + realID + ".txt",
			Status:   "recebido",
		}
		FileRepoMock.Files[realID] = file
		f.lastFileID = realID
		return nil
	})
	// Steps para usar o último ID criado em requisições subsequentes
	ctx.Step(`^faço uma requisição (GET|PUT|DELETE) para "/books/ultimo"$`, func(method string) error {
		id := f.lastBookID
		if id == "" {
			return fmt.Errorf("nenhum livro criado para usar o ID")
		}
		path := "/books/" + id
		return f.iMakeARequestTo(method, path)
	})
	ctx.Step(`^faço uma requisição (PUT) para "/books/ultimo" com o corpo:$`, func(method, body string) error {
		id := f.lastBookID
		if id == "" {
			return fmt.Errorf("nenhum livro criado para usar o ID")
		}
		path := "/books/" + id
		return f.iMakeARequestToWithBody(method, path, body)
	})
	ctx.Step(`^faço uma requisição (GET|PUT|DELETE) para "/files/ultimo"$`, func(method string) error {
		id := f.lastFileID
		if id == "" {
			return fmt.Errorf("nenhum arquivo criado para usar o ID")
		}
		path := "/files/" + id
		return f.iMakeARequestTo(method, path)
	})
	ctx.Step(`^faço uma requisição (PUT) para "/files/ultimo" com o corpo:$`, func(method, body string) error {
		id := f.lastFileID
		if id == "" {
			return fmt.Errorf("nenhum arquivo criado para usar o ID")
		}
		path := "/files/" + id
		return f.iMakeARequestToWithBody(method, path, body)
	})
	ctx.Step(`^a resposta deve conter o livro com ID "\$\{ultimo_id\}"$`, func() error {
		var resp map[string]interface{}
		if err := json.Unmarshal(f.response.Body.Bytes(), &resp); err != nil {
			return err
		}
		if v, ok := resp["id"]; !ok || fmt.Sprintf("%v", v) != f.lastBookID {
			return fmt.Errorf("esperado id = %s, obteve %v", f.lastBookID, v)
		}
		return nil
	})
	ctx.Step(`^a resposta deve conter o arquivo com ID "\$\{ultimo_id\}"$`, func() error {
		var resp map[string]interface{}
		if err := json.Unmarshal(f.response.Body.Bytes(), &resp); err != nil {
			return err
		}
		if v, ok := resp["id"]; !ok || fmt.Sprintf("%v", v) != f.lastFileID {
			return fmt.Errorf("esperado id = %s, obteve %v", f.lastFileID, v)
		}
		return nil
	})
	ctx.Step(`^a resposta deve conter o campo "([^"]*)" com valor "([^"]*)"$`, func(field, value string) error {
		var resp map[string]interface{}
		if err := json.Unmarshal(f.response.Body.Bytes(), &resp); err != nil {
			return err
		}
		v, ok := resp[field]
		if !ok {
			return fmt.Errorf("campo %s não encontrado na resposta", field)
		}
		if fmt.Sprintf("%v", v) != value {
			return fmt.Errorf("esperado %s = %s, obteve %v", field, value, v)
		}
		return nil
	})
	// Adicione outros steps conforme necessário
}
