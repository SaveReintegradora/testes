package controllers

import (
	"minha-api/models"
	"minha-api/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientCRUDController struct {
	repo *repositories.ClientRepository
}

func NewClientCRUDController(repo *repositories.ClientRepository) *ClientCRUDController {
	return &ClientCRUDController{repo: repo}
}

// GetAllClients godoc
// @Summary      Lista todos os clientes
// @Description  Retorna todos os clientes cadastrados
// @Tags         clients
// @Produce      json
// @Success      200 {array} models.Client
// @Router       /clients [get]
func (c *ClientCRUDController) GetAll(ctx *gin.Context) {
	clients, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
		return
	}
	ctx.JSON(http.StatusOK, clients)
}

// GetClientByID godoc
// @Summary      Busca cliente por ID
// @Description  Retorna um cliente pelo ID
// @Tags         clients
// @Produce      json
// @Param        id path string true "ID do cliente" example("b3e1c2d0-1234-4abc-9def-1234567890ab")
// @Success      200 {object} models.Client
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /clients/{id} [get]
func (c *ClientCRUDController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido. Use um UUID válido."})
		return
	}
	client, err := c.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar cliente"})
		}
		return
	}
	ctx.JSON(http.StatusOK, client)
}

// CreateClient godoc
// @Summary      Cria um novo cliente
// @Description  Cria um cliente a partir de um JSON
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        client body models.Client true "Dados do cliente"
// @Success      201 {object} models.Client
// @Failure      400 {object} map[string]string
// @Router       /clients [post]
func (c *ClientCRUDController) Create(ctx *gin.Context) {
	var client models.Client
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	if err := c.repo.Create(&client); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar cliente"})
		return
	}
	ctx.JSON(http.StatusCreated, client)
}

// UpdateClient godoc
// @Summary      Atualiza um cliente
// @Description  Atualiza os dados de um cliente pelo ID
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do cliente"
// @Param        client body models.Client true "Dados do cliente"
// @Success      200 {object} models.Client
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /clients/{id} [put]
func (c *ClientCRUDController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var client models.Client
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	client.ID = id
	if err := c.repo.Update(&client); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar cliente"})
		return
	}
	ctx.JSON(http.StatusOK, client)
}

// DeleteClient godoc
// @Summary      Deleta um cliente
// @Description  Remove um cliente pelo ID
// @Tags         clients
// @Produce      json
// @Param        id path string true "ID do cliente"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /clients/{id} [delete]
func (c *ClientCRUDController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.repo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar cliente"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cliente deletado com sucesso"})
}
