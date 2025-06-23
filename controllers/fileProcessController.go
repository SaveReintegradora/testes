package controllers

import (
	"minha-api/models"
	"minha-api/repositories"

	//	"minha-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileProcessController struct {
	repo *repositories.FileProcessRepository
}

func NewFileProcessController(repo *repositories.FileProcessRepository) *FileProcessController {
	return &FileProcessController{repo: repo}
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

	// s3URL, err := utils.UploadToS3(ctx, file.Filename, src)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar para S3", "details": err.Error()})
	// 	return
	// }
	s3URL := "mocked-s3-url" // TODO: Replace with actual S3 upload logic or import the correct package

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
	var f models.FileProcess
	if err := ctx.ShouldBindJSON(&f); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	f.ID = id
	if err := c.repo.Update(&f); err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar registro"})
		}
		return
	}
	ctx.JSON(http.StatusOK, f)
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
