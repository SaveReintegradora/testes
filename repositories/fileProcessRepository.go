package repositories

import (
	"context"
	"errors"
	"minha-api/database"
	"minha-api/models"
)

type FileProcessRepository struct{}

func NewFileProcessRepository() *FileProcessRepository {
	return &FileProcessRepository{}
}

func (r *FileProcessRepository) GetAll() ([]models.FileProcess, error) {
	rows, err := database.Conn.Query(
		context.Background(),
		`SELECT id, file_name, file_path, received_at, status, error_msg FROM file_processes ORDER BY received_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.FileProcess
	for rows.Next() {
		var f models.FileProcess
		if err := rows.Scan(&f.ID, &f.FileName, &f.FilePath, &f.ReceivedAt, &f.Status, &f.ErrorMsg); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *FileProcessRepository) GetByID(id string) (*models.FileProcess, error) {
	var f models.FileProcess
	err := database.Conn.QueryRow(
		context.Background(),
		`SELECT id, file_name, file_path, received_at, status, error_msg FROM file_processes WHERE id = $1`,
		id,
	).Scan(&f.ID, &f.FileName, &f.FilePath, &f.ReceivedAt, &f.Status, &f.ErrorMsg)
	if err != nil {
		return nil, errors.New("not found")
	}
	return &f, nil
}

func (r *FileProcessRepository) Create(f *models.FileProcess) error {
	_, err := database.Conn.Exec(
		context.Background(),
		`INSERT INTO file_processes (id, file_name, file_path, received_at, status, error_msg) VALUES ($1, $2, $3, $4, $5, $6)`,
		f.ID, f.FileName, f.FilePath, f.ReceivedAt, f.Status, f.ErrorMsg,
	)
	return err
}

func (r *FileProcessRepository) Update(f *models.FileProcess) error {
	cmd, err := database.Conn.Exec(
		context.Background(),
		`UPDATE file_processes SET file_name=$1, file_path=$2, received_at=$3, status=$4, error_msg=$5 WHERE id=$6`,
		f.FileName, f.FilePath, f.ReceivedAt, f.Status, f.ErrorMsg, f.ID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *FileProcessRepository) Delete(id string) error {
	cmd, err := database.Conn.Exec(
		context.Background(),
		`DELETE FROM file_processes WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

type FileProcessRepositoryInterface interface {
	GetAll() ([]models.FileProcess, error)
	GetByID(id string) (*models.FileProcess, error)
	Create(f *models.FileProcess) error
	Update(f *models.FileProcess) error
	Delete(id string) error
}
