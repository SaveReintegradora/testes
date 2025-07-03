package controllers

import (
	"context"
	"fmt"
	"minha-api/repositories"
	"minha-api/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ClientExportController struct {
	repo        *repositories.ClientRepository
	s3uploader  utils.S3Uploader
	s3presigner utils.S3Presigner
}

func NewClientExportController(repo *repositories.ClientRepository, uploader utils.S3Uploader, presigner utils.S3Presigner) *ClientExportController {
	return &ClientExportController{repo: repo, s3uploader: uploader, s3presigner: presigner}
}

// ExportClients godoc
// @Summary      Exporta todos os clientes em XLS e salva no S3
// @Description  Gera um arquivo XLS com todos os clientes do banco e retorna um link temporário para download do arquivo salvo no S3. O arquivo contém as colunas: ID, Nome, Email, Telefone, Endereço, CNPJ.
// @Tags         clients
// @Produce      json
// @Success      200 {object} map[string]string "Exemplo de resposta: {\"download_url\":\"https://bucket.s3.amazonaws.com/clientes_export_20250703_153000.xlsx\"}"
// @Failure      500 {object} map[string]string "Erro ao buscar clientes, gerar XLS ou enviar para S3"
// @Router       /clients/export [get]
func (c *ClientExportController) ExportClients(ctx *gin.Context) {
	clients, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
		return
	}

	xl := excelize.NewFile()
	sheet := "Clientes"
	idx, err := xl.NewSheet(sheet)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar nova aba no XLS"})
		return
	}
	xl.SetActiveSheet(idx)
	headers := []string{"ID", "Nome", "Email", "Telefone", "Endereço", "CNPJ"}
	for i, h := range headers {
		col, _ := excelize.CoordinatesToCellName(i+1, 1)
		xl.SetCellValue(sheet, col, h)
	}
	for idx, client := range clients {
		row := idx + 2
		xl.SetCellValue(sheet, fmt.Sprintf("A%d", row), client.ID)
		xl.SetCellValue(sheet, fmt.Sprintf("B%d", row), client.Name)
		xl.SetCellValue(sheet, fmt.Sprintf("C%d", row), client.Email)
		xl.SetCellValue(sheet, fmt.Sprintf("D%d", row), client.Phone)
		xl.SetCellValue(sheet, fmt.Sprintf("E%d", row), client.Address)
		xl.SetCellValue(sheet, fmt.Sprintf("F%d", row), client.CNPJ)
	}

	fileName := "clientes_export_" + time.Now().Format("20060102_150405") + ".xlsx"
	buf, err := xl.WriteToBuffer()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar XLS"})
		return
	}

	_, err = c.s3uploader.UploadToS3(context.Background(), fileName, buf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar XLS para S3", "details": err.Error()})
		return
	}

	bucket := os.Getenv("AWS_BUCKET_NAME")
	bucket = strings.TrimSpace(strings.Trim(bucket, "."))
	presignedURL, err := c.s3presigner.PresignGetObject(context.Background(), bucket, fileName, 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar link temporário para download", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"download_url": presignedURL})
}
