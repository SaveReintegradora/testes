package repositories

import (
	"context"
	"errors"
	"minha-api/database"
	"minha-api/models"
)

type BookRepository struct{}

func NewBookRepository() *BookRepository {
	return &BookRepository{}
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	rows, err := database.Conn.Query(
		context.Background(),
		`SELECT id, title, author, created_at 
		 FROM books 
		 ORDER BY title COLLATE "pt-BR-x-icu"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *BookRepository) GetByID(id string) (*models.Book, error) {
	var book models.Book
	err := database.Conn.QueryRow(
		context.Background(),
		`SELECT id, title, author, created_at 
		 FROM books WHERE id = $1`,
		id,
	).Scan(&book.ID, &book.Title, &book.Author, &book.CreatedAt)

	if err != nil {
		return nil, errors.New("not found")
	}
	return &book, nil
}

func (r *BookRepository) Create(book *models.Book) error {
	_, err := database.Conn.Exec(
		context.Background(),
		`INSERT INTO books (id, title, author) 
		 VALUES ($1, $2, $3)`,
		book.ID, book.Title, book.Author,
	)
	return err
}

func (r *BookRepository) Update(book *models.Book) error {
	cmd, err := database.Conn.Exec(
		context.Background(),
		`UPDATE books 
		 SET title = $1, author = $2 
		 WHERE id = $3`,
		book.Title, book.Author, book.ID,
	)

	if cmd.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return err
}

func (r *BookRepository) Delete(id string) error {
	cmd, err := database.Conn.Exec(
		context.Background(),
		`DELETE FROM books WHERE id = $1`,
		id,
	)

	if cmd.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return err
}
