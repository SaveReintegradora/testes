package main

import (
	"fmt"
	"minha-api/database"
	_ "minha-api/docs" // ajuste para o nome do seu módulo
	"minha-api/routes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	fmt.Println("Inicializando a API...")
	database.InitDB()
	defer database.CloseDB()

	r := routes.SetupRoutes()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
