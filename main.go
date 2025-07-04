package main

import (
	"log"
	"minha-api/database"
	_ "minha-api/docs" // ajuste para o nome do seu módulo
	"minha-api/routes"
	"os"

	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Só carrega o .env fora do Docker
	if os.Getenv("RUNNING_IN_DOCKER") != "1" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Aviso: .env não carregado (pode ser normal em produção):", err)
		}
	}

	database.ConnectDatabase() // Use apenas ConnectDatabase()

	r := routes.SetupRoutes()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	certFile := os.Getenv("SSL_CERT_FILE")
	keyFile := os.Getenv("SSL_KEY_FILE")
	if certFile != "" && keyFile != "" {
		// HTTPS ativado
		r.RunTLS(":5000", certFile, keyFile)
	} else {
		// HTTP padrão
		r.Run(":5000")
	}
}
