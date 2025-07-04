package database

import (
	"fmt"
	"minha-api/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "postgres://allan:agripa99@books_db:5432/minha_api_books?sslmode=disable"
	}
	fmt.Println("[INFO] DSN utilizado para conexão:", dsn)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Erro ao conectar no banco com GORM:", err)
		panic(err)
	}
	// Migração automática
	DB.AutoMigrate(&models.Client{})
}
