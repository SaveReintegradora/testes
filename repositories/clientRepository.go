package repositories

import (
	"minha-api/database"
	"minha-api/models"
)

type ClientRepository struct{}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{}
}

func (r *ClientRepository) Create(client *models.Client) error {
	return database.DB.Create(client).Error
}

func (r *ClientRepository) GetAll() ([]models.Client, error) {
	var clients []models.Client
	err := database.DB.Find(&clients).Error
	return clients, err
}
