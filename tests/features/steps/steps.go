package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"minha-api/models"
)

type apiFeature struct {
	server   http.Handler
	response *httptest.ResponseRecorder
	request  *http.Request
	body     *bytes.Buffer
	apiKey   string
}

func (a *apiFeature) iSetTheApiKeyTo(arg1 string) error {
	a.apiKey = arg1
	return nil
}

func (a *apiFeature) iCreateABookWithTitleAndAuthor(arg1, arg2 string) error {
	book := models.Book{
		Title:  arg1,
		Author: arg2,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(book)
	if err != nil {
		return err
	}
	a.body = &buf
	return nil
}

func (a *apiFeature) iSendARequestToTheEndpoint(arg1, arg2 string) error {
	var err error
	a.request, err = http.NewRequest(arg1, fmt.Sprintf("http://localhost:8080%s", arg2), a.body)
	if err != nil {
		return err
	}
	a.request.Header.Set("Content-Type", "application/json")
	a.request.Header.Set("Accept", "application/json")
	a.request.Header.Set("X-API-Key", a.apiKey)
	return nil
}

func (a *apiFeature) iSendARequest() error {
	rec := httptest.NewRecorder()
	a.server.ServeHTTP(rec, a.request)
	a.response = rec
	return nil
}

func (a *apiFeature) theResponseStatusShouldBe(arg1 int) error {
	if a.response.Code != arg1 {
		return fmt.Errorf("o status code da resposta deveria ser %d mas foi %d", arg1, a.response.Code)
	}
	return nil
}

func (a *apiFeature) theResponseBodyShouldContain(arg1 string) error {
	var body map[string]interface{}
	if err := json.Unmarshal(a.response.Body.Bytes(), &body); err != nil {
		return err
	}
	if _, ok := body[arg1]; !ok {
		return fmt.Errorf("o corpo da resposta não contém o campo %s", arg1)
	}
	return nil
}
