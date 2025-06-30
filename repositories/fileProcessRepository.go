package repositories

import (
	"minha-api/database"
	"minha-api/models"
)

type FileProcessRepository struct{}

func NewFileProcessRepository() *FileProcessRepository {
	return &FileProcessRepository{}
}

func (r *FileProcessRepository) GetAll() ([]models.FileProcess, error) {
	var files []models.FileProcess
	result := database.DB.Order("received_at DESC").Find(&files)
	return files, result.Error
}

func (r *FileProcessRepository) GetByID(id string) (*models.FileProcess, error) {
	var f models.FileProcess
	result := database.DB.First(&f, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &f, nil
}

func (r *FileProcessRepository) Create(f *models.FileProcess) error {
	return database.DB.Create(f).Error
}

func (r *FileProcessRepository) Update(f *models.FileProcess) error {
	return database.DB.Save(f).Error
}

func (r *FileProcessRepository) Delete(id string) error {
	return database.DB.Delete(&models.FileProcess{}, "id = ?", id).Error
}

type FileProcessRepositoryInterface interface {
	GetAll() ([]models.FileProcess, error)
	GetByID(id string) (*models.FileProcess, error)
	Create(f *models.FileProcess) error
	Update(f *models.FileProcess) error
	Delete(id string) error
}
