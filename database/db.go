package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func InitDB() {
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL não definida")
	}

	var err error
	Conn, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
}

func CloseDB() {
	if Conn != nil {
		Conn.Close(context.Background())
	}
}
