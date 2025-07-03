package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"minha-api/models"
)

var DB *gorm.DB

func InitGorm() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "host=books_db user=allan dbname=minha_api_books sslmode=disable password=agripa99"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco com GORM: %v", err)
	}
	DB = db

	// Migração automática das tabelas
	err = DB.AutoMigrate(&models.Book{}, &models.FileProcess{}, &models.Client{})
	if err != nil {
		log.Fatalf("Erro ao migrar tabelas: %v", err)
	}
}
