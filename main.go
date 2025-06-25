package main

import (
	"fmt"
	"log"
	"minha-api/database"
	_ "minha-api/docs" // ajuste para o nome do seu módulo
	"minha-api/routes"
	"os"
	"time"

	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func waitForDB(maxAttempts int, delay time.Duration) {
	for i := 1; i <= maxAttempts; i++ {
		err := database.PingDB()
		if err == nil {
			return
		}
		log.Printf("Tentativa %d: Banco ainda não disponível (%v)", i, err)
		time.Sleep(delay)
	}
	log.Fatal("Banco de dados não respondeu após várias tentativas.")
}

func main() {
	// Só carrega o .env fora do Docker
	if os.Getenv("RUNNING_IN_DOCKER") != "1" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Aviso: .env não carregado (pode ser normal em produção):", err)
		}
	}

	// Aguarda o banco de dados ficar disponível de forma robusta
	waitForDB(15, 2*time.Second)

	fmt.Println("Inicializando a API...")
	database.InitDB()
	defer database.CloseDB()

	r := routes.SetupRoutes()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run() // Inicia o servidor na porta padrão 8080 (ou a porta definida na variável de ambiente PORT)    docker compose restart
}
