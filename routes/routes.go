package routes

import (
	"minha-api/controllers"
	middlewares "minha-api/middleware"
	"minha-api/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	repo := repositories.NewBookRepository()
	controller := controllers.NewBookController(repo)

	books := r.Group("/books", middlewares.ApiKeyMiddleware())
	{
		books.GET("", controller.GetBooks)
		books.GET("/:id", controller.GetBookByID)
		books.POST("", controller.CreateBook)
		books.PUT("/:id", controller.UpdateBook)
		books.DELETE("/:id", controller.DeleteBook)
	}

	return r
}
