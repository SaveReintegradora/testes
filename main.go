package main

import (
	"minha-api/database"
	"minha-api/routes"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	r := routes.SetupRoutes()
	r.Run(":8080")
}
