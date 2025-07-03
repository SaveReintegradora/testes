package database

import (
	"fmt"
	"minha-api/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=books_db user=allan dbname=minha_api_books sslmode=disable password=allan"
	}
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Erro ao conectar no banco com GORM:", err)
		panic(err)
	}
	// Migração automática
	DB.AutoMigrate(&models.Client{})
}
