package routes

import (
	"minha-api/controllers"
	middlewares "minha-api/middleware"
	"minha-api/repositories"
	"minha-api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.RateLimitMiddleware()) // Rate limiting global

	// CORS customizado para permitir x-api-key
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "x-api-key"},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: false,
	}
	r.Use(cors.New(corsConfig))

	repo := repositories.NewBookRepository()
	controller := controllers.NewBookController(repo)

	fileRepo := repositories.NewFileProcessRepository()
	fileController := controllers.NewFileProcessController(fileRepo, &utils.RealS3Uploader{}, &utils.RealS3Presigner{})

	clientRepo := repositories.NewClientRepository()
	clientController := controllers.NewClientController(clientRepo)

	clientCRUDController := controllers.NewClientCRUDController(clientRepo)
	clientExportController := controllers.NewClientExportController(clientRepo, &utils.RealS3Uploader{}, &utils.RealS3Presigner{})

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

	// Protege todas as rotas de clientes com API Key
	clients := r.Group("/clients", middlewares.ApiKeyMiddleware())
	{
		clients.POST("/upload", clientController.UploadClients) // upload de clientes
		clients.GET("", clientCRUDController.GetAll)
		clients.GET(":id", clientCRUDController.GetByID)
		clients.POST("", clientCRUDController.Create)
		clients.PUT(":id", clientCRUDController.Update)
		clients.DELETE(":id", clientCRUDController.Delete)
		clients.GET("/export", clientExportController.ExportClients) // exportação de clientes
		clients.GET("/list", func(ctx *gin.Context) {
			clients, err := clientRepo.GetAll()
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, clients)
		}) // endpoint temporário para listar todos os clientes
	}

	return r
}

// Ajuste: Remove interfaces indefinidas e usa tipos concretos dos mocks
func SetupRoutesWithReposAndS3(bookRepo repositories.BookRepositoryInterface, fileRepo repositories.FileProcessRepositoryInterface, s3uploader utils.S3Uploader, s3presigner utils.S3Presigner) *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.RateLimitMiddleware()) // Rate limiting global

	// CORS customizado para permitir x-api-key
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "x-api-key"},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: false,
	}
	r.Use(cors.New(corsConfig))

	controller := controllers.NewBookController(bookRepo)
	fileController := controllers.NewFileProcessController(fileRepo, s3uploader, s3presigner)

	// Adiciona instâncias de repositório e controllers de clientes
	clientRepo := repositories.NewClientRepository()
	clientController := controllers.NewClientController(clientRepo)
	clientCRUDController := controllers.NewClientCRUDController(clientRepo)
	clientExportController := controllers.NewClientExportController(clientRepo, s3uploader, s3presigner)

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

	// Protege todas as rotas de clientes com API Key
	clients := r.Group("/clients", middlewares.ApiKeyMiddleware())
	{
		clients.POST("/upload", clientController.UploadClients) // upload de clientes
		clients.GET("", clientCRUDController.GetAll)
		clients.GET(":id", clientCRUDController.GetByID)
		clients.POST("", clientCRUDController.Create)
		clients.PUT(":id", clientCRUDController.Update)
		clients.DELETE(":id", clientCRUDController.Delete)
		clients.GET("/export", clientExportController.ExportClients) // exportação de clientes
		clients.GET("/list", func(ctx *gin.Context) {
			clients, err := clientRepo.GetAll()
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, clients)
		}) // endpoint temporário para listar todos os clientes
	}

	return r
}

// Alias para facilitar uso nos testes BDD
func SetupRoutesWithMocks(bookRepo repositories.BookRepositoryInterface, fileRepo repositories.FileProcessRepositoryInterface, s3uploader utils.S3Uploader, s3presigner utils.S3Presigner) *gin.Engine {
	return SetupRoutesWithReposAndS3(bookRepo, fileRepo, s3uploader, s3presigner)
}
