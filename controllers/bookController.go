package controllers

import (
	"minha-api/models"
	"minha-api/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookController struct {
	repo *repositories.BookRepository
}

func NewBookController(repo *repositories.BookRepository) *BookController {
	return &BookController{repo: repo}
}

// GetBooks godoc
// @Summary      Lista todos os livros
// @Description  Retorna todos os livros cadastrados
// @Tags         books
// @Produce      json
// @Success      200  {array}   models.Book
// @Failure      401  {object}  map[string]string
// @Router       /books [get]
// @Security     ApiKeyAuth
func (c *BookController) GetBooks(ctx *gin.Context) {
	books, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar livros"})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

// GetBookByID godoc
// @Summary      Busca livro por ID
// @Description  Retorna um livro específico pelo ID
// @Tags         books
// @Produce      json
// @Param        id   path      string  true  "ID do livro"
// @Success      200  {object}  models.Book
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /books/{id} [get]
// @Security     ApiKeyAuth
func (c *BookController) GetBookByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	book, err := c.repo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary      Cria um novo livro
// @Description  Adiciona um novo livro ao banco de dados
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      models.Book  true  "Livro a ser criado"
// @Success      201   {object}  models.Book
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /books [post]
// @Security     ApiKeyAuth
func (c *BookController) CreateBook(ctx *gin.Context) {
	var newBook models.Book

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if newBook.Title == "" || newBook.Author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Título e autor são obrigatórios"})
		return
	}

	newBook.ID = uuid.New().String()
	newBook.CreatedAt = time.Now()

	if err := c.repo.Create(&newBook); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar livro"})
		return
	}

	ctx.JSON(http.StatusCreated, newBook)
}

// UpdateBook godoc
// @Summary      Atualiza um livro
// @Description  Atualiza os dados de um livro existente
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      string      true  "ID do livro"
// @Param        book  body      models.Book true  "Dados atualizados"
// @Success      200   {object}  models.Book
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /books/{id} [put]
// @Security     ApiKeyAuth
func (c *BookController) UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedBook models.Book
	if err := ctx.ShouldBindJSON(&updatedBook); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	updatedBook.ID = id

	if err := c.repo.Update(&updatedBook); err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar livro"})
		}
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

// DeleteBook godoc
// @Summary      Remove um livro
// @Description  Remove um livro pelo ID
// @Tags         books
// @Param        id   path      string  true  "ID do livro"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /books/{id} [delete]
// @Security     ApiKeyAuth
func (c *BookController) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.repo.Delete(id); err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Livro não encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar livro"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
