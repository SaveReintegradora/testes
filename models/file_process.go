package models

import (
	"time"

	"gorm.io/gorm"
)

type FileProcess struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FileName   string         `json:"fileName"`
	FilePath   string         `json:"file_path"`
	ReceivedAt time.Time      `json:"received_at"`
	Status     string         `json:"status"` // pendente, em processamento, concluido com erros, concluido sem erros
	ErrorMsg   string         `json:"error_msg,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
