package models

import "gorm.io/gorm"

type Client struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
	CNPJ      string         `json:"cnpj"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
