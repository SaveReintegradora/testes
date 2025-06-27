package mocks

import (
	"errors"
	"minha-api/models"
	"minha-api/repositories"
	"time"
)

type FileProcessRepositoryMock struct {
	Files map[string]models.FileProcess
}

// Garante que FileProcessRepositoryMock implementa FileProcessRepositoryInterface
var _ repositories.FileProcessRepositoryInterface = (*FileProcessRepositoryMock)(nil)

func NewFileProcessRepositoryMock() *FileProcessRepositoryMock {
	return &FileProcessRepositoryMock{
		Files: map[string]models.FileProcess{
			"1": {ID: "1", FileName: "mock.txt", FilePath: "https://mock-s3.local/mock.txt", ReceivedAt: time.Now(), Status: "recebido"},
		},
	}
}

func (m *FileProcessRepositoryMock) GetAll() ([]models.FileProcess, error) {
	files := make([]models.FileProcess, 0, len(m.Files))
	for _, f := range m.Files {
		files = append(files, f)
	}
	return files, nil
}

func (m *FileProcessRepositoryMock) GetByID(id string) (*models.FileProcess, error) {
	if f, ok := m.Files[id]; ok {
		return &f, nil
	}
	return nil, errors.New("not found")
}

func (m *FileProcessRepositoryMock) Create(f *models.FileProcess) error {
	m.Files[f.ID] = *f
	return nil
}

func (m *FileProcessRepositoryMock) Update(f *models.FileProcess) error {
	if _, ok := m.Files[f.ID]; ok {
		m.Files[f.ID] = *f
		return nil
	}
	return errors.New("not found")
}

func (m *FileProcessRepositoryMock) Delete(id string) error {
	if _, ok := m.Files[id]; ok {
		delete(m.Files, id)
		return nil
	}
	return errors.New("not found")
}

func (m *FileProcessRepositoryMock) Reset() {
	m.Files = map[string]models.FileProcess{}
	m.Files["1"] = models.FileProcess{
		ID:       "1",
		FileName: "mock.txt",
		FilePath: "https://mock-s3.local/mock.txt",
		Status:   "recebido",
	}
}
