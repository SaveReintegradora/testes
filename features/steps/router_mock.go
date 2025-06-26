package steps

import (
	"minha-api/repositories"
	"minha-api/routes"
	"minha-api/utils"
	"net/http"
)

var BookRepoMock *repositories.BookRepositoryMock
var FileRepoMock *repositories.FileProcessRepositoryMock

func setupRouterWithMocks() http.Handler {
	BookRepoMock = repositories.NewBookRepositoryMock()
	FileRepoMock = repositories.NewFileProcessRepositoryMock()
	s3mock := &utils.MockS3Uploader{}
	return routes.SetupRoutesWithMocks(BookRepoMock, FileRepoMock, s3mock)
}
