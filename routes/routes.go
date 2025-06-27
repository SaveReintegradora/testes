package routes

import (
	"minha-api/controllers"
	middlewares "minha-api/middleware"
	"minha-api/repositories"
	"minha-api/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	repo := repositories.NewBookRepository()
	controller := controllers.NewBookController(repo)

	fileRepo := repositories.NewFileProcessRepository()
	fileController := controllers.NewFileProcessController(fileRepo, &utils.RealS3Uploader{})

	books := r.Group("/books", middlewares.ApiKeyMiddleware())
	{
		books.GET("", controller.GetBooks)
		books.GET(":id", controller.GetBookByID)
		books.POST("", controller.CreateBook)
		books.PUT(":id", controller.UpdateBook)
		books.DELETE(":id", controller.DeleteBook)
	}

	files := r.Group("/files", middlewares.ApiKeyMiddleware())
	{
		files.GET("", fileController.GetAll)
		files.GET(":id", fileController.GetByID)
		files.POST("sendFiles", fileController.Create)
		files.PUT(":id", fileController.Update)
		files.DELETE(":id", fileController.Delete)
		files.GET(":id/download", fileController.DownloadFile) // nova rota de download
	}

	return r
}

// Ajuste: Remove interfaces indefinidas e usa tipos concretos dos mocks
func SetupRoutesWithReposAndS3(bookRepo repositories.BookRepositoryInterface, fileRepo repositories.FileProcessRepositoryInterface, s3uploader utils.S3Uploader) *gin.Engine {
	r := gin.Default()

	controller := controllers.NewBookController(bookRepo)
	fileController := controllers.NewFileProcessController(fileRepo, s3uploader)

	books := r.Group("/books", middlewares.ApiKeyMiddleware())
	{
		books.GET("", controller.GetBooks)
		books.GET(":id", controller.GetBookByID)
		books.POST("", controller.CreateBook)
		books.PUT(":id", controller.UpdateBook)
		books.DELETE(":id", controller.DeleteBook)
	}

	files := r.Group("/files", middlewares.ApiKeyMiddleware())
	{
		files.GET("", fileController.GetAll)
		files.GET(":id", fileController.GetByID)
		files.POST("sendFiles", fileController.Create)
		files.PUT(":id", fileController.Update)
		files.DELETE(":id", fileController.Delete)
		files.GET(":id/download", fileController.DownloadFile) // nova rota de download
	}

	return r
}

// Alias para facilitar uso nos testes BDD
func SetupRoutesWithMocks(bookRepo repositories.BookRepositoryInterface, fileRepo repositories.FileProcessRepositoryInterface, s3uploader utils.S3Uploader) *gin.Engine {
	return SetupRoutesWithReposAndS3(bookRepo, fileRepo, s3uploader)
}
