package routes

import (
	"minha-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	books := r.Group("/books")
	{
		books.GET("", controllers.GetBooks)
		books.GET("/:id", controllers.GetBookByID)
		books.POST("", controllers.CreateBook)
		books.DELETE("/:id", controllers.DeleteBook)
	}

	return r
}
