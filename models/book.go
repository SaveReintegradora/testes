package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title     string         `json:"title"`
	Author    string         `json:"author"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
