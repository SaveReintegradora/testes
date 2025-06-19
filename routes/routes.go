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

	fileRepo := repositories.NewFileProcessRepository()
	fileController := controllers.NewFileProcessController(fileRepo)

	books := r.Group("/books", middlewares.ApiKeyMiddleware())
	{
		books.GET("", controller.GetBooks)
		books.GET("/:id", controller.GetBookByID)
		books.POST("", controller.CreateBook)
		books.PUT("/:id", controller.UpdateBook)
		books.DELETE("/:id", controller.DeleteBook)
	}

	files := r.Group("/files", middlewares.ApiKeyMiddleware())
	{
		files.GET("", fileController.GetAll)
		files.GET("/:id", fileController.GetByID)
		files.POST("", fileController.Create)
		files.PUT("/:id", fileController.Update)
		files.DELETE("/:id", fileController.Delete)
	}

	return r
}
