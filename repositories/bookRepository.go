package repositories

import (
	"minha-api/database"
	"minha-api/models"
)

type BookRepository struct{}

func NewBookRepository() *BookRepository {
	return &BookRepository{}
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	result := database.DB.Find(&books)
	return books, result.Error
}

func (r *BookRepository) GetByID(id string) (*models.Book, error) {
	var book models.Book
	result := database.DB.First(&book, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func (r *BookRepository) Create(book *models.Book) error {
	return database.DB.Create(book).Error
}

func (r *BookRepository) Update(book *models.Book) error {
	return database.DB.Save(book).Error
}

func (r *BookRepository) Delete(id string) error {
	return database.DB.Delete(&models.Book{}, "id = ?", id).Error
}

type BookRepositoryInterface interface {
	GetAll() ([]models.Book, error)
	GetByID(id string) (*models.Book, error)
	Create(book *models.Book) error
	Update(book *models.Book) error
	Delete(id string) error
}
