package models

import "time"

type FileProcess struct {
	ID         string    `json:"id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	ReceivedAt time.Time `json:"received_at"`
	Status     string    `json:"status"` // pendente, em processamento, concluido com erros, concluido sem erros
	ErrorMsg   string    `json:"error_msg,omitempty"`
}
