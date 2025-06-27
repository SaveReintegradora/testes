package steps

import (
	"minha-api/tests/mocks"
)

var BookRepoMock *mocks.BookRepositoryMock
var FileRepoMock *mocks.FileProcessRepositoryMock

func SetupMocks() {
	BookRepoMock = mocks.NewBookRepositoryMock()
	FileRepoMock = mocks.NewFileProcessRepositoryMock()
}
