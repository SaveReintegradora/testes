package repositories

import (
	"errors"
	"minha-api/models"
	"time"
)

type BookRepositoryMock struct {
	Books map[string]models.Book
}

// Garante que BookRepositoryMock implementa BookRepositoryInterface
var _ BookRepositoryInterface = (*BookRepositoryMock)(nil)

func NewBookRepositoryMock() *BookRepositoryMock {
	return &BookRepositoryMock{
		Books: map[string]models.Book{
			"1": {ID: "1", Title: "Livro 1", Author: "Autor 1", CreatedAt: time.Now()},
		},
	}
}

func (m *BookRepositoryMock) GetAll() ([]models.Book, error) {
	books := make([]models.Book, 0, len(m.Books))
	for _, b := range m.Books {
		books = append(books, b)
	}
	return books, nil
}

func (m *BookRepositoryMock) GetByID(id string) (*models.Book, error) {
	if b, ok := m.Books[id]; ok {
		return &b, nil
	}
	return nil, errors.New("not found")
}

func (m *BookRepositoryMock) Create(book *models.Book) error {
	m.Books[book.ID] = *book
	return nil
}

func (m *BookRepositoryMock) Update(book *models.Book) error {
	if _, ok := m.Books[book.ID]; ok {
		m.Books[book.ID] = *book
		return nil
	}
	return errors.New("not found")
}

func (m *BookRepositoryMock) Delete(id string) error {
	if _, ok := m.Books[id]; ok {
		delete(m.Books, id)
		return nil
	}
	return errors.New("not found")
}

// Reset limpa o estado do mock
func (m *BookRepositoryMock) Reset() {
	m.Books = map[string]models.Book{}
	// Adiciona um livro padrão para testes de listagem
	m.Books["1"] = models.Book{
		ID:     "1",
		Title:  "Livro 1",
		Author: "Autor 1",
	}
}
