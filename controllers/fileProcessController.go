package controllers

import (
	"context"
	"minha-api/models"
	"minha-api/repositories"
	"minha-api/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileProcessController struct {
	repo        repositories.FileProcessRepositoryInterface
	s3uploader  utils.S3Uploader
	s3presigner utils.S3Presigner
}

func NewFileProcessController(repo repositories.FileProcessRepositoryInterface, uploader utils.S3Uploader, presigner utils.S3Presigner) *FileProcessController {
	return &FileProcessController{repo: repo, s3uploader: uploader, s3presigner: presigner}
}

// GetAll godoc
// @Summary      Lista todos os arquivos
// @Description  Retorna todos os arquivos processados
// @Tags         files
// @Produce      json
// @Success      200  {array}   models.FileProcess
// @Failure      401  {object}  map[string]string
// @Router       /files [get]
// @Security     ApiKeyAuth
func (c *FileProcessController) GetAll(ctx *gin.Context) {
	files, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar arquivos"})
		return
	}
	ctx.JSON(http.StatusOK, files)
}

// GetByID godoc
// @Summary      Busca arquivo por ID
// @Description  Retorna um arquivo específico pelo ID
// @Tags         files
// @Produce      json
// @Param        id   path      string  true  "ID do arquivo"
// @Success      200  {object}  models.FileProcess
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /files/{id} [get]
// @Security     ApiKeyAuth
func (c *FileProcessController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	file, err := c.repo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, file)
}

// Create godoc
// @Summary      Envia arquivo para processamento
// @Description  Faz upload de um arquivo e registra no sistema
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        nomeArquivo formData file true "Arquivo a ser enviado"
// @Success      201   {object}  models.FileProcess
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /files/sendFiles [post]
// @Security     ApiKeyAuth
func (c *FileProcessController) Create(ctx *gin.Context) {
	file, err := ctx.FormFile("nomeArquivo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado ou inválido"})
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir arquivo"})
		return
	}
	defer src.Close()

	// Upload direto para S3 usando o utilitário
	s3URL, err := c.s3uploader.UploadToS3(context.Background(), file.Filename, src)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar para S3", "details": err.Error()})
		return
	}

	var f models.FileProcess
	f.ID = uuid.New().String()
	f.FileName = file.Filename
	f.FilePath = s3URL
	f.Status = "recebido"
	f.ReceivedAt = time.Now()

	if err := c.repo.Create(&f); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar registro"})
		return
	}
	ctx.JSON(http.StatusCreated, f)
}

// Update godoc
// @Summary      Atualiza um arquivo
// @Description  Atualiza os dados de um arquivo existente
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "ID do arquivo"
// @Param        file  body      models.FileProcess true  "Dados atualizados"
// @Success      200   {object}  models.FileProcess
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /files/{id} [put]
// @Security     ApiKeyAuth
func (c *FileProcessController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	type fileUpdateInput struct {
		FileName string `json:"fileName"`
		FilePath string `json:"filePath,omitempty"`
		Status   string `json:"status,omitempty"`
	}
	var input fileUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	existing, err := c.repo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	if input.FileName != "" {
		existing.FileName = input.FileName
	}
	if input.FilePath != "" {
		existing.FilePath = input.FilePath
	}
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := c.repo.Update(existing); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar registro"})
		return
	}
	ctx.JSON(http.StatusOK, existing)
}

// Delete godoc
// @Summary      Remove um arquivo
// @Description  Remove um arquivo pelo ID
// @Tags         files
// @Param        id   path      string  true  "ID do arquivo"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /files/{id} [delete]
// @Security     ApiKeyAuth
func (c *FileProcessController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	if err := c.repo.Delete(id); err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar registro"})
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}

// DownloadFile godoc
// @Summary      Download do arquivo
// @Description  Realiza o download do arquivo original enviado para o S3
// @Tags         files
// @Produce      octet-stream
// @Param        id   path      string  true  "ID do arquivo"
// @Success      302  {string}  string  "Redirect para o arquivo no S3"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /files/{id}/download [get]
// @Security     ApiKeyAuth
func (c *FileProcessController) DownloadFile(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	file, err := c.repo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	if file.FilePath == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Arquivo sem URL de download"})
		return
	}
	bucket := os.Getenv("AWS_BUCKET_NAME")
	key := file.FileName
	url, err := c.s3presigner.PresignGetObject(ctx, bucket, key, 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar link de download"})
		return
	}
	ctx.Header("Location", url)
	ctx.Status(http.StatusFound)
}
