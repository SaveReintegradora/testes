package controllers_test

import (
	"minha-api/controllers"
	"minha-api/models"
	"minha-api/tests/mocks"
	"minha-api/utils"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	bookRepo := &mocks.BookRepositoryMock{Books: map[string]models.Book{}}
	fileRepo := &mocks.FileProcessRepositoryMock{Files: map[string]models.FileProcess{}}
	s3mock := &utils.MockS3Uploader{}
	bookController := controllers.NewBookController(bookRepo)
	fileController := controllers.NewFileProcessController(fileRepo, s3mock, &utils.MockS3Presigner{})

	r = gin.Default()
	r.GET("/books", bookController.GetBooks)
	r.POST("/books", bookController.CreateBook)
	// ...adicione outras rotas conforme necessário...
	r.GET("/files", fileController.GetAll)
	r.POST("/files/sendFiles", fileController.Create)
	// ...adicione outras rotas conforme necessário...
}
