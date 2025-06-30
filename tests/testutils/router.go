package testutils

import (
	"minha-api/repositories"
	"minha-api/routes"
	"minha-api/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	return routes.SetupRoutes()
}

func SetupRouterWithMocks() *gin.Engine {
	bookRepo := repositories.NewBookRepositoryMock()
	fileRepo := repositories.NewFileProcessRepositoryMock()
	s3mock := &utils.MockS3Uploader{}
	s3presign := &utils.MockS3Presigner{}
	return routes.SetupRoutesWithReposAndS3(bookRepo, fileRepo, s3mock, s3presign)
}

func SetupRouterWithReposAndS3(bookRepo repositories.BookRepositoryInterface, fileRepo repositories.FileProcessRepositoryInterface, s3uploader utils.S3Uploader, s3presigner utils.S3Presigner) *gin.Engine {
	return routes.SetupRoutesWithReposAndS3(bookRepo, fileRepo, s3uploader, s3presigner)
}

func GetFileProcessRepositoryMock() *repositories.FileProcessRepositoryMock {
	return repositories.NewFileProcessRepositoryMock()
}
