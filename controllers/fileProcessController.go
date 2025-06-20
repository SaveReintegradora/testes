package controllers

import (
	"minha-api/models"
	"minha-api/repositories"
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
	var f models.FileProcess
	if err := ctx.ShouldBindJSON(&f); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	if f.FileName == "" || f.FilePath == "" || f.Status == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Campos obrigatórios faltando"})
		return
	}
	f.ID = uuid.New().String()
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
