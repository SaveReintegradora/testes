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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	// Testa a collation em português
	if _, err := conn.Exec(context.Background(), `SELECT * FROM (VALUES ('ázimo'), ('abelha')) AS t(palavra) ORDER BY palavra COLLATE "pt-BR-x-icu"`); err != nil {
		log.Println("Aviso: Collation pt-BR-x-icu não disponível. Usando ordenação padrão.")
	}

	Conn = conn
	log.Println("Conectado ao PostgreSQL!")
}

func CloseDB() {
	Conn.Close(context.Background())
}
