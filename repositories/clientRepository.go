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

func (r *ClientRepository) ExistsByNameAndCPF(name, cpf string) (bool, error) {
	var count int64
	type Result struct {
		Name string
		CPF  string
	}
	// Se o campo CPF não existir no model, sempre retorna false
	db := database.DB.Table("clients").Where("name = ? AND cpf = ?", name, cpf)
	db.Count(&count)
	return count > 0, db.Error
}

func (r *ClientRepository) ExistsByNameAndCNPJ(name, cnpj string) (bool, error) {
	var count int64
	db := database.DB.Table("clients").Where("name = ? AND cnpj = ?", name, cnpj)
	db.Count(&count)
	return count > 0, db.Error
}

func (r *ClientRepository) ExistsByNameAndEmail(name, email string) (bool, error) {
	var count int64
	db := database.DB.Table("clients").Where("name = ? AND email = ?", name, email)
	db.Count(&count)
	return count > 0, db.Error
}
