package main

import (
	"minha-api/routes"
)

func main() {
	r := routes.SetupRoutes()
	r.Run(":8080") // Sobe na porta 8080
}
