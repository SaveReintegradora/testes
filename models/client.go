package models

import "gorm.io/gorm"

// Client representa um cliente do sistema.
type Client struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string         `json:"name" validate:"required,min=2"`                                              // Nome do cliente (obrigatório, mínimo 2 caracteres)
	Email     string         `gorm:"uniqueIndex:idx_email_unique" json:"email" validate:"required,email"`         // Email único e válido
	Phone     string         `json:"phone" validate:"required"`                                                   // Telefone obrigatório
	Address   string         `json:"address" validate:"required"`                                                 // Endereço obrigatório
	CNPJ      string         `gorm:"uniqueIndex:idx_cnpj_unique" json:"cnpj" validate:"omitempty,len=14,numeric"` // CNPJ único, 14 dígitos
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                                              // Soft delete
}
