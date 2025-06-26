package controllers

import (
	"context"
	"fmt"
	"minha-api/models"
	"minha-api/repositories"
	"minha-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileProcessController struct {
	repo       repositories.FileProcessRepositoryInterface
	s3uploader utils.S3Uploader
}

func NewFileProcessController(repo repositories.FileProcessRepositoryInterface, uploader utils.S3Uploader) *FileProcessController {
	return &FileProcessController{repo: repo, s3uploader: uploader}
}

func (c *FileProcessController) GetAll(ctx *gin.Context) {
	files, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar arquivos"})
		return
	}
	ctx.JSON(http.StatusOK, files)
}

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

func (c *FileProcessController) Create(ctx *gin.Context) {
	fmt.Println(">>> Entrou no método Create do controller")
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
	// s3URL, err := c.s3uploader.UploadToS3(context.Background(), file.Filename, src)
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

func (c *FileProcessController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	type fileUpdateInput struct {
		FileName string `json:"fileName"`
		Name     string `json:"name"`
		FilePath string `json:"filePath,omitempty"`
		Status   string `json:"status,omitempty"`
	}
	var input fileUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	if input.FileName == "" && input.Name != "" {
		input.FileName = input.Name
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
