package controllers

import (
	"minha-api/models"
	"minha-api/repositories"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type ClientController struct {
	repo *repositories.ClientRepository
}

func NewClientController(repo *repositories.ClientRepository) *ClientController {
	return &ClientController{repo: repo}
}

// UploadClients godoc
// @Summary      Upload de clientes via arquivo Excel
// @Description  Recebe um arquivo .xls, lê os dados e cadastra clientes no banco
// @Tags         clients
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Arquivo de clientes (.xls)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /clients/upload [post]
func (c *ClientController) UploadClients(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
		return
	}

	// Verifica extensão
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xls" && ext != ".xlsx" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "O arquivo deve ser .xls ou .xlsx"})
		return
	}

	tempPath := "/tmp/" + uuid.New().String() + ext
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao salvar arquivo temporário"})
		return
	}
	defer os.Remove(tempPath)

	xl, err := excelize.OpenFile(tempPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao abrir arquivo Excel"})
		return
	}

	sheet := xl.GetSheetName(0)
	rows, err := xl.GetRows(sheet)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler linhas do Excel"})
		return
	}

	if len(rows) == 0 || len(rows[0]) < 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modelo de arquivo inválido. Esperado: Nome, Email, Telefone, Endereço"})
		return
	}

	header := rows[0]
	expected := []string{"Nome", "Email", "Telefone", "Endereço"}
	for i, col := range expected {
		if len(header) <= i || strings.TrimSpace(strings.ToLower(header[i])) != strings.TrimSpace(strings.ToLower(col)) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cabeçalho do arquivo inválido. Esperado: Nome, Email, Telefone, Endereço"})
			return
		}
	}

	var count int
	for i, row := range rows {
		if i == 0 {
			continue // pula header
		}
		if len(row) < 4 {
			continue // ignora linhas incompletas
		}
		client := models.Client{
			ID:      uuid.New().String(),
			Name:    row[0],
			Email:   row[1],
			Phone:   row[2],
			Address: row[3],
		}
		if err := c.repo.Create(&client); err == nil {
			count++
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"importados": count})
}
